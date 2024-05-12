package queues

import (
	"github.com/woodybriggs/texe/types"
)

type FifoSliceQueue struct {
	types.Queue
	items []*types.TaskRunInfo
}

func NewFifoSliceQueue(buffsize int) *FifoSliceQueue {
	return &FifoSliceQueue{
		items: make([]*types.TaskRunInfo, 0, buffsize),
	}
}

func (fifo *FifoSliceQueue) Enqueue(ctx *types.TaskRunInfo) error {
	fifo.items = append(fifo.items, ctx)
	return nil
}

func (fifo *FifoSliceQueue) Dequeue() *types.TaskRunInfo {
	count := len(fifo.items)
	if count < 1 {
		return nil
	}

	popped := fifo.items[0]
	fifo.items = fifo.items[1:]

	return popped
}
