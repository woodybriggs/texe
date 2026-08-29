package queues

import (
	"context"
	"sync"

	"github.com/woodybriggs/texe/types"
)

type FifoSliceQueue struct {
	types.Queue
	mu    sync.Mutex
	items []*types.TaskRunInfo
}

func NewFifoSliceQueue(buffsize int) *FifoSliceQueue {
	return &FifoSliceQueue{
		items: make([]*types.TaskRunInfo, 0, buffsize),
	}
}

func (fifo *FifoSliceQueue) Enqueue(ctx context.Context, task *types.TaskRunInfo) error {
	fifo.mu.Lock()
	defer fifo.mu.Unlock()
	fifo.items = append(fifo.items, task)
	return nil
}

func (fifo *FifoSliceQueue) Dequeue() *types.TaskRunInfo {
	fifo.mu.Lock()
	defer fifo.mu.Unlock()
	count := len(fifo.items)
	if count < 1 {
		return nil
	}

	popped := fifo.items[0]
	fifo.items = fifo.items[1:]

	return popped
}

func (fifo *FifoSliceQueue) Len() int {
	return len(fifo.items)
}
