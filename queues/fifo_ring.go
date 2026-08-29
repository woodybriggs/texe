package queues

import (
	"runtime"
	"sync"

	"github.com/woodybriggs/texe/types"
)

type FifoRingQueue struct {
	types.Queue
	mu     sync.Mutex
	buffer []*types.TaskRunInfo
	head   int
	tail   int
}

func NewFifoRingQueue(capacity int) *FifoRingQueue {
	return &FifoRingQueue{
		buffer: make([]*types.TaskRunInfo, capacity),
		head:   0,
		tail:   0,
	}
}

func (queue *FifoRingQueue) expand() {
	newhead := cap(queue.buffer)
	newcap := cap(queue.buffer) * 2
	newbuf := make([]*types.TaskRunInfo, newcap)
	copy(newbuf, queue.buffer)
	runtime.GC()
	queue.buffer = newbuf
	queue.head = newhead
}

func (queue *FifoRingQueue) Enqueue(task *types.TaskRunInfo) error {
	queue.mu.Lock()
	defer queue.mu.Unlock()
	if queue.buffer[queue.head] != nil {
		queue.expand()
	}

	queue.buffer[queue.head] = task
	if (queue.head + 1) == cap(queue.buffer) {
		queue.head = 0
	} else {
		queue.head++
	}

	return nil
}

func (queue *FifoRingQueue) Dequeue() *types.TaskRunInfo {
	queue.mu.Lock()
	defer queue.mu.Unlock()
	if queue.head == queue.tail {
		return nil
	}
	task := queue.buffer[queue.tail]
	queue.buffer[queue.tail] = nil

	if (queue.tail + 1) == cap(queue.buffer) {
		queue.tail = 0
	} else {
		queue.tail++
	}

	return task
}

func (queue *FifoRingQueue) Len() int {
	if queue.head >= queue.tail {
		return queue.head - queue.tail
	}
	return cap(queue.buffer) - queue.tail + queue.head
}
