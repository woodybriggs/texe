package texe

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/woodybriggs/texe/queues"
	"github.com/woodybriggs/texe/types"
)

type mockTaskExe struct {
	startErr     error
	startCalled  bool
	startingCall bool
	completeCall bool
}

func (m *mockTaskExe) TaskStartingCallback(info *types.TaskRunInfo) {
	m.startingCall = true
}

func (m *mockTaskExe) Start(info *types.TaskRunInfo) error {
	m.startCalled = true
	return m.startErr
}

func (m *mockTaskExe) TaskCompleteCallback(info *types.TaskRunInfo, err error) {
	m.completeCall = true
}

func TestTexe_QueueTask(t *testing.T) {
	engine := NewTexe()
	task := &types.Task{
		Exe:         &mockTaskExe{},
		Description: "test task",
	}

	info, err := engine.QueueTask(context.Background(), task)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if info.Status != types.TexeStatus_Queued {
		t.Errorf("expected status Queued, got %v", info.Status)
	}
}

func TestTexe_StartWithContext_ExecutesTask(t *testing.T) {
	exe := &mockTaskExe{}
	engine := NewTexe(
		WithMaxWorkers(2),
		WithQueue(queues.NewFifoSliceQueue(8)),
	)

	task := &types.Task{
		Exe:         exe,
		Description: "test task",
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	_, err := engine.QueueTask(context.Background(), task)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	go func() {
		time.Sleep(100 * time.Millisecond)
		cancel()
	}()

	engine.StartWithContext(ctx)

	if !exe.startCalled {
		t.Error("expected Start to be called")
	}
	if !exe.startingCall {
		t.Error("expected TaskStartingCallback to be called")
	}
	if !exe.completeCall {
		t.Error("expected TaskCompleteCallback to be called")
	}
}

func TestTexe_StartWithContext_TaskError(t *testing.T) {
	expectedErr := errors.New("task failed")
	exe := &mockTaskExe{startErr: expectedErr}
	engine := NewTexe(
		WithMaxWorkers(2),
		WithQueue(queues.NewFifoSliceQueue(8)),
	)

	task := &types.Task{
		Exe:         exe,
		Description: "failing task",
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	info, err := engine.QueueTask(context.Background(), task)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	go func() {
		time.Sleep(100 * time.Millisecond)
		cancel()
	}()

	engine.StartWithContext(ctx)

	if info.Error != expectedErr {
		t.Errorf("expected error %v, got %v", expectedErr, info.Error)
	}
	if info.Status != types.TexeStatus_Error {
		t.Errorf("expected status Error, got %v", info.Status)
	}
}

func TestTexe_StartWithContext_StatusTransitions(t *testing.T) {
	exe := &mockTaskExe{}
	engine := NewTexe(
		WithMaxWorkers(2),
		WithQueue(queues.NewFifoSliceQueue(8)),
	)

	task := &types.Task{
		Exe:         exe,
		Description: "status test",
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	info, _ := engine.QueueTask(context.Background(), task)

	go func() {
		time.Sleep(100 * time.Millisecond)
		cancel()
	}()

	engine.StartWithContext(ctx)

	// After execution, status should be Complete (since no error)
	if info.Status != types.TexeStatus_Complete {
		t.Errorf("expected status Complete, got %v", info.Status)
	}
}
