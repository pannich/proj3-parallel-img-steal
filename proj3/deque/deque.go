package deque

import (
	"image"
	"sync/atomic"
)

type Task struct {
	Bounds image.Rectangle
}

// Double-ended queue
//Start at bottom=0, top=0, the length of the queue never changed. It is the length of the array.
//PushBottom, bottom++. (+ is growing downwards)
//When popbottom, bottom-- (shrink upwards). When poptop, top++, meaning the next element to poptop is the one lower.

type DEQueue struct {
	Tasks  []*Task
	top    int32
	bottom int32
	stamp  int64 // Atomic stamp to manage versioning and detect ABA problems
}

func NewDEQueue(size int) *DEQueue {
	return &DEQueue{
		Tasks:  make([]*Task, size),
		top:    0,
		bottom: 0,
		stamp:  0,
	}
}

func (q *DEQueue) GetTop() int32 {
	return atomic.LoadInt32(&q.top)
}

func (q *DEQueue) GetBottom() int32 {
	return atomic.LoadInt32(&q.bottom)
}

// PushBottom
// : noconcurrency when push task
func (q *DEQueue) PushBottom(task *Task) bool {
	localBottom := atomic.LoadInt32(&q.bottom)
	if localBottom == int32(len(q.Tasks)) {
		return false // The queue is full, indicate failure
	}
	// Prepare to store the task
	q.Tasks[localBottom] = task
	atomic.AddInt32(&q.bottom, 1)
	return true
}

func (q *DEQueue) PopTop() (*Task, bool) {
	localTop := atomic.LoadInt32(&q.top)
	localBottom := atomic.LoadInt32(&q.bottom)
	if localBottom <= localTop {
		return nil, false // Queue is empty or in an inconsistent state
	}

	task := q.Tasks[localTop]
	oldStamp := atomic.LoadInt64(&q.stamp)
	newStamp := oldStamp + 1

	if atomic.CompareAndSwapInt32(&q.top, localTop, localTop+1) {
		atomic.CompareAndSwapInt64(&q.stamp, oldStamp, newStamp)
		return task, true
	}
	return nil, false
}

func (q *DEQueue) PopBottom() (*Task, bool) {
	localBottom := atomic.LoadInt32(&q.bottom)
	if localBottom == 0 {
		return nil, false // Queue is empty
	}
	localBottom--
	atomic.StoreInt32(&q.bottom, localBottom)

	localTop := atomic.LoadInt32(&q.top)
	task := q.Tasks[localBottom]
	oldStamp := atomic.LoadInt64(&q.stamp)
	newStamp := oldStamp + 1

	// Top and bottom one or mar apart, no conflict
	if localBottom > localTop {
		return task, true
	}
	if localBottom == localTop { // last element
		// if I win, bottom is 0. If I lose, thief must have won, bottom is 0.
		atomic.StoreInt32(&q.bottom, 0)
		if atomic.CompareAndSwapInt32(&q.top, localTop, 0) {
			atomic.CompareAndSwapInt64(&q.stamp, oldStamp, newStamp)
			return task, true
		}
	}
	// Failed to pop last task, restore the bottom and top since DEQueue is empty
	atomic.StoreInt32(&q.top, 0)
	atomic.StoreInt64(&q.stamp, int64(newStamp))
	atomic.AddInt32(&q.bottom, 0)
	return nil, false
}

func (q *DEQueue) IsEmpty() bool {
	localTop := atomic.LoadInt32(&q.top)
	localBottom := atomic.LoadInt32(&q.bottom)
	return localTop >= localBottom
}
