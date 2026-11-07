package main

import (
	"bytes"
	"errors"
	"io/fs"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/otnc/fiXmArgeScript/internal/interp"
)

// osProcess は fxma:ieviljio に渡す、このプロセスの情報と魔女 (子プロセス) を生む仕組み。
func osProcess(args []string) *interp.Process {
	return &interp.Process{
		Args:     args,
		Env:      os.Environ(),
		Pid:      os.Getpid(),
		Platform: runtime.GOOS,
		Version:  moduleVersion(),
		Spawn:    spawn,
	}
}

type osChild struct {
	cmd            *exec.Cmd
	stdout, stderr bytes.Buffer
}

func spawn(c interp.Command) (interp.Child, error) {
	path, err := exec.LookPath(c.Name)
	if errors.Is(err, exec.ErrNotFound) {
		return nil, fs.ErrNotExist
	}
	if err != nil {
		return nil, err
	}
	child := &osChild{cmd: exec.Command(path, c.Args...)}
	child.cmd.Dir = c.Dir
	child.cmd.Env = c.Env
	child.cmd.Stdin = strings.NewReader(c.Stdin)
	child.cmd.Stdout = &child.stdout
	child.cmd.Stderr = &child.stderr
	if err := child.cmd.Start(); err != nil {
		return nil, err
	}
	return child, nil
}

func (c *osChild) Pid() int { return c.cmd.Process.Pid }

// Wait は終わるまで待つ。0 以外の終了コードはエラーではなく結果として返す (Node.js の spawn と同じ)。
func (c *osChild) Wait() (interp.ChildResult, error) {
	err := c.cmd.Wait()
	if _, ok := errors.AsType[*exec.ExitError](err); ok {
		err = nil
	}
	return interp.ChildResult{Stdout: c.stdout.String(), Stderr: c.stderr.String(), Code: c.cmd.ProcessState.ExitCode()}, err
}

func (c *osChild) Kill() error { return c.cmd.Process.Kill() }
