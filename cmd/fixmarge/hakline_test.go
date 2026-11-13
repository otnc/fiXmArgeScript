package main

import (
	"errors"
	"io"
	"testing"

	"github.com/nyaosorg/go-readline-ny"
	"github.com/nyaosorg/go-readline-ny/keys"
	"github.com/nyaosorg/go-ttyadapter/auto"
)

// 端末の読み手を、打つキーの台本を渡して確かめる。
func TestEditedLines(t *testing.T) {
	typed := []string{
		"x", " ", "i", "z", "u", " ", "2", keys.Left, keys.Left, "\b", "\r", // 途中に戻って1文字消す
		"魔", "女", keys.Left, "\b", "\r", // 全角の文字も1文字ずつ動く
		" ", "\r", // 空白だけの行は履歴に入れない
		keys.Up, keys.Up, "\r", // 履歴を2つ戻る
		"a", keys.CtrlC, // 打ちかけを捨てる
		keys.CtrlD, // 入力の終わり
	}
	lines := newEditedLines(&readline.Editor{Tty: &auto.Pilot{Text: typed, Width: 80, Height: 24}, Writer: io.Discard})

	for _, want := range []string{"x iz 2\n", "女\n", " \n", "x iz 2\n"} {
		got, err := lines.readLine("> ")
		if err != nil || got != want {
			t.Fatalf("readLine = %q, %v, want %q", got, err, want)
		}
	}
	if _, err := lines.readLine("> "); !errors.Is(err, errInterrupted) {
		t.Fatalf("Ctrl+C: %v, want errInterrupted", err)
	}
	if line, err := lines.readLine("... "); line != "" || !errors.Is(err, io.EOF) {
		t.Fatalf("Ctrl+D: %q, %v, want io.EOF", line, err)
	}
}
