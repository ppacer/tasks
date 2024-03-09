# ppacer tasks

This repository contains Go packages with generic and commonly used ppacer
`dag.Task` implementations.


## List of available `dag.Task` implementations

* `bash`
    * `Bash` - task for running Bash command and redirect stdout and stderr to
      the task logger.
* `aws`
    * ...


## Usage example

To add a dependency to your Go module, you can run (example for bash):

```
go get github.com/ppacer/tasks/bash@latest
```

Simplest DAG using bash task can look like the following:

```
package main

import (
        "os/exec"

        "github.com/ppacer/core/dag"
        "github.com/ppacer/tasks/bash"
)

func prepBashDagExample() dag.Dag {
        ls := bash.New("ls_tmp", exec.Command("ls", "/tmp"))
        cp := bash.New("copy_files", exec.Command("cp", "/tmp/f1", "/tmp/f2"))

        root := dag.NewNode(ls).Next(dag.NewNode(cp))

        lsDag := dag.New(dag.Id("ls_dag")).AddRoot(root).Done()
        return lsDag
}
```
