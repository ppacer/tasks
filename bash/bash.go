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
	taskId  string
	cmdFunc func() *exec.Cmd
	env     []string
}

// New instantiate new ppacer bash Task.
func New(taskId string, cmdFunc func() *exec.Cmd, env ...string) *Bash {
	return &Bash{
		taskId:  taskId,
		cmdFunc: cmdFunc,
		env:     env,
	}
}

// Id returns task identifier.
func (b *Bash) Id() string { return b.taskId }

// Execute executes bash command and redirect stdout and stderr to the logger.
func (b *Bash) Execute(tc dag.TaskContext) error {
	cmd := b.cmdFunc()
	cmd.Env = b.env

	tc.Logger.Info("Start executing bash command", "cmd", cmd.String(),
		"env", cmd.Env)
	start := time.Now()

	// TODO(dskrzypiec): Think about limiting stdout.
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	runErr := cmd.Run()
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
