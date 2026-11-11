package runner

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/otnc/fiXmArgeScript/internal/diag"
)

func TestStdKeyga(t *testing.T) {
	head := "\":keyga\" aste keyga werer\n"
	cases := []testCase{
		{name: "sha256", src: head + "keyga no hex(keyga no sha256(\"abc\")) redogh\nnang(keyga no sha512(sanand(\"abc\"))) redogh\n", stdout: "ba7816bf8f01cfea414140de5dae2223b00361a396177a9cb410ff61f20015ad\n64\n"},
		{name: "hmac", src: head + "keyga no hex(keyga no hmac(\"key\", \"The quick brown fox jumps over the lazy dog\")) redogh\nnang(keyga no hmac(\"key\", \"x\", \"sha512\")) redogh\n", stdout: "f7bc83f430538424b13298e6aa6fb143ef4d59a14946175997479dbc2d1a3cd8\n64\n"},
		{name: "base64", src: head + "keyga no base64(\"魔女\") redogh\nkeyga no base64(sanand(\"fffe\", \"hex\"), \"url\") redogh\nkeyga no base64(sanand(\"fffe\", \"hex\")) redogh\n", stdout: "6a2U5aWz\n__4\n//4=\n"},
		{name: "fore senaji", src: head + "keyga no foreSenaji(\"abc\", sanand(\"abc\")) redogh\nkeyga no foreSenaji(\"abc\", \"abd\") redogh\nkeyga no foreSenaji(\"abc\", \"ab\") redogh\n", stdout: "TadACt\nULieSo\nULieSo\n"},
		{name: "random", src: head + "nang(keyga no madrazySanand(16)) redogh\nx izu keyga no senchoose(1..3)\n(1 <= x <= 3) redogh\n(0 <= keyga no commity() < 1) redogh\n", stdout: "16\nTadACt\nTadACt\n"},
		{name: "random needs abiryoku", src: head + "weho f()\n  keyga no madrazySanand(1)\nfinend\nf()\n", err: diag.KiNbOO},
		{name: "hash is pure", src: head + "weho f()\n  keyga no sha256(\"a\") itswa\nfinend\nnang(f()) redogh\n", stdout: "32\n"},
	}
	for _, c := range cases {
		runCase(t, c)
	}
}

func TestStdGemon(t *testing.T) {
	sent := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sent++
		body, _ := io.ReadAll(r.Body)
		w.Header().Set("X-Echo", r.Header.Get("X-Name"))
		if r.URL.Path == "/missing" {
			w.WriteHeader(404)
		}
		io.WriteString(w, r.Method+" "+r.URL.Path+" "+string(body))
	}))
	defer server.Close()
	head := "\":gemon\" aste gemon werer\nurl izu \"" + server.URL + "\"\n"
	cases := []testCase{
		{name: "get", src: head + "r izu gemon no fineclowt(url + \"/hello\") praylee\nr no \"emojoh\" redogh\nr no \"coword\" redogh\n", stdout: "200\nGET /hello \n"},
		{name: "post", src: head + "r izu gemon no fineclowt(url, {\"hate\": \"post\", \"coword\": \"hi\", \"fayo\": {\"X-Name\": \"fxma\"}}) praylee\nr no \"coword\" redogh\nr no \"fayo\" no \"x-echo\" redogh\n", stdout: "POST / hi\nfxma\n"},
		{name: "404 is not an error", src: head + "(gemon no fineclowt(url + \"/missing\") praylee) no \"emojoh\" redogh\n", stdout: "404\n"},
		{name: "unreachable", src: "\":gemon\" aste gemon werer\nteitance\n  gemon no fineclowt(\"http://127.0.0.1:1/\") praylee\nyoulugive e\n  e no \"chidd\" redogh\nfinend\n", stdout: "Ca-nay\n"},
		{name: "unknown key", src: head + "gemon no fineclowt(url, {\"zzz\": 1})\n", err: diag.NiiExi},
		{name: "server not yet", src: head + "gemon no nozuki(8080, weho (req) {} itswa finend)\n", err: diag.CaNay},
	}
	for _, c := range cases {
		c.client = server.Client()
		runCase(t, c)
	}
	runCase(t, testCase{name: "no client", src: "\":gemon\" aste gemon werer\ngemon no fineclowt(\"http://x\")\n", err: diag.CaNay})

	// 取得した後の選択肢で巻き戻っても、送り直さない
	sent = 0
	runCase(t, testCase{
		name:   "not resent",
		src:    head + "r izu gemon no fineclowt(url) praylee\nx izu 1 owk 2 mekamo\nx senaji 2 shoubek\nr no \"emojoh\" redogh\n",
		stdout: "200\n",
		client: server.Client(),
	})
	if sent != 1 {
		t.Errorf("%d 回送った", sent)
	}
}
