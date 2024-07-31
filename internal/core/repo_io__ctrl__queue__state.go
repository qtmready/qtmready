package core

import (
	"go.temporal.io/sdk/workflow"

	"go.breu.io/quantm/internal/shared"
)

type (
	// Node represents a single node in the queue.
	Node struct {
		pr   *RepoIOPullRequest
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
		PR       *RepoIOPullRequest `json:"pr"`
		Position int                `json:"position"`
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
func (q *Queue) push(ctx workflow.Context, pr *RepoIOPullRequest) {
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

	return pr
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
// initial state: a <-> b <-> c <-> d
// promote node d
//
// Example:
//
//	q := NewQueue()
//	ctx := workflow.Context{}
//	a := RepoIOPullRequest{Number: 1}
//	b := RepoIOPullRequest{Number: 2}
//	c := RepoIOPullRequest{Number: 3}
//	d := RepoIOPullRequest{Number: 4}
//
//	q.push(ctx, a)
//	q.push(ctx, b)
//	q.push(ctx, c)
//	q.push(ctx, d)
//
//	q.promote(d)
//
//	the queue should now be: a <-> b <-> d <-> c
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
// initial state: a <-> b <-> c <-> d,
// demote node b.
//
// Example:
//
//	q := NewQueue()
//	ctx := workflow.Context{}
//	a := RepoIOPullRequest{Number: 1}
//	b := RepoIOPullRequest{Number: 2}
//	c := RepoIOPullRequest{Number: 3}
//	d := RepoIOPullRequest{Number: 4}
//
//	q.push(ctx, a)
//	q.push(ctx, b)
//	q.push(ctx, c)
//	q.push(ctx, d)
//
//	q.demote(b)
//
//	the queue should now be: a <-> c <-> b <-> d
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

// remove removes a specific pull request from the queue.
func (q *Queue) remove(ctx workflow.Context, prNumber int64) {
	_ = q.mutex.Lock(ctx)
	defer q.mutex.Unlock()

	node, exists := q.index[prNumber]
	if !exists {
		return // Item not in queue, do nothing
	}

	if node.prev != nil {
		node.prev.next = node.next
	} else {
		q.head = node.next
	}

	if node.next != nil {
		node.next.prev = node.prev
	} else {
		q.tail = node.prev
	}

	delete(q.index, prNumber)
}

/**
 * QueueCtrlState methods
 */

// Signal handlers

// on_add handles the addition of a new pull request to the primary queue.
//
// Usage:
//
//	state := NewQueueCtrlState(ctx, repo, branch)
//	add_handler := state.on_add(ctx)
//	selector.AddReceive(add_channel, add_handler)
func (s *QueueCtrlState) on_add(ctx workflow.Context) shared.ChannelHandler {
	return func(c workflow.ReceiveChannel, more bool) {
		payload := &RepoIOPullRequest{}

		s.rx(ctx, c, payload)
		s.push(ctx, payload, false)
	}
}

// on_add_priority handles the addition of a new pull request to the priority queue.
//
// Usage:
//
//	state := NewQueueCtrlState(ctx, repo, branch)
//	add_priority_handler := state.on_add_priority(ctx)
//	selector.AddReceive(add_priority_channel, add_priority_handler)
func (s *QueueCtrlState) on_add_priority(ctx workflow.Context) shared.ChannelHandler {
	return func(c workflow.ReceiveChannel, more bool) {
		payload := &RepoIOPullRequest{}

		s.rx(ctx, c, payload)
		s.push(ctx, payload, true)
	}
}

// on_promote handles the promotion of a pull request in the primary queue.
//
// Usage:
//
//	state := NewQueueCtrlState(ctx, repo, branch)
//	promote_handler := state.on_promote(ctx)
//	selector.AddReceive(promote_channel, promote_handler)
func (s *QueueCtrlState) on_promote(ctx workflow.Context) shared.ChannelHandler {
	return func(c workflow.ReceiveChannel, more bool) {
		payload := &RepoIOPullRequest{}

		s.rx(ctx, c, payload)
		s.primary.reorder(ctx, *payload, true)
	}
}

// on_demote handles the demotion of a pull request in the primary queue.
//
// Usage:
//
//	state := NewQueueCtrlState(ctx, repo, branch)
//	demote_handler := state.on_demote(ctx)
//	selector.AddReceive(demote_channel, demote_handler)
func (s *QueueCtrlState) on_demote(ctx workflow.Context) shared.ChannelHandler {
	return func(c workflow.ReceiveChannel, more bool) {
		payload := &RepoIOPullRequest{}

		s.rx(ctx, c, payload)
		s.primary.reorder(ctx, *payload, false)
	}
}

/**
 * Other QueueCtrlState methods
 */

// push adds a new pull request to either the primary or priority queue.
//
// Example:
//
//	state := NewQueueCtrlState(ctx, repo, branch)
//	pr := RepoIOPullRequest{Number: 123}
//	state.push(ctx, pr, false) // Add to primary queue
//	state.push(ctx, pr, true)  // Add to priority queue
func (s *QueueCtrlState) push(ctx workflow.Context, pr *RepoIOPullRequest, urgent bool) {
	_ = s.mutex.Lock(ctx)
	defer s.mutex.Unlock()

	if urgent {
		s.priority.push(ctx, pr)
	} else {
		s.primary.push(ctx, pr)
	}
}

// next waits for the next item to be available in either queue.
//
// Example:
//
//	state := NewQueueCtrlState(ctx, repo, branch)
//	err := state.next(ctx)
//	if err != nil {
//	    // Handle error
//	}
func (s *QueueCtrlState) next(ctx workflow.Context) error {
	return workflow.Await(ctx, func() bool {
		return s.priority.peek() || s.primary.peek()
	})
}

// pop removes and returns the next pull request from either queue.
//
// Example:
//
//	state := NewQueueCtrlState(ctx, repo, branch)
//	pr := state.pop(ctx)
//	if pr != nil {
//	    // Process the pull request
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
func NewQueue(ctx workflow.Context) Queue {
	return Queue{
		mutex: workflow.NewMutex(ctx),
		head:  &Node{},
		tail:  &Node{},
		index: make(map[int64]*Node),
	}
}

// NewQueueCtrlState creates a new QueueCtrlState and sets the branch.
// It returns the updated context and the new QueueCtrlState.
//
// Example:
//
//	ctx := workflow.Context{}
//	repo := &Repo{}
//	branch := "main"
//	ctx, state := NewQueueCtrlState(ctx, repo, branch)
func NewQueueCtrlState(ctx workflow.Context, repo *Repo, branch string) (workflow.Context, *QueueCtrlState) {
	ctrl := &QueueCtrlState{
		base_ctrl: NewBaseCtrl(ctx, "queue_ctrl", repo),
		primary:   NewQueue(ctx),
		priority:  NewQueue(ctx),
	}

	return ctrl.set_branch(ctx, branch), ctrl
}
