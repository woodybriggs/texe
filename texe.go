package texe

import (
	"context"

	"github.com/woodybriggs/texe/queues"
	"github.com/woodybriggs/texe/types"
)

type Texe struct {
	types.TexeOpts
	workers chan struct{}
}

func defaultOpts() types.TexeOpts {
	return types.TexeOpts{
		MaxWorkers: 8,
		Queue:      queues.NewFifoSliceQueue(8),
	}
}

func WithMaxWorkers(max uint) func(*types.TexeOpts) {
	return func(to *types.TexeOpts) {
		to.MaxWorkers = max
	}
}

func WithQueue(queue types.Queue) func(*types.TexeOpts) {
	return func(to *types.TexeOpts) {
		to.Queue = queue
	}
}

func NewTexe(opts ...func(*types.TexeOpts)) *Texe {

	o := defaultOpts()
	for _, fn := range opts {
		fn(&o)
	}

	return &Texe{
		TexeOpts: o,
		workers:  make(chan struct{}, o.MaxWorkers),
	}
}

func (tex *Texe) QueueTask(task *types.Task) (*types.TaskRunInfo, error) {
	taskruninfo := &types.TaskRunInfo{
		Task:   *task,
		Status: types.TexeStatus_Unknown,
		Error:  nil,
	}

	err := tex.Queue.Enqueue(taskruninfo)
	if err != nil {
		taskruninfo.Status = types.TexeStatus_Error
		taskruninfo.Error = err
		return taskruninfo, err
	}

	taskruninfo.Status = types.TexeStatus_Queued

	return taskruninfo, nil
}

func (tex *Texe) StartWithContext(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			{
				return ctx.Err()
			}
		default:
			{
				taskruninfo := tex.Queue.Dequeue()
				if taskruninfo == nil {
					continue
				}

				tex.workers <- struct{}{}
				taskruninfo.Status = types.TexeStatus_Running
				taskruninfo.Exe.TaskStartingCallback(taskruninfo)

				go func(tri *types.TaskRunInfo) {
					err := tri.Exe.Start(tri)
					if err != nil {
						tri.Error = err
						tri.Status = types.TexeStatus_Error
					} else {
						tri.Status = types.TexeStatus_Complete
					}
					tri.Exe.TaskCompleteCallback(tri, err)
					<-tex.workers
				}(taskruninfo)
			}
		}
	}
}
