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

func (tex *Texe) QueueTask(def *types.TaskDef) (*types.TaskContext, error) {
	texctx := &types.TaskContext{
		Def:    def,
		Status: types.TexeStatus_Unknown,
	}

	err := tex.Queue.Enqueue(texctx)
	if err != nil {
		texctx.Status = types.TexeStatus_Error
		texctx.Error = err
		return texctx, err
	}

	texctx.Status = types.TexeStatus_Queued

	return texctx, nil
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
				texctx := tex.Queue.Dequeue()
				if texctx == nil {
					continue
				}

				tex.workers <- struct{}{}

				go func(tc *types.TaskContext) {
					tc.Status = types.TexeStatus_Running
					err := tc.Def.Exe.Start(tc)
					<-tex.workers
					if err != nil {
						tc.Error = err
						tc.Status = types.TexeStatus_Error
					} else {
						tc.Status = types.TexeStatus_Complete
					}
				}(texctx)
			}
		}
	}
}
