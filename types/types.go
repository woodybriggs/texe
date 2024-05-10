package types

type Queue interface {
	Enqueue(*TaskContext) error
	Dequeue() *TaskContext
}

type TexeOpts struct {
	MaxWorkers uint
	Queue      Queue
}

type TaskExe interface {
	Start(*TaskContext) error
}

type TaskDef struct {
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
	return [...]string{"error", "stopped", "queued", "running", "complete"}[s]
}

type TaskContext struct {
	Def    *TaskDef
	Status TexeStatus
	Error  error
}
