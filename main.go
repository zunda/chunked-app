package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"io"
)

type sourceCodeServer struct{}

// This seems to respond with Transfer-Encoding: chunked
func (sourceCodeServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	file, err := os.Open("main.go")
	if err != nil {
		w.WriteHeader(503)
		fmt.Fprint(w, err.Error() + "\n")
		return
	}
	_, err = io.Copy(w, file)
	if err != nil {
		w.WriteHeader(503)
		fmt.Fprint(w, err.Error() + "\n")
		return
	}
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	h := http.NewServeMux()
	h.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		toc := `
<html><body><h1>chunked-app</h1>
<ul>
<li><a href="/">/</a> : without <tt>Transfer-Encoding: chunked</tt> and with <tt>Content-length</tt>
<li><a href="/code">server source code</a> : with <tt>Transfer-Encoding: chunked</tt> and without <tt>Content-length</tt>
</ul>
</body></html>
`
		fmt.Fprint(w, toc)
	})
	h.Handle("/code", &sourceCodeServer{})

	err := http.ListenAndServe(":" + port, h)
	log.Fatal(err)
}
