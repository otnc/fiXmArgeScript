//go:build js && wasm

// fixmarge-wasm は WebAssembly 版の入口。JavaScript に run・build・words を出す。npm の fixmarge パッケージ (npm/) がこれを包む。
// ブラウザでもNode.jsでも動くので、ファイルはメモリ上に置き、通信はホストの fetch を使う (Go の net/http が GOOS=js でそうする)。
package main

import (
	"errors"
	"io/fs"
	"net/http"
	"path"
	"strings"
	"syscall/js"

	"github.com/spf13/afero"

	"github.com/otnc/fiXmArgeScript/internal/build"
	"github.com/otnc/fiXmArgeScript/internal/interp"
	"github.com/otnc/fiXmArgeScript/internal/lexicon"
	"github.com/otnc/fiXmArgeScript/internal/runner"
)

// mainFile は run と build に渡したソースのファイル名。エラーの場所に出る。
const mainFile = "main.fxma"

func main() {
	api := js.Global().Get("Object").New()
	api.Set("run", js.FuncOf(run))
	api.Set("build", js.FuncOf(buildSource))
	api.Set("words", js.FuncOf(words))
	js.Global().Set("__fxma", api)
	// npm/index.js が、読み込み終わるのを待っている
	if ready := js.Global().Get("__fxmaReady"); ready.Type() == js.TypeFunction {
		ready.Invoke()
	}
	select {}
}

// promise は work を別の goroutine で動かし、その結果で解決する Promise を返す。
// JavaScript から呼ばれた関数の中で待つと、Go の実行が止まってしまうため。
func promise(work func() map[string]any) js.Value {
	return js.Global().Get("Promise").New(js.FuncOf(func(_ js.Value, args []js.Value) any {
		resolve := args[0]
		go func() { resolve.Invoke(work()) }()
		return nil
	}))
}

// run(source, { stdin, files, args, maxSteps }) は下書きを実行し、{ stdout, stderr, exitCode, files } で解決する。
// files は名前 → 中身の辞書で、werer でも memoid でも読める。実行した後のファイルも files で返す。
func run(_ js.Value, args []js.Value) any {
	source := args[0].String()
	opts := optionsArg(args)
	maxSteps := 0
	if n := opts.Get("maxSteps"); n.Type() == js.TypeNumber {
		maxSteps = n.Int()
	}
	return promise(func() map[string]any {
		fsys := afero.NewMemMapFs()
		for name, content := range stringMap(opts.Get("files")) {
			afero.WriteFile(fsys, path.Join("/", name), []byte(content), 0o644)
		}
		// werer の名前は main.fxma からの相対なので、memoid と同じく / から数える
		readFile := func(name string) ([]byte, error) { return afero.ReadFile(fsys, path.Join("/", name)) }
		var stdout, stderr strings.Builder
		o := interp.Options{
			Stdin:      strings.NewReader(stringOr(opts.Get("stdin"), "")),
			Stdout:     &stdout,
			Stderr:     &stderr,
			FS:         fsys,
			Process:    &interp.Process{Args: stringList(opts.Get("args")), Platform: platform()},
			HTTPClient: http.DefaultClient,
			MaxSteps:   maxSteps,
		}
		warnings, err := runner.Run(source, mainFile, readFile, o)
		for _, w := range warnings {
			stderr.WriteString(w.String() + "\n")
		}
		code := 0
		if exit, ok := errors.AsType[*interp.Exit](err); ok {
			code = exit.Code
		} else if err != nil {
			stderr.WriteString(err.Error() + "\n")
			code = 1
		}
		return map[string]any{"stdout": stdout.String(), "stderr": stderr.String(), "exitCode": code, "files": readFiles(fsys)}
	})
}

// build(source) は清書を作り、{ output, stderr } で解決する。エラーなら output は null。
func buildSource(_ js.Value, args []js.Value) any {
	source := args[0].String()
	return promise(func() map[string]any {
		out, warnings, err := build.Build(source, mainFile)
		var stderr strings.Builder
		for _, w := range warnings {
			stderr.WriteString(w.String() + "\n")
		}
		var output any = out
		if err != nil {
			stderr.WriteString(err.Error() + "\n")
			output = nil
		}
		return map[string]any{"output": output, "stderr": stderr.String()}
	})
}

// uses は words() に出す語の使い道。エディタの色分けに使うものだけ。
var uses = []struct {
	use  lexicon.Use
	name string
}{
	{lexicon.Keyword, "keyword"},
	{lexicon.Builtin, "builtin"},
	{lexicon.Package, "package"},
	{lexicon.Function, "function"},
	{lexicon.Key, "key"},
	{lexicon.Error, "error"},
}

// words() は下書きに書ける名前の一覧 [{ name, spelling, uses }] を返す。エディタの色分けに使う。
// キーワードを JavaScript 側にもう一度書かないよう、語の表 (lexicon) から作る。
func words(js.Value, []js.Value) any {
	var list []any
	for _, e := range lexicon.Entries() {
		var names []any
		for _, u := range uses {
			if e.Uses&u.use != 0 {
				names = append(names, u.name)
			}
		}
		list = append(list, map[string]any{"name": e.Name, "spelling": e.Spelling, "uses": names})
	}
	return list
}

// platform は dv no "wid" の値。docs/std/ieviljio.md の表のとおり、JavaScript の環境を見て決める。
func platform() string {
	g := js.Global()
	switch {
	case g.Get("document").Truthy():
		return "browser"
	case g.Get("Deno").Truthy():
		return "deno"
	case g.Get("Bun").Truthy():
		return "bun"
	}
	return "node"
}

func optionsArg(args []js.Value) js.Value {
	if len(args) > 1 && args[1].Type() == js.TypeObject {
		return args[1]
	}
	return js.Global().Get("Object").New()
}

func stringOr(v js.Value, fallback string) string {
	if v.Type() != js.TypeString {
		return fallback
	}
	return v.String()
}

func stringMap(v js.Value) map[string]string {
	m := map[string]string{}
	if v.Type() != js.TypeObject {
		return m
	}
	keys := js.Global().Get("Object").Call("keys", v)
	for i := range keys.Length() {
		k := keys.Index(i).String()
		m[k] = stringOr(v.Get(k), "")
	}
	return m
}

func stringList(v js.Value) []string {
	var list []string
	if !js.Global().Get("Array").Call("isArray", v).Bool() {
		return list
	}
	for i := range v.Length() {
		list = append(list, stringOr(v.Index(i), ""))
	}
	return list
}

// readFiles は実行した後のファイルを、名前 (/ を除いたもの) → 中身の辞書にする。
func readFiles(fsys afero.Fs) map[string]any {
	files := map[string]any{}
	afero.Walk(fsys, "/", func(p string, info fs.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}
		content, err := afero.ReadFile(fsys, p)
		if err == nil {
			files[strings.TrimPrefix(p, "/")] = string(content)
		}
		return nil
	})
	return files
}
