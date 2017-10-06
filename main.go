package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
)

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
<li><a href="/buf">Buffered server</a>
<li><a href="/stream">Streaming server</a>
<li><a href="/mix">Streaming server with content-length</a>
</ul>
</body></html>
`
		fmt.Fprint(w, toc)
	})

	h.HandleFunc("/buf", func(w http.ResponseWriter, r *http.Request) {
		code, _ := ioutil.ReadFile("main.go")
		fmt.Fprint(w, string(code))
	})

	h.HandleFunc("/stream", func(w http.ResponseWriter, r *http.Request) {
		file, _ := os.Open("main.go")
		defer file.Close()
		io.Copy(w, file)
	})

	// TODO: Serve this with both Transfer-Encoding: chunked and Content-Length
	h.HandleFunc("/mix", func(w http.ResponseWriter, r *http.Request) {
		file, _ := os.Open("main.go")
		defer file.Close()

		stat, _ := file.Stat()
		w.Header().Set("Content-Length", strconv.FormatInt(stat.Size(), 10))
		io.Copy(w, file)
	})

	err := http.ListenAndServe(":"+port, h)
	log.Fatal(err)
}
