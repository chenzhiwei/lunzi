package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

func rootHandler(w http.ResponseWriter, r *http.Request) {
	host, _ := os.Hostname()
	resp := ""

	verbose := os.Getenv("WEBSERVER_VERBOSE")
	if verbose == "true" {
		body, _ := ioutil.ReadAll(r.Body)
		resp = fmt.Sprintf("time=%v\nhost=%s\nmethod=%s\nuri=%s\nprotocol=%s\nheader=%v\nbody=%s\n", time.Now().Format("20060102T150405Z0700"), host, r.Method, r.RequestURI, r.Proto, r.Header, body)
	} else {
		resp = fmt.Sprintf("time=%v, host=%s, uri=%s\n", time.Now().Format("20060102T150405Z0700"), host, r.RequestURI)
	}

	fmt.Println(resp)
	io.WriteString(w, resp)
}

func main() {
	addr := os.Getenv("WEBSERVER_ADDR")
	if addr == "" {
		addr = ":8080"
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", rootHandler)
	http.ListenAndServe(addr, mux)
}
