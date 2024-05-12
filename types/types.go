package types

type Queue interface {
	Enqueue(*TaskRunInfo) error
	Dequeue() *TaskRunInfo
}

type TexeOpts struct {
	MaxWorkers uint
	Queue      Queue
}

type TaskExe interface {
	TaskStartingCallback(*TaskRunInfo)
	Start(*TaskRunInfo) error
	TaskCompleteCallback(*TaskRunInfo, error)
}

type Task struct {
	Exe         TaskExe
	Description string
}

type TexeStatus int

const (
	TexeStatus_Error   TexeStatus = -1
	TexeStatus_Unknown TexeStatus = iota
	TexeStatus_Stopped
	TexeStatus_Queued
	TexeStatus_Running
	TexeStatus_Complete
)

func (s TexeStatus) String() string {
	return [...]string{"error", "unknown", "stopped", "queued", "running", "complete"}[s]
}

type TaskRunInfo struct {
	Task
	Status TexeStatus
	Error  error
}
