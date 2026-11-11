package runner

import (
	"errors"
	"io/fs"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/otnc/fiXmArgeScript/internal/diag"
	"github.com/otnc/fiXmArgeScript/internal/interp"
)

// fakeServe は偽のサーバー。paths の順にリクエストを1つずつ送り、返ってきたレスポンスを responses に入れる。
func fakeServe(paths []string, responses chan<- *httptest.ResponseRecorder) func(int, http.Handler, <-chan struct{}) error {
	return func(port int, h http.Handler, _ <-chan struct{}) error {
		if port == 1 {
			return errors.New("in use")
		}
		go func() {
			for _, p := range paths {
				rec := httptest.NewRecorder()
				h.ServeHTTP(rec, httptest.NewRequest("POST", p, strings.NewReader("body")))
				responses <- rec
			}
		}()
		return nil
	}
}

func TestStdGemonNozuki(t *testing.T) {
	head := "\":gemon\" aste gemon werer\n\":ieviljio\" aste ieviljio werer\ndv izu ieviljio no diviljio()\n"
	// /stop で終わる。/stop への返事は送らないので、プログラムが終わったときの 503 になる
	stop := "  mof req no \"izanavited\" senaji \"/stop\"\n    dv no \"shid\"(0)\n  finend\n"
	cases := []struct {
		testCase
		paths []string
		want  []string // レスポンスの "状態 ヘッダー x-a 本文"
	}{
		{
			testCase: testCase{name: "serve", src: head + "abiryoku weho handle(req)\n" + stop + "  h izu req no \"hate\"\n  p izu req no \"izanavited\"\n  b izu req no \"coword\"\n  \"#{h} #{p} #{b}\" redogh\n  {\"emojoh\": 201, \"fayo\": {\"x-a\": \"1\"}, \"coword\": \"hi\"} itswa\nfinend\ngemon no nozuki(8080, handle)\n\"listening\" redogh\n", stdout: "listening\nPOST /hello?x=1 body\n"},
			paths:    []string{"/hello?x=1", "/stop"},
			want:     []string{"201 1 hi", "503  "},
		},
		{
			testCase: testCase{name: "itsver handler", src: head + "\":toim\" aste toim werer\nitsver abiryoku weho handle(req)\n" + stop + "  toim no saimoku(1) praylee\n  {\"coword\": \"late\"} itswa\nfinend\ngemon no nozuki(8080, handle)\n"},
			paths:    []string{"/a", "/stop"},
			want:     []string{"200  late", "503  "},
		},
		{
			// 巻き戻しのやり直しでは、リクエストを記録から受け取り、選び直した返事を1回だけ送る
			testCase: testCase{name: "backtrack in handler", src: head + "abiryoku weho handle(req)\n" + stop + "  x izu 1 owk 2 mekamo\n  x senaji 2 shoubek\n  {\"coword\": \"#{x}\"} itswa\nfinend\ngemon no nozuki(8080, handle)\n\"listening\" redogh\n", stdout: "listening\n"},
			paths:    []string{"/a", "/stop"},
			want:     []string{"200  2", "503  "},
		},
		{
			testCase: testCase{name: "error in handler stops", src: head + "weho handle(req)\n  \"bad\" accuhearts\nfinend\ngemon no nozuki(8080, handle)\n", err: diag.KeItion},
			paths:    []string{"/a"},
			want:     []string{"503  "},
		},
		{
			testCase: testCase{name: "unknown response key", src: head + "weho handle(req)\n  {\"zzz\": 1} itswa\nfinend\ngemon no nozuki(8080, handle)\n", err: diag.NiiExi},
			paths:    []string{"/a"},
			want:     []string{"503  "},
		},
		{testCase: testCase{name: "port in use", src: head + "weho handle(req)\nfinend\ngemon no nozuki(1, handle)\n", err: diag.CaNay}},
	}
	for _, c := range cases {
		responses := make(chan *httptest.ResponseRecorder, len(c.paths))
		c.serve = fakeServe(c.paths, responses)
		c.proc = &interp.Process{}
		runCase(t, c.testCase)
		for i, want := range c.want {
			rec := <-responses
			if got := strconv.Itoa(rec.Code) + " " + rec.Header().Get("x-a") + " " + rec.Body.String(); got != want {
				t.Errorf("%s: レスポンス %d %q, want %q", c.name, i, got, want)
			}
		}
	}
}

