package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
)

// serve は gemon no nozuki のサーバーを立てる。受け付けを始めたら戻り、done が閉じたら止める。
func serve(port int, handler http.Handler, done <-chan struct{}) error {
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}
	// net/http のログは英語なので出さない (コンソールにはフィクスマージ語だけを出す)
	srv := &http.Server{Handler: handler, ErrorLog: log.New(io.Discard, "", 0)}
	go srv.Serve(ln)
	go func() {
		<-done
		srv.Close()
	}()
	return nil
}
