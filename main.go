package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

type throttlingHandler struct {
	d time.Duration
}

func (th *throttlingHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("Responding cunked lines with delays: %d\n", th.d)
	flusher, _ := w.(http.Flusher)

	file, _ := os.Open("main.go")
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fmt.Fprintln(w, scanner.Text())
		flusher.Flush()
		time.Sleep(th.d)
	}
}

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
	<ul>
	<li><a href="/slow">slowly</a>
	</ul>
<li><a href="/mix">Streaming server with content-length</a>
</ul>
</body></html>
`
		fmt.Fprint(w, toc)
	})

	h.HandleFunc("/favicon.ico", http.NotFound)

	h.HandleFunc("/buf", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Responding with buffered payload")
		code, _ := ioutil.ReadFile("main.go")
		fmt.Fprint(w, string(code))
	})

	h.Handle("/chunked", &throttlingHandler{0 * time.Millisecond})

	h.Handle("/slow", &throttlingHandler{100 * time.Millisecond})

	// TODO: Serve this with both Transfer-Encoding: chunked and Content-Length
	h.HandleFunc("/mix", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Responding through IO stream with content-length")

		file, _ := os.Open("main.go")
		stat, _ := file.Stat()
		w.Header().Set("Content-Length", strconv.FormatInt(stat.Size(), 10))
		file.Close()
		w.Header().Set("Transfer-Encoding", "chunked")	// This erases Content-Length

		log.Println("Calling throttlingHandler")
		th := &throttlingHandler{0 * time.Millisecond}
		th.ServeHTTP(w, r)
	})

	log.Println("Listening at port " + port)
	err := http.ListenAndServe(":"+port, h)
	log.Fatal(err)
}
