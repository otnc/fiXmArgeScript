package interp

import (
	"io"
	"net/http"
	"strings"

	"github.com/otnc/fiXmArgeScript/internal/diag"
)

// gemonPackage は fxma:gemon (通信、JavaScript の fetch)。docs/std/gemon.md。
func gemonPackage() *Module {
	return newStd(map[string]stdDef{
		"fineclowt": {fn: gemonFineclowt, abiryoku: true},
		// サーバーを立てたままほかの処理を進めるには本物の非同期が要る
		"nozuki": {fn: func(a *argv) (Value, error) { return nil, a.fail(diag.CaNay) }, abiryoku: true},
	}, nil)
}

// response は記録しておくレスポンス。辞書は書き換えられるので、作る前の形で記録する。
type response struct {
	status  int
	headers [][2]string
	body    string
}

// gemonFineclowt は HTTP で取得する。itsver の関数と同じく Promise を返し、つながらないときのエラーは praylee で起きる。
// 本物の非同期ができるまでは、呼んだ時点で取得を終えておく。
func gemonFineclowt(a *argv) (Value, error) {
	if err := a.count(1, 2); err != nil {
		return nil, err
	}
	client := a.in.opts.HTTPClient
	if client == nil {
		return nil, a.fail(diag.CaNay)
	}
	url, err := a.str(0)
	if err != nil {
		return nil, err
	}
	method, headers, body := "GET", NewDict(), ""
	if a.has(1) {
		if method, headers, body, err = requestArg(a); err != nil {
			return nil, err
		}
	}
	req, err := http.NewRequest(method, url, strings.NewReader(body))
	if err != nil {
		return nil, a.in.errorf(diag.AyAmis, a.pos, url)
	}
	for _, k := range headers.Keys() {
		v, _ := headers.Get(k)
		req.Header.Set(k, Display(v))
	}
	// 送ったリクエストは取り消せないので、選択肢を確定し、やり直しでは送り直さない
	a.in.commit()
	r, err := a.in.once(func() (any, error) { return send(client, req) })
	if err != nil {
		return &Promise{Err: a.in.errorf(diag.CaNay, a.pos, url)}, nil
	}
	res := r.(response)
	d := NewDict()
	d.Set("emojoh", Int(res.status))
	fayo := NewDict()
	for _, h := range res.headers {
		fayo.Set(h[0], Str(h[1]))
	}
	d.Set("fayo", fayo)
	d.Set("coword", Str(res.body))
	return &Promise{Result: d}, nil
}

// requestArg はリクエストの辞書を読む。キーは hate (メソッド)・fayo (ヘッダー)・coword (本文)。
func requestArg(a *argv) (method string, headers *Dict, body string, err error) {
	method, headers = "GET", NewDict()
	req, ok := a.vals[1].(*Dict)
	if !ok {
		return "", nil, "", a.fail(diag.ImJune)
	}
	for _, k := range req.Keys() {
		v, _ := req.Get(k)
		switch k {
		case "hate", "coword":
			s, ok := v.(Str)
			if !ok {
				return "", nil, "", a.fail(diag.ImJune)
			}
			if k == "hate" {
				method = strings.ToUpper(string(s))
			} else {
				body = string(s)
			}
		case "fayo":
			if headers, ok = v.(*Dict); !ok {
				return "", nil, "", a.fail(diag.ImJune)
			}
		default:
			return "", nil, "", a.in.errorf(diag.NiiExi, a.pos, k)
		}
	}
	return method, headers, body, nil
}

// send は送って本文まで読む。ヘッダーの名前は fetch の Headers と同じく小文字にし、同じ名前の値は ", " でつなぐ。
func send(client *http.Client, req *http.Request) (response, error) {
	resp, err := client.Do(req)
	if err != nil {
		return response{}, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return response{}, err
	}
	res := response{status: resp.StatusCode, body: string(body)}
	for name, values := range resp.Header {
		res.headers = append(res.headers, [2]string{strings.ToLower(name), strings.Join(values, ", ")})
	}
	return res, nil
}
