package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)


func extraHeaderHandleFunc(w http.ResponseWriter, r *http.Request) {
	hj, ok := w.(http.Hijacker)
	if !ok {
		log.Fatal("Could not obtain Hijacker")
	}
	conn, bufrw, err := hj.Hijack()
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	message := "Hello World.\n"
	bufrw.WriteString("HTTP/1.1 200 OK\r\n")
	bufrw.WriteString("Content-Type: text/plain\r\n")
	fmt.Fprintf(bufrw, "Content-Length: %d\r\n", len([]byte(message)))
	bufrw.WriteString("\r\n")
	bufrw.WriteString(message)
	bufrw.Flush()
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
	h.HandleFunc("/mix", extraHeaderHandleFunc)

	log.Println("Listening at port " + port)
	err := http.ListenAndServe(":"+port, h)
	log.Fatal(err)
}
