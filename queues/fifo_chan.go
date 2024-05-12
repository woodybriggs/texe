package queues

import "github.com/woodybriggs/texe/types"

type FifoChanQueue struct {
	types.Queue
	items chan *types.TaskRunInfo
}

func NewFifoChanQueue(buffsize int) *FifoChanQueue {
	return &FifoChanQueue{
		items: make(chan *types.TaskRunInfo, buffsize),
	}
}

func (q *FifoChanQueue) Enqueue(task *types.TaskRunInfo) error {
	q.items <- task
	return nil
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
