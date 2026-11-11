package interp

import (
	"io"
	"net/http"
	"slices"
	"strings"

	"github.com/otnc/fiXmArgeScript/internal/diag"
)

// gemonPackage は :gemon (通信、JavaScript の fetch)。docs/std/gemon.md。
func gemonPackage() *Module {
	return newStd(map[string]stdDef{
		"fineclowt": {fn: gemonFineclowt, abiryoku: true, itsver: true},
		"nozuki":    {fn: gemonNozuki, abiryoku: true},
	}, nil)
}

// response は記録しておくレスポンス。辞書は書き換えられるので、作る前の形で記録する。
type response struct {
	status  int
	headers [][2]string
	body    string
}

// gemonFineclowt は HTTP で取得する。itsver の関数と同じく Promise を返し、つながらないときのエラーは praylee で起きる。
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
	// 送ったリクエストは取り消せないので、startOp が選択肢を確定し、やり直しでは送り直さない
	return a.in.startOp(func() (any, error) { return send(client, req) }, func(r any, err error) (Value, error) {
		if err != nil {
			return nil, a.in.errorf(diag.CaNay, a.pos, url)
		}
		return responseDict(r.(response)), nil
	}), nil
}

// responseDict は fineclowt の結果の辞書。キーは emojoh (状態コード)・fayo (ヘッダー)・coword (本文)。
func responseDict(res response) *Dict {
	d := NewDict()
	d.Set("emojoh", Int(res.status))
	fayo := NewDict()
	for _, h := range res.headers {
		fayo.Set(h[0], Str(h[1]))
	}
	d.Set("fayo", fayo)
	d.Set("coword", Str(res.body))
	return d
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
	return response{status: resp.StatusCode, headers: headerPairs(resp.Header), body: string(body)}, nil
}

// headerPairs はヘッダーを記録できる形にする。名前は fetch の Headers と同じく小文字にし、同じ名前の値は ", " でつなぐ。
func headerPairs(h http.Header) [][2]string {
	var pairs [][2]string
	for name, values := range h {
		pairs = append(pairs, [2]string{strings.ToLower(name), strings.Join(values, ", ")})
	}
	// map の順は決まらないので、名前の順にそろえる
	slices.SortFunc(pairs, func(a, b [2]string) int { return strings.Compare(a[0], b[0]) })
	return pairs
}

// serverRequest はサーバーが受け取ったリクエストの知らせ。返事は reply に1回だけ送る。
type serverRequest struct {
	method, path string
	headers      [][2]string
	body         string
	reply        chan response
}

// gemonNozuki はサーバーを立て、リクエストに気付くたびに f(req) を呼んで、返した辞書で応答する。
// サーバーはプログラムが終わるまで動き続ける (Node.js と同じく、プログラムはサーバーがある限り終わらない)。
func gemonNozuki(a *argv) (Value, error) {
	if err := a.count(2, 2); err != nil {
		return nil, err
	}
	serve := a.in.opts.Serve
	if serve == nil {
		return nil, a.fail(diag.CaNay)
	}
	port, err := a.index(0)
	if err != nil {
		return nil, err
	}
	if port < 0 || port > 65535 {
		return nil, a.in.errorf(diag.ImJune, a.pos, Display(a.vals[0]))
	}
	f, err := a.fn(1)
	if err != nil {
		return nil, err
	}
	call := a.callbackLater(f)
	err = a.in.listen(true, func(send func(any), done <-chan struct{}) error {
		return serve(int(port), serveHandler(send, done), done)
	}, func(v any) {
		req := v.(serverRequest)
		call([]Value{serverRequestDict(req)}, func(result Value) error {
			res, err := responseArg(a, result)
			if err != nil {
				return err
			}
			// 返事は取り消せないので選択肢を確定し、やり直しでは送り直さない
			a.in.commit()
			a.in.once(func() (any, error) {
				req.reply <- res
				return nil, nil
			})
			return nil
		})
	})
	if err != nil {
		return nil, a.in.errorf(diag.CaNay, a.pos, Display(a.vals[0]))
	}
	return null, nil
}

// serveHandler はリクエストを知らせにして送り、返事を待って応答する。プログラムが終わったら 503 で応答する。
func serveHandler(send func(any), done <-chan struct{}) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		req := serverRequest{method: r.Method, path: r.URL.RequestURI(), headers: headerPairs(r.Header), body: string(body), reply: make(chan response, 1)}
		send(req)
		select {
		case res := <-req.reply:
			for _, h := range res.headers {
				w.Header().Set(h[0], h[1])
			}
			w.WriteHeader(res.status)
			io.WriteString(w, res.body)
		case <-done:
			w.WriteHeader(http.StatusServiceUnavailable)
		}
	})
}

// serverRequestDict は f に渡すリクエストの辞書。キーは hate・izanavited・fayo・coword。
func serverRequestDict(req serverRequest) *Dict {
	d := NewDict()
	d.Set("hate", Str(req.method))
	d.Set("izanavited", Str(req.path))
	fayo := NewDict()
	for _, h := range req.headers {
		fayo.Set(h[0], Str(h[1]))
	}
	d.Set("fayo", fayo)
	d.Set("coword", Str(req.body))
	return d
}

// responseArg は f が返したレスポンスの辞書を読む。キーは emojoh (省くと 200)・fayo・coword。
func responseArg(a *argv, v Value) (response, error) {
	res := response{status: http.StatusOK}
	d, ok := v.(*Dict)
	if !ok {
		return res, a.fail(diag.ImJune)
	}
	for _, k := range d.Keys() {
		v, _ := d.Get(k)
		switch k {
		case "emojoh":
			n, ok := toIndex(v)
			if !ok || n < 100 || n > 999 {
				return res, a.in.errorf(diag.ImJune, a.pos, Display(v))
			}
			res.status = int(n)
		case "fayo":
			fayo, ok := v.(*Dict)
			if !ok {
				return res, a.fail(diag.ImJune)
			}
			for _, name := range fayo.Keys() {
				val, _ := fayo.Get(name)
				res.headers = append(res.headers, [2]string{name, Display(val)})
			}
		case "coword":
			s, ok := v.(Str)
			if !ok {
				return res, a.fail(diag.ImJune)
			}
			res.body = string(s)
		default:
			return res, a.in.errorf(diag.NiiExi, a.pos, k)
		}
	}
	return res, nil
}
