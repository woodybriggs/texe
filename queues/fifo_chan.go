package queues

import "github.com/woodybriggs/texe/types"

type FifoChanQueue struct {
	types.Queue
	items chan *types.TaskContext
}

func NewFifoChanQueue(buffsize int) *FifoChanQueue {
	return &FifoChanQueue{
		items: make(chan *types.TaskContext, buffsize),
	}
}

func (q *FifoChanQueue) Enqueue(task *types.TaskContext) error {
	q.items <- task
	return nil
}

func (q *FifoChanQueue) Dequeue() *types.TaskContext {
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
