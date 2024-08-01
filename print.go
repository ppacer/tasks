package tasks

import (
	"fmt"
	"time"

	"github.com/ppacer/core/dag"
)

// PrintTask simply prints a message on stdout.
type PrintTask struct {
	taskId  string
	message string
}

// NewPrintTask initialize new PrintTask.
func NewPrintTask(taskId, message string) *PrintTask {
	return &PrintTask{
		taskId:  taskId,
		message: message,
	}
}

// Id return task identifier.
func (pt *PrintTask) Id() string { return pt.taskId }

// Execute prints message to stdout.
func (pt *PrintTask) Execute(tc dag.TaskContext) error {
	fmt.Printf("[%s] %s\n", pt.taskId, pt.message)
	tc.Logger.Info("PrintTask finished!", "ts", time.Now())
	return nil
}
