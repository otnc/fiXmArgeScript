package main

import (
	"bytes"
	"errors"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"slices"
	"strings"
	"sync"
	"syscall"
	"unicode/utf8"

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
		Signal:   notifySignal,
	}
}

type osChild struct {
	cmd            *exec.Cmd
	stdout, stderr bytes.Buffer
	forwarding     sync.WaitGroup // OnOutput に出力を渡している goroutine
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
	if c.OnOutput == nil {
		child.cmd.Stdout = &child.stdout
		child.cmd.Stderr = &child.stderr
		if err := child.cmd.Start(); err != nil {
			return nil, err
		}
		return child, nil
	}
	stdout, err := child.cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	stderr, err := child.cmd.StderrPipe()
	if err != nil {
		return nil, err
	}
	if err := child.cmd.Start(); err != nil {
		return nil, err
	}
	child.forwarding.Go(func() { forward(stdout, &child.stdout, func(s string) { c.OnOutput(false, s) }) })
	child.forwarding.Go(func() { forward(stderr, &child.stderr, func(s string) { c.OnOutput(true, s) }) })
	return child, nil
}

// forward は出力を読みながら on に渡し、まとめて buf にも入れる。
// 文字の途中で読み終わったら、その文字は次に読んだものと一緒に渡す (UTF-8 の文字を割らないため)。
func forward(r io.Reader, buf *bytes.Buffer, on func(string)) {
	chunk := make([]byte, 4096)
	var rest []byte
	for {
		n, err := r.Read(chunk)
		buf.Write(chunk[:n])
		data := append(rest, chunk[:n]...)
		cut := completeRunes(data)
		if err != nil {
			cut = len(data)
		}
		if cut > 0 {
			on(string(data[:cut]))
		}
		rest = slices.Clone(data[cut:])
		if err != nil {
			return
		}
	}
}

// completeRunes は、data の頭から、最後まで揃った文字だけの長さを返す。
func completeRunes(data []byte) int {
	for i := len(data) - 1; i >= 0 && i >= len(data)-utf8.UTFMax; i-- {
		if utf8.RuneStart(data[i]) {
			if utf8.FullRune(data[i:]) {
				return len(data)
			}
			return i
		}
	}
	return len(data)
}

func (c *osChild) Pid() int { return c.cmd.Process.Pid }

// Wait は終わるまで待つ。0 以外の終了コードはエラーではなく結果として返す (Node.js の spawn と同じ)。
func (c *osChild) Wait() (interp.ChildResult, error) {
	// os/exec は、パイプを読み終える前に Wait してはいけない
	c.forwarding.Wait()
	err := c.cmd.Wait()
	if _, ok := errors.AsType[*exec.ExitError](err); ok {
		err = nil
	}
	return interp.ChildResult{Stdout: c.stdout.String(), Stderr: c.stderr.String(), Code: c.cmd.ProcessState.ExitCode()}, err
}

func (c *osChild) Kill() error { return c.cmd.Process.Kill() }

// signals は dv no "nozuki" で使えるシグナル。Node.js と同じ名前で書く。
var signals = map[string]os.Signal{
	"SIGINT":  os.Interrupt,
	"SIGTERM": syscall.SIGTERM,
	"SIGHUP":  syscall.SIGHUP,
	"SIGQUIT": syscall.SIGQUIT,
}

// notifySignal はシグナルに気付くたびに send を呼ぶ。気付いている間は、Ctrl+C などで終了しなくなる (Node.js と同じ)。
func notifySignal(name string, send func(), done <-chan struct{}) error {
	sig, ok := signals[name]
	if !ok {
		return fs.ErrNotExist
	}
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, sig)
	go func() {
		defer signal.Stop(ch)
		for {
			select {
			case <-ch:
				send()
			case <-done:
				return
			}
		}
	}()
	return nil
}
