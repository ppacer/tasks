package bash

import (
	"bytes"
	"os/exec"
	"time"

	"github.com/ppacer/core/dag"
)

// Bash represents ppacer Task for executing bash commands. Both stdout and
// stderr streams are redirect to the task logger.
type Bash struct {
	taskId string
	cmd    *exec.Cmd
}

// NewBash instantiate new ppacer bash Task.
func NewBash(taskId string, bashCmd *exec.Cmd) *Bash {
	return &Bash{
		taskId: taskId,
		cmd:    bashCmd,
	}
}

// Execute executes bash command and redirect stdout and stderr to the logger.
func (b *Bash) Execute(tc dag.TaskContext) error {
	tc.Logger.Info("Start executing bash command", "cmd", b.cmd.String(),
		"env", b.cmd.Env)
	start := time.Now()

	// TODO(dskrzypiec): Think about limiting stdout.
	var stdout, stderr bytes.Buffer
	b.cmd.Stdout = &stdout
	b.cmd.Stderr = &stderr

	runErr := b.cmd.Run()
	if runErr != nil {
		tc.Logger.Error("Bash command has failed", "duration",
			time.Since(start), "stderr", stderr.String())
		if stdout.Len() > 0 {
			tc.Logger.Warn("Stdout content", "stdout", stdout.String())
		}
		return runErr
	}
	tc.Logger.Info("Successfully finished executing bash command", "durationMs",
		time.Since(start).Milliseconds())
	tc.Logger.Info("Standard output", "stdout", stdout.String())

	return nil
}
