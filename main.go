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

type extraHeaderResponseWriter struct {
	origWriter http.ResponseWriter
}

func (eh *extraHeaderResponseWriter) WriteHeader(rc int) {
	log.Println("Hey")
	eh.origWriter.WriteHeader(rc)
}

func (eh *extraHeaderResponseWriter) Write(b []byte) (int, error) {
	log.Println("Yo")
	return eh.origWriter.Write(b)
}

func (eh *extraHeaderResponseWriter) Header() http.Header {
	log.Println("Ya")
	return eh.origWriter.Header()
}


func extraHeaderHandler(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		w = &extraHeaderResponseWriter{origWriter: w}
		h.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

type throttlingHandler struct {
	d time.Duration
}

func (th *throttlingHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("Responding chunked lines with delays: %d\n", th.d)
	flusher, ok := w.(http.Flusher)
	if !ok {
		log.Fatal("Could not obtain Flusher")
	}

	file, _ := os.Open("main.go")
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fmt.Fprintln(w, scanner.Text())
		flusher.Flush()
		time.Sleep(th.d)
	}
}

type bufferedHandler struct {
}

func (*bufferedHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("Responding with buffered payload")
	code, _ := ioutil.ReadFile("main.go")
	fmt.Fprint(w, string(code))
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
	h.Handle("/buf", &bufferedHandler{})
	h.Handle("/chunked", &throttlingHandler{0 * time.Millisecond})
	h.Handle("/slow", &throttlingHandler{100 * time.Millisecond})
	h.HandleFunc("/mix", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Responding through IO stream with content-length")

		file, _ := os.Open("main.go")
		stat, _ := file.Stat()
		w.Header().Set("Content-Length", strconv.FormatInt(stat.Size(), 10))
		file.Close()
		w.Header().Set("Transfer-Encoding", "chunked")

		log.Println("Calling extraHeaderHandler")
		extraHeaderHandler(&bufferedHandler{}).ServeHTTP(w, r)
	})

	log.Println("Listening at port " + port)
	err := http.ListenAndServe(":"+port, h)
	log.Fatal(err)
}
