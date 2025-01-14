package scheduler

type Queue struct {
	tasks []ImageTask		// Array of ImageTask type
}

// NewQueue creates a new queue.
func NewQueue() *Queue {
	return &Queue{}
}

func (q *Queue) GetTasks() []ImageTask {
	return q.tasks
}

// IsEmpty checks if the queue is empty.
func (q *Queue) GetLength() int {
	return len(q.tasks)
}

// Enqueue adds a task to the end of the queue.
func (q *Queue) Enqueue(task ImageTask) {
	q.tasks = append(q.tasks, task)
}

// Dequeue removes and returns a task from the front of the queue.
// Returns Task, bool - the bool is false if the queue is empty.
func (q *Queue) Dequeue() (ImageTask) {
	task := q.tasks[0]  // Get the first task
	q.tasks = q.tasks[1:]  // Remove it from the queue
	return task
}

// IsEmpty checks if the queue is empty.
func (q *Queue) IsEmpty() bool {
	return len(q.tasks) == 0
}
