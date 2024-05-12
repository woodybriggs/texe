package main

import (
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/woodybriggs/texe"
	"github.com/woodybriggs/texe/queues"
	"github.com/woodybriggs/texe/types"
)

type SayHelloTaskExe struct {
	types.TaskExe
	Id    string
	Hello string
}

func (task *SayHelloTaskExe) TaskStartingCallback(info *types.TaskRunInfo) {
	fmt.Println("Starting", task.Id)
}

func (task *SayHelloTaskExe) Start(info *types.TaskRunInfo) error {
	fmt.Println("Hello", task.Id, task.Hello)
	return nil
}

func (task *SayHelloTaskExe) TaskCompleteCallback(info *types.TaskRunInfo, err error) {
	fmt.Println("task complete", info)
	if err != nil {
		log.Println(err)
	}
}

func main() {

	taskcount := 10

	tex := texe.NewTexe(
		texe.WithMaxWorkers(1),
		texe.WithQueue(queues.NewFifoChanQueue(taskcount)),
	)

	for id := range taskcount {

		ident, err := uuid.NewV7()

		task := &types.Task{
			Description: "say hello",
			Exe: &SayHelloTaskExe{
				Id:    ident.String(),
				Hello: fmt.Sprint(id),
			},
		}
		_, err = tex.QueueTask(task)
		fmt.Println(err)
	}

	tex.StartWithContext(context.Background())
}
