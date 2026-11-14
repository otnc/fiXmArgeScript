package main

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestReplaceExecutable(t *testing.T) {
	dir := t.TempDir()
	exe := filepath.Join(dir, "fixmarge.exe")
	next := filepath.Join(dir, "download")
	if err := os.WriteFile(exe, []byte("old"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(next, []byte("new"), 0o644); err != nil {
		t.Fatal(err)
	}
	// 前の更新で残った .old は消してから退かす
	if err := os.WriteFile(exe+".old", []byte("older"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := replaceExecutable(exe, next); err != nil {
		t.Fatal(err)
	}
	if data, _ := os.ReadFile(exe); string(data) != "new" {
		t.Errorf("exe = %q, want new", data)
	}
	if info, _ := os.Stat(exe); runtime.GOOS != "windows" && info.Mode().Perm() != 0o755 {
		t.Errorf("mode = %v, want the old exe's 0755", info.Mode().Perm())
	}
	if _, err := os.Stat(exe + ".new"); !os.IsNotExist(err) {
		t.Errorf(".new is left: %v", err)
	}
	old, err := os.ReadFile(exe + ".old")
	if runtime.GOOS == "windows" && string(old) != "old" {
		t.Errorf(".old = %q, %v; want the replaced exe", old, err)
	}
	if runtime.GOOS != "windows" && !os.IsNotExist(err) {
		t.Errorf(".old should be removed: %v", err)
	}
}
