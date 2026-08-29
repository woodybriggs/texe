package queues

import (
	"testing"

	"github.com/woodybriggs/texe/types"
)

func TestFifoChanQueue_EnqueueDequeue(t *testing.T) {
	queue := NewFifoChanQueue(8)

	task1 := &types.TaskRunInfo{Status: types.TexeStatus_Queued}
	task2 := &types.TaskRunInfo{Status: types.TexeStatus_Queued}

	queue.Enqueue(task1)
	queue.Enqueue(task2)

	result := queue.Dequeue()
	if result != task1 {
		t.Errorf("expected task1, got %v", result)
	}

	result = queue.Dequeue()
	if result != task2 {
		t.Errorf("expected task2, got %v", result)
	}
}

func TestFifoChanQueue_EmptyDequeue(t *testing.T) {
	queue := NewFifoChanQueue(8)

	result := queue.Dequeue()
	if result != nil {
		t.Errorf("expected nil, got %v", result)
	}
}
