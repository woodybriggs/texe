package queues

import (
	"context"

	"github.com/woodybriggs/texe/types"
)

type FifoChanQueue struct {
	types.Queue
	items chan *types.TaskRunInfo
}

func NewFifoChanQueue(buffsize int) *FifoChanQueue {
	return &FifoChanQueue{
		items: make(chan *types.TaskRunInfo, buffsize),
	}
}

func (q *FifoChanQueue) Enqueue(ctx context.Context, task *types.TaskRunInfo) error {
	select {
	case q.items <- task:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (q *FifoChanQueue) Dequeue() *types.TaskRunInfo {
	select {
	case task := <-q.items:
		{
			return task
		}
	default:
		{
			return nil
		}
	}
}

func (q *FifoChanQueue) Len() int {
	return len(q.items)
}
