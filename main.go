package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"time"
)

func handler(w http.ResponseWriter, r *http.Request) {
	dump, err := httputil.DumpRequest(r, true)
	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
		return
	}
	fmt.Println(string(dump))
	fmt.Fprintf(w, "<html><body>hello</body></html>\n")
}

func handlerChunkedResponseFlushSwitch(w http.ResponseWriter, r *http.Request, flush bool) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		panic("expected http.ResponseWriter to be an http.Flusher")
	}
	for i := 1; i <= 10; i++ {
		fmt.Fprintf(w, "Chunk #%d\n", i)
		if flush {
			flusher.Flush()
		}
		time.Sleep(100 * time.Millisecond)
	}
	flusher.Flush()
}

func handlerChunkedResponse(w http.ResponseWriter, r *http.Request) {
	handlerChunkedResponseFlushSwitch(w, r, true)
}

func handlerChunkedResponseNoFlush(w http.ResponseWriter, r *http.Request) {
	handlerChunkedResponseFlushSwitch(w, r, false)
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	var httpServer http.Server
	http.HandleFunc("/", handler)
	http.HandleFunc("/chunked", handlerChunkedResponse)
	http.HandleFunc("/chunked/noflush", handlerChunkedResponseNoFlush)
	log.Println("start http listening :" + port)
	httpServer.Addr = ":" + port
	log.Println(httpServer.ListenAndServe())
}
