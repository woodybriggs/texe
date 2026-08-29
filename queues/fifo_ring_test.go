package queues

import (
	"context"
	"sync"
	"testing"

	"github.com/woodybriggs/texe/types"
)

func TestFifoRingQueue_EnqueueDequeue(t *testing.T) {
	queue := NewFifoRingQueue(4)

	task1 := &types.TaskRunInfo{Status: types.TexeStatus_Queued}
	task2 := &types.TaskRunInfo{Status: types.TexeStatus_Queued}

	queue.Enqueue(context.Background(), task1)
	queue.Enqueue(context.Background(), task2)

	result := queue.Dequeue()
	if result != task1 {
		t.Errorf("expected task1, got %v", result)
	}

	result = queue.Dequeue()
	if result != task2 {
		t.Errorf("expected task2, got %v", result)
	}
}

func TestFifoRingQueue_EmptyDequeue(t *testing.T) {
	queue := NewFifoRingQueue(4)

	result := queue.Dequeue()
	if result != nil {
		t.Errorf("expected nil, got %v", result)
	}
}

func TestFifoRingQueue_Expand(t *testing.T) {
	queue := NewFifoRingQueue(2)

	task1 := &types.TaskRunInfo{Status: types.TexeStatus_Queued}
	task2 := &types.TaskRunInfo{Status: types.TexeStatus_Queued}
	task3 := &types.TaskRunInfo{Status: types.TexeStatus_Queued}

	queue.Enqueue(context.Background(), task1)
	queue.Enqueue(context.Background(), task2)
	queue.Enqueue(context.Background(), task3) // This should trigger expand

	result := queue.Dequeue()
	if result != task1 {
		t.Errorf("expected task1, got %v", result)
	}

	result = queue.Dequeue()
	if result != task2 {
		t.Errorf("expected task2, got %v", result)
	}

	result = queue.Dequeue()
	if result != task3 {
		t.Errorf("expected task3, got %v", result)
	}
}

func TestFifoRingQueue_Ordering(t *testing.T) {
	queue := NewFifoRingQueue(8)

	for i := 0; i < 10; i++ {
		queue.Enqueue(context.Background(), &types.TaskRunInfo{Status: types.TexeStatus_Queued})
	}

	for i := 0; i < 10; i++ {
		result := queue.Dequeue()
		if result == nil {
			t.Fatalf("expected non-nil at position %d", i)
		}
	}

	if queue.Dequeue() != nil {
		t.Error("expected nil after dequeuing all items")
	}
}

func TestFifoRingQueue_ConcurrentAccess(t *testing.T) {
	queue := NewFifoRingQueue(16)
	var wg sync.WaitGroup

	// Concurrent enqueues
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			queue.Enqueue(context.Background(), &types.TaskRunInfo{Status: types.TexeStatus_Queued})
		}()
	}
	wg.Wait()

	// Concurrent dequeues - no panic or data race is the success condition
	dequeued := make(chan *types.TaskRunInfo, 50)
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			result := queue.Dequeue()
			if result != nil {
				dequeued <- result
			}
		}()
	}
	wg.Wait()
	close(dequeued)

	count := 0
	for range dequeued {
		count++
	}
	if count == 0 {
		t.Error("expected at least one dequeued item")
	}
}
