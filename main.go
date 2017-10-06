package main

import (
	"bufio"
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
		log.Println("Responding the index")
		toc := `
<html><body><h1>chunked-app</h1>
<ul>
<li><a href="/buf">Buffered server</a>
<li><a href="/chunked">Streaming server</a>
<li><a href="/mix">Streaming server with content-length</a>
</ul>
</body></html>
`
		fmt.Fprint(w, toc)
	})

	h.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
	})

	h.HandleFunc("/buf", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Responding with buffered payload")
		code, _ := ioutil.ReadFile("main.go")
		fmt.Fprint(w, string(code))
	})

	h.HandleFunc("/chunked", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Responding with chunked payload")

		flusher, _ := w.(http.Flusher)

		file, _ := os.Open("main.go")
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			fmt.Fprintln(w, scanner.Text())
			flusher.Flush()
		}
	})

	// TODO: Serve this with both Transfer-Encoding: chunked and Content-Length
	h.HandleFunc("/mix", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Responding through IO stream with content-length")

		file, _ := os.Open("main.go")
		defer file.Close()

		stat, _ := file.Stat()
		w.Header().Set("Content-Length", strconv.FormatInt(stat.Size(), 10))
		io.Copy(w, file)
	})

	log.Println("Listening at port " + port)
	err := http.ListenAndServe(":"+port, h)
	log.Fatal(err)
}
