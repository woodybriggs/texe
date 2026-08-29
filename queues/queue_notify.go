package queues

import (
	"context"

	"github.com/woodybriggs/texe/types"
)

// QueueWithNotify wraps a Queue and signals a channel on every Enqueue,
// allowing consumers to block until work is available instead of polling.
type QueueWithNotify struct {
	types.Queue
	Ready chan struct{}
}

// NewQueueWithNotify wraps the given queue and returns a QueueWithNotify.
// The Ready channel is buffered to 1 so Enqueue never blocks on the signal.
func NewQueueWithNotify(q types.Queue) *QueueWithNotify {
	return &QueueWithNotify{
		Queue: q,
		Ready: make(chan struct{}, 1),
	}
}

func (w *QueueWithNotify) Enqueue(ctx context.Context, task *types.TaskRunInfo) error {
	if err := w.Queue.Enqueue(ctx, task); err != nil {
		return err
	}
	select {
	case w.Ready <- struct{}{}:
	default:
	}
	return nil
}
