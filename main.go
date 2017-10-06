package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

type helloHandler struct {}

func (helloHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello, you've hit %s\n", r.URL.Path)
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	err := http.ListenAndServe(":" + port, helloHandler{})
	log.Fatal(err)
}
