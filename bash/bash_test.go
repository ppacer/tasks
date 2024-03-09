package bash

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/ppacer/core/dag"
)

func TestBashExecLs(t *testing.T) {
	cmd := exec.Command("/bin/sh", "-c", "ls", "-la")
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

	testStdoutContains(llMaps, "bash_test.go", t)
}

func TestBashWriteIntoFile(t *testing.T) {
	const fileName = "tmp.tmp"
	cmd1 := exec.Command("/bin/sh", "-c", "echo 'hello' > "+fileName)
	cmd2 := exec.Command("/bin/sh", "-c", "ls", "-la")
	cmd3 := exec.Command("/bin/sh", "-c", "rm "+fileName)
	bashWrite := New("write", cmd1)
	bashList := New("list", cmd2)
	bashRemove := New("remove", cmd3)

	var logs bytes.Buffer
	tc := dag.TaskContext{Logger: jsonLoggerToBufio(&logs, slog.LevelInfo)}

	wErr := bashWrite.Execute(tc)
	if wErr != nil {
		t.Errorf("Executing bash cmd %s failed: %s", cmd1.String(),
			wErr.Error())
	}

	lsErr := bashList.Execute(tc)
	if lsErr != nil {
		t.Errorf("Executing bash cmd %s failed: %s", cmd2.String(),
			lsErr.Error())
	}
	fInfo, fErr := os.Stat(fileName)
	if os.IsNotExist(fErr) {
		t.Errorf("Expected file %s to exists, but it does not", fileName)
	}
	if fInfo.Size() == 0 {
		t.Error("Expected non-empty file")
	}

	rErr := bashRemove.Execute(tc)
	if rErr != nil {
		t.Errorf("Executing bash cmd %s failed: %s", cmd3.String(),
			rErr.Error())
	}

	// checking for stdout "tmp.tmp" from ls command
	llMaps, llErr := logLineMaps(&logs)
	if llErr != nil {
		t.Errorf("Cannot parse log lines: %s", llErr.Error())
	}
	testStdoutContains(llMaps, fileName, t)
}

func testStdoutContains(llMaps []map[string]any, phrase string, t *testing.T) {
	t.Helper()

	for _, logLineMap := range llMaps {
		stdoutMsg, ok := logLineMap["stdout"]
		if !ok {
			continue
		}
		stdoutMsgStr, castOk := stdoutMsg.(string)
		if !castOk {
			continue
		}
		if strings.Contains(stdoutMsgStr, phrase) {
			// found it!
			return
		}
	}
	t.Errorf("Phrase <%s> not found in stdout", phrase)
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
