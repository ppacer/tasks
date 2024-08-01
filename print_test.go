package tasks

import "testing"

func TestNewPrintTask(t *testing.T) {
	data := []struct {
		taskId  string
		message string
	}{
		{"t1", "msg1"},
		{"t2", ""},
	}

	for _, input := range data {
		pt := NewPrintTask(input.taskId, input.message)
		if pt.Id() != input.taskId {
			t.Errorf("Expected task ID %s, got: %s", input.taskId, pt.Id())
		}
		if pt.message != input.message {
			t.Errorf("Expected message [%s], got [%s]", input.message, pt.message)
		}
	}
}
