package tasks

import "github.com/ppacer/core/dag"

// Empty represents empty Task which does nothing.
type Empty struct {
	taskId string
}

// NewEmpty initialize new Empty Task.
func NewEmpty(taskId string) Empty {
	return Empty{taskId: taskId}
}

// Id returns Task identifier.
func (e Empty) Id() string { return e.taskId }

// Execute does nothing and immediately finishes with success.
func (e Empty) Execute(_ dag.TaskContext) error {
	return nil
}
