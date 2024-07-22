package tasks

import (
	"testing"

	"github.com/ppacer/core/dag"
)

func TestEmptyTask(t *testing.T) {
	const taskId = "task1"
	var et dag.Task = NewEmpty(taskId)

	if et.Id() != taskId {
		t.Errorf("Expected taskId, to be %s, got: %s", taskId,
			et.Id())
	}
	execErr := et.Execute(dag.TaskContext{})
	if execErr != nil {
		t.Errorf("Error while executing Empty task: %s", execErr.Error())
	}
}
