package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/mattn/go-colorable"
	"github.com/nyaosorg/go-readline-ny"
	"github.com/nyaosorg/go-readline-ny/simplehistory"
	"golang.org/x/term"
)

// errInterrupted は、打ちかけの行を Ctrl+C で捨てたとき。
var errInterrupted = errors.New("interrupted")

// lineReader は hak の入力を1行読む。行は改行で終わる。入力の終わりは io.EOF で、
// 改行のない最後の行は、行と io.EOF を一緒に返す。
type lineReader interface {
	readLine(prompt string) (string, error)
}

// newLineReader は、端末なら矢印キーと履歴の使える読み手を、パイプやファイルならそのまま読む読み手を返す。
func (s streams) newLineReader(stdin *bufio.Reader) lineReader {
	in, inIsFile := s.stdin.(*os.File)
	out, outIsFile := s.stdout.(*os.File)
	if inIsFile && outIsFile && term.IsTerminal(int(in.Fd())) && term.IsTerminal(int(out.Fd())) {
		// Windows の古いコンソールでもエスケープシーケンスが効くように colorable を通す
		return newEditedLines(&readline.Editor{Writer: colorable.NewColorable(out)})
	}
	return pipedLines{in: stdin, out: s.stdout}
}

// pipedLines は端末でない入力をそのまま読む。
type pipedLines struct {
	in  *bufio.Reader
	out io.Writer
}

func (p pipedLines) readLine(prompt string) (string, error) {
	fmt.Fprint(p.out, prompt)
	line, err := p.in.ReadString('\n')
	if err != nil {
		// 端末なら最後に打った Enter で改行されるが、ここでは誰も改行しない
		if line == "" {
			fmt.Fprintln(p.out)
		}
		return line, io.EOF
	}
	return line, nil
}

// editedLines は端末から読む。行の編集と、上下の矢印キーでの履歴 (終わるまでの間だけ) が使える。
type editedLines struct {
	editor  *readline.Editor
	history *simplehistory.Container
}

func newEditedLines(editor *readline.Editor) editedLines {
	history := simplehistory.New()
	editor.History = history
	return editedLines{editor: editor, history: history}
}

func (e editedLines) readLine(prompt string) (string, error) {
	e.editor.PromptWriter = func(w io.Writer) (int, error) { return io.WriteString(w, prompt) }
	line, err := e.editor.ReadLine(context.Background())
	if errors.Is(err, readline.CtrlC) {
		return "", errInterrupted
	}
	if err != nil {
		return "", io.EOF
	}
	if strings.TrimSpace(line) != "" {
		e.history.Add(line)
	}
	return line + "\n", nil
}