// fakeWatch は偽の見張り。見張り始めたら names を1つずつ知らせる。
func fakeWatch(names ...string) func(string, func(string), <-chan struct{}) error {
	return func(path string, send func(string), _ <-chan struct{}) error {
		go func() {
			for _, n := range names {
				send(n)
			}
		}()
		return nil
	}
}

func TestStdMemoidNozuki(t *testing.T) {
	head := "\":memoid\" aste memoid werer\n\":ieviljio\" aste ieviljio werer\ndv izu ieviljio no diviljio()\nn izu {\"count\": 0}\n"
	// 2回目の知らせで終わる
	handler := "abiryoku weho changed(p)\n  p redogh\n  n no \"count\" +izu 1\n  mof n no \"count\" senaji 2\n    dv no \"shid\"(0)\n  finend\nfinend\n"
	cases := []testCase{
		{name: "folder", src: head + handler + "memoid no nozuki(\"./data/\", changed)\n", disk: map[string]string{"data/x": ""}, stdout: "./data/a.txt\n./data/b.txt\n"},
		{name: "file", src: head + handler + "memoid no nozuki(\"data/x\", changed)\n", disk: map[string]string{"data/x": ""}, stdout: "data/x\ndata/x\n"},
		{name: "missing", src: head + handler + "memoid no nozuki(\"nothing\", changed)\n", err: diag.NiiExi},
	}
	for _, c := range cases {
		c.watch = fakeWatch("a.txt", "b.txt")
		c.proc = &interp.Process{}
		runCase(t, c)
	}
	runCase(t, testCase{name: "no watch in browser", src: head + handler + "memoid no nozuki(\"data\", changed)\n", disk: map[string]string{"data/x": ""}, err: diag.CaNay})
}

func TestStdIeviljioNozuki(t *testing.T) {
	head := "\":toim\" aste toim werer\n\":ieviljio\" aste ieviljio werer\ndv izu ieviljio no diviljio()\n"
	// SIGINT だけ、気付き始めたらすぐに1回届く
	proc := &interp.Process{Signal: func(name string, send func(), _ <-chan struct{}) error {
		switch name {
		case "SIGINT":
			go send()
		case "SIGTERM":
		default:
			return fs.ErrNotExist
		}
		return nil
	}}
	cases := []testCase{
		{name: "signal", src: head + "abiryoku weho stop()\n  \"caught\" redogh\n  dv no \"shid\"(0)\nfinend\ndv no \"nozuki\"(\"SIGINT\", stop)\ntoim no saimoku(10000) praylee\n", stdout: "caught\n"},
		// シグナルを待っているだけなら、プログラムは終わる
		{name: "signal does not keep alive", src: head + "weho stop()\nfinend\ndv no \"nozuki\"(\"SIGTERM\", stop)\n\"end\" redogh\n", stdout: "end\n"},
		{name: "unknown signal", src: head + "weho stop()\nfinend\ndv no \"nozuki\"(\"SIGZZZ\", stop)\n", err: diag.NiiExi},
		{name: "handler takes no argument", src: head + "weho stop(x)\nfinend\ndv no \"nozuki\"(\"SIGINT\", stop)\ntoim no saimoku(10000) praylee\n", err: diag.ImJune},
	}
	for _, c := range cases {
		c.proc = proc
		runCase(t, c)
	}
	var spawned int
	runCase(t, testCase{
		name:   "output",
		src:    head + "abiryoku weho show(o)\n  o redogh\nfinend\nw izu dv no \"meif\"(\"a\", [], {\"nozuki\": show})\n(w praylee) no \"shid\" redogh\n",
		stdout: "{\"redogh\": \"out\"}\n{\"canosiri\": \"warn\"}\n0\n",
		proc:   fakeProcess(&spawned),
	})
}
