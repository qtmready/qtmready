package core

import (
	"go.temporal.io/sdk/workflow"
)

type (
	// Node represents a single node in the queue.
	Node struct {
		pr   RepoIOPullRequest
		prev *Node
		next *Node
	}

	// Queue represents a thread-safe queue of pull requests.
	Queue struct {
		mutex workflow.Mutex
		head  *Node
		tail  *Node
		index map[int64]*Node // Map PR number to Node
	}

	// QueueCtrlState represents the state of the queue controller,
	// managing both a primary and a priority queue.
	QueueCtrlState struct {
		*base_ctrl
		primary  Queue
		priority Queue
	}

	// QueueItem represents a single item in the queue for frontend representation.
	QueueItem struct {
		PR       RepoIOPullRequest `json:"pr"`
		Position int               `json:"position"`
	}
)

// push adds a new pull request to the end of the queue.
//
// Example:
//
//	q := NewQueue()
//	ctx := workflow.Context{}
//	pr := RepoIOPullRequest{Number: 123}
//	q.push(ctx, pr)
func (q *Queue) push(ctx workflow.Context, pr RepoIOPullRequest) {
	_ = q.mutex.Lock(ctx)
	defer q.mutex.Unlock()

	node := &Node{pr: pr}
	if q.tail == nil {
		q.head = node
		q.tail = node
	} else {
		node.prev = q.tail
		q.tail.next = node
		q.tail = node
	}

	q.index[pr.Number.Int64()] = node
}

// pop removes and returns the first pull request in the queue.
// Returns nil if the queue is empty.
//
// Example:
//
//	q := NewQueue()
//	ctx := workflow.Context{}
//	pr := q.pop(ctx)
//	if pr != nil {
//	    fmt.Printf("Popped PR number: %d\n", pr.Number)
//	}
func (q *Queue) pop(ctx workflow.Context) *RepoIOPullRequest {
	_ = q.mutex.Lock(ctx)
	defer q.mutex.Unlock()

	if q.head == nil {
		return nil
	}

	pr := q.head.pr
	q.head = q.head.next

	if q.head == nil {
		q.tail = nil
	} else {
		q.head.prev = nil
	}

	delete(q.index, pr.Number.Int64())

	return &pr
}

// peek returns true if the queue is not empty, false otherwise.
//
// Example:
//
//	q := NewQueue()
//	if q.peek() {
//	    fmt.Println("Queue is not empty")
//	}
func (q *Queue) peek() bool {
	return q.head != nil
}

// reorder moves a pull request one position up or down in the queue.
// If promote is true, the item is moved up; if false, it's moved down.
// Silently ignores if the item is already at the top/bottom of the queue.
//
// Example:
//
//	q := NewQueue()
//	ctx := workflow.Context{}
//	pr := RepoIOPullRequest{Number: 123}
//	q.reorder(ctx, pr, true) // Promote PR
func (q *Queue) reorder(ctx workflow.Context, pr RepoIOPullRequest, promote bool) {
	_ = q.mutex.Lock(ctx)
	defer q.mutex.Unlock()

	node, exists := q.index[pr.Number.Int64()]
	if !exists {
		return // Item not in queue, do nothing
	}

	if promote {
		q.promote(node)
	} else {
		q.demote(node)
	}
}

// promote moves a Node one position up in the queue.
func (q *Queue) promote(node *Node) {
	if node.prev != nil {
		prev_prev := node.prev.prev
		next := node.next

		node.prev.next = next
		node.prev.prev = node
		node.next = node.prev
		node.prev = prev_prev

		if prev_prev != nil {
			prev_prev.next = node
		} else {
			q.head = node
		}

		if next != nil {
			next.prev = node.next
		} else {
			q.tail = node.next
		}
	}
}

// demote moves a Node one position down in the queue.
func (q *Queue) demote(node *Node) {
	if node.next != nil {
		prev := node.prev
		next_next := node.next.next

		node.next.prev = prev
		node.next.next = node
		node.prev = node.next
		node.next = next_next

		if prev != nil {
			prev.next = node.prev
		} else {
			q.head = node.prev
		}

		if next_next != nil {
			next_next.prev = node
		} else {
			q.tail = node
		}
	}
}

// items returns a list representation of the queue for frontend use.
//
// Example:
//
//	q := NewQueue()
//	ctx := workflow.Context{}
//	items := q.items(ctx)
//	for _, item := range items {
//	    fmt.Printf("PR %d at position %d\n", item.PR.Number, item.Position)
//	}
func (q *Queue) items(ctx workflow.Context) []QueueItem {
	_ = q.mutex.Lock(ctx)
	defer q.mutex.Unlock()

	items := make([]QueueItem, 0, len(q.index))
	position := 1

	for node := q.head; node != nil; node = node.next {
		items = append(items, QueueItem{
			PR:       node.pr,
			Position: position,
		})
		position++
	}

	return items
}

// length returns the number of items in the queue.
//
// Example:
//
//	q := NewQueue()
//	fmt.Printf("Queue length: %d\n", q.length())
func (q *Queue) length() int {
	return len(q.index)
}

// push adds a new pull request to either the primary or priority queue.
//
// Example:
//
//	s := NewQueueCtrlState(base)
//	ctx := workflow.Context{}
//	pr := RepoIOPullRequest{Number: 123}
//	s.push(ctx, pr, true) // Push to priority queue
func (s *QueueCtrlState) push(ctx workflow.Context, pr RepoIOPullRequest, urgent bool) {
	_ = s.mutex.Lock(ctx)
	defer s.mutex.Unlock()

	if urgent {
		s.priority.push(ctx, pr)
	} else {
		s.primary.push(ctx, pr)
	}
}

// peek returns true if any of the queues (priority or primary) has an item in it.
//
// Example:
//
//	s := NewQueueCtrlState(base)
//	if s.peek() {
//	    fmt.Println("At least one queue has items")
//	}
func (s *QueueCtrlState) peek() bool {
	return s.priority.peek() || s.primary.peek()
}

// pop removes and returns the next pull request from the queues.
// It prioritizes the priority queue over the primary queue.
//
// Example:
//
//	s := NewQueueCtrlState(base)
//	ctx := workflow.Context{}
//	pr := s.pop(ctx)
//	if pr != nil {
//	    fmt.Printf("Popped PR number: %d\n", pr.Number)
//	}
func (s *QueueCtrlState) pop(ctx workflow.Context) *RepoIOPullRequest {
	_ = s.mutex.Lock(ctx)
	defer s.mutex.Unlock()

	if s.priority.peek() {
		return s.priority.pop(ctx)
	}

	return s.primary.pop(ctx)
}

// NewQueue creates a new Queue.
//
// Example:
//
//	q := NewQueue()
func NewQueue() Queue {
	return Queue{
		index: make(map[int64]*Node),
	}
}

// NewQueueCtrlState creates a new QueueCtrlState.
//
// Example:
//
//	base := &base_ctrl{}
//	s := NewQueueCtrlState(base)
func NewQueueCtrlState(base *base_ctrl) *QueueCtrlState {
	return &QueueCtrlState{
		base_ctrl: base,
		primary:   NewQueue(),
		priority:  NewQueue(),
	}
}
