package main

type (
	Queue struct {
		Branches []*string `json:"branches"`
	}

	MergeQueue []*Queue

	MergeQueueWorkflows struct {
		MergeQueue MergeQueue
	}
)

// return true if the branch exit.
func (q *Queue) next() bool {
	if len(q.Branches) == 0 {
		return false
	}

	return true
}

// pop the branch to queue.
func (q *Queue) pop() *string {
	if q.next() {
		branch := q.Branches[0]
		q.Branches = q.Branches[1:]

		return branch
	}

	return nil
}

// push the branch to queue.
func (q *Queue) push(branch string) {
	q.Branches = append(q.Branches, &branch)
}

// front returns the branch at the front of the queue.
func (q *Queue) front() *string {
	if q.is_empty() {
		return nil
	}

	return q.Branches[0]
}

// is_empty returns true if the queue is empty.
func (q *Queue) is_empty() bool {
	return len(q.Branches) == 0
}

// size returns the number of branches in the queue.
func (q *Queue) size() int {
	return len(q.Branches)
}

// New creates a new Queue.
func NewQueue() *Queue {
	return &Queue{Branches: []*string{}}
}
