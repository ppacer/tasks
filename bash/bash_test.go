package bash

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"os/exec"
	"strings"
	"testing"

	"github.com/ppacer/core/dag"
)

func TestBashExecLs(t *testing.T) {
	cmd := exec.Command("ls", "-la")
	bt := New("mock_task", cmd)
	var logs bytes.Buffer
	tc := dag.TaskContext{Logger: jsonLoggerToBufio(&logs, slog.LevelInfo)}

	execErr := bt.Execute(tc)
	if execErr != nil {
		t.Errorf("Executing bash cmd %s failed: %s", cmd.String(),
			execErr.Error())
	}

	llMaps, llErr := logLineMaps(&logs)
	if llErr != nil {
		t.Errorf("Cannot parse log lines: %s", llErr.Error())
	}
	if len(llMaps) != 3 {
		t.Errorf("Expected 3 log lines, got: %d", len(llMaps))
	}

	// Expecting stdout with "bash_test.go" file
	stdoutMsg, ok := llMaps[2]["stdout"]
	if !ok {
		t.Errorf("Expected <msg> key to be in JSON log line: %v", llMaps[2])
	}
	stdoutMsgStr, castOk := stdoutMsg.(string)
	if !castOk {
		t.Errorf("Expected value for key 'msg', to be string, but it's not: %+v",
			stdoutMsg)
	}
	if !strings.Contains(stdoutMsgStr, "bash_test.go") {
		t.Errorf("Expected <bash_test.go>, to be listed by <ls -la> command. Got: %s",
			stdoutMsgStr)
	}
}

func jsonLoggerToBufio(b *bytes.Buffer, lvl slog.Level) *slog.Logger {
	h := slog.HandlerOptions{Level: lvl}
	return slog.New(slog.NewJSONHandler(b, &h))
}

func logLineMaps(b *bytes.Buffer) ([]map[string]any, error) {
	lineMaps := make([]map[string]any, 0)
	logLines := bytes.Split(b.Bytes(), []byte{'\n'})

	for _, ll := range logLines {
		if len(ll) > 0 {
			var m map[string]any
			jErr := json.Unmarshal(ll, &m)
			if jErr != nil {
				return lineMaps, jErr
			}
			lineMaps = append(lineMaps, m)
		}
	}

	return lineMaps, nil
}
