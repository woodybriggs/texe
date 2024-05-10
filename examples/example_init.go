package main

import (
	"context"
	"fmt"

	"github.com/woodybriggs/texe"
	"github.com/woodybriggs/texe/queues"
	"github.com/woodybriggs/texe/types"
)

type SayHelloTaskExe struct {
	types.TaskExe
	Hello string
}

func (task SayHelloTaskExe) Start(texctx *types.TaskContext) error {
	fmt.Println("Hello", task.Hello)
	return nil
}

func main() {

	taskcount := 10

	tex := texe.NewTexe(
		texe.WithMaxWorkers(1),
		texe.WithQueue(queues.NewFifoChanQueue(taskcount)),
	)

	for id := range taskcount {
		_, err := tex.QueueTask(&types.TaskDef{
			Description: "say hello",
			Exe: &SayHelloTaskExe{
				Hello: fmt.Sprint(id),
			},
		})
		fmt.Println(err)
	}

	tex.StartWithContext(context.Background())
}
