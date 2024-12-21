package states

import (
	"go.temporal.io/sdk/workflow"
)

type (
	// Node[E events.Payload] represents a doubly linked list node containing a payload of type E.
	Node[E any] struct {
		Item     *E       `json:"item"`     // Pointer to the payload item.
		Previous *Node[E] `json:"previous"` // Pointer to the previous node in the list.
		Next     *Node[E] `json:"next"`     // Pointer to the next node in the list.
	}

	// Sequencer[K comparable, E events.Payload] provides a thread-safe, FIFO queue with indexed access.
	// It utilizes a doubly linked list for queue management and a map for O(1) key-based lookup.
	Sequencer[K comparable, E any] struct {
		Head *Node[E]       `json:"head"` // Pointer to the head (front) of the queue.
		Tail *Node[E]       `json:"tail"` // Pointer to the tail (back) of the queue.
		Map  map[K]*Node[E] `json:"map"`  // Map providing key-to-node associations.

		mutex workflow.Mutex // mutex for thread-safe operations.
	}
)

// - Queue Manipulation -

// Push adds an item to the end of the queue.
func (q *Sequencer[K, E]) Push(ctx workflow.Context, key K, item *E) {
	_ = q.mutex.Lock(ctx)
	defer q.mutex.Unlock()

	node := &Node[E]{Item: item}
	if q.Tail == nil {
		q.Head = node
		q.Tail = node
	} else {
		q.Tail.Next = node
		node.Previous = q.Tail
		q.Tail = node
	}

	q.Map[key] = node
}

// Priority adds an item to the front of the queue.
func (q *Sequencer[K, E]) Priority(ctx workflow.Context, key K, item *E) {
	_ = q.mutex.Lock(ctx)
	defer q.mutex.Unlock()

	node := &Node[E]{Item: item}
	if q.Head == nil {
		q.Head = node
		q.Tail = node
	} else {
		node.Next = q.Head
		q.Head.Previous = node
		q.Head = node
	}

	q.Map[key] = node
}

// Pop removes and returns the item at the front of the queue.
func (q *Sequencer[K, E]) Pop(ctx workflow.Context) *E {
	_ = q.mutex.Lock(ctx)
	defer q.mutex.Unlock()

	if q.Head == nil {
		return nil
	}

	node := q.Head
	q.Head = node.Next

	if q.Head != nil {
		q.Head.Previous = nil
	} else {
		q.Tail = nil
	}

	return node.Item
}

// Remove removes a specific item from the queue based on its key.
func (q *Sequencer[K, E]) Remove(ctx workflow.Context, key K) {
	_ = q.mutex.Lock(ctx)
	defer q.mutex.Unlock()

	node, ok := q.Map[key]
	if !ok {
		return
	}

	delete(q.Map, key)

	if node.Previous != nil {
		node.Previous.Next = node.Next
	} else {
		q.Head = node.Next
	}

	if node.Next != nil {
		node.Next.Previous = node.Previous
	} else {
		q.Tail = node.Previous
	}
}

// - Queue Item Reordering -

// Promote moves an item one position forward in the queue.
func (q *Sequencer[K, E]) Promote(ctx workflow.Context, key K) {
	_ = q.mutex.Lock(ctx)
	defer q.mutex.Unlock()

	node, ok := q.Map[key]
	if !ok || node.Previous == nil {
		return
	}

	prev := node.Previous

	if prev.Previous != nil {
		prev.Previous.Next = node
	} else {
		q.Head = node
	}

	if node.Next != nil {
		node.Next.Previous = prev
	} else {
		q.Tail = prev
	}

	node.Previous = prev.Previous
	prev.Next = node.Next
	node.Next = prev
	prev.Previous = node
}

// Demote moves an item one position backward in the queue.
func (q *Sequencer[K, E]) Demote(ctx workflow.Context, key K) {
	_ = q.mutex.Lock(ctx)
	defer q.mutex.Unlock()

	node, ok := q.Map[key]
	if !ok || node.Next == nil {
		return // Node not found or already at the tail
	}

	next := node.Next

	if next.Next != nil {
		next.Next.Previous = node
	} else {
		q.Tail = node
	}

	if node.Previous != nil {
		node.Previous.Next = next
	} else {
		q.Head = next
	}

	node.Next = next.Next
	next.Previous = node.Previous
	node.Previous = next
	next.Next = node
}

// - Queue Inspection -

// Peek returns the item at the front of the queue without removing it.
func (q *Sequencer[K, E]) Peek(ctx workflow.Context) *E {
	return q.Head.Item
}

// Position returns the position of the key in the queue (starting from 1).
// Returns 0 if the key is not found.
func (q *Sequencer[K, E]) Position(ctx workflow.Context, key K) int {
	node, ok := q.Map[key]
	if !ok {
		return 0
	}

	position := 1

	for current := q.Head; current != nil; current = current.Next {
		if current == node {
			return position
		}

		position++
	}

	// We should never reach this point.
	return 0
}

// Length returns the number of items in the queue.
func (q *Sequencer[K, E]) Length(ctx workflow.Context) int {
	length := 0
	for current := q.Head; current != nil; current = current.Next {
		length++
	}

	return length
}

// All returns all items in the queue.
func (q *Sequencer[K, E]) All(ctx workflow.Context) []*E {
	_ = q.mutex.Lock(ctx)
	defer q.mutex.Unlock()

	items := make([]*E, 0)
	for current := q.Head; current != nil; current = current.Next {
		items = append(items, current.Item)
	}

	return items
}

// - Initialization and Creation -

// Init restores the lock mutex.
func (q *Sequencer[K, E]) Init(ctx workflow.Context) {
	q.mutex = workflow.NewMutex(ctx)
}

// NewSequencer[K, E] creates a new Sequencer.
func NewSequencer[K comparable, E any]() *Sequencer[K, E] {
	return &Sequencer[K, E]{
		Map: make(map[K]*Node[E]),
	}
}
