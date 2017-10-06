package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

type numberDumper int

func (n numberDumper) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "It is %d\n", n)
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	h := http.NewServeMux()
	h.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hello.\n")
	})
	h.HandleFunc("/404", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
		fmt.Fprint(w, "Nothing is here.\n")
	})
	h.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "You're at the root path.\n")
	})
	h.Handle("/one", numberDumper(1))
	h.Handle("/two", numberDumper(2))

	err := http.ListenAndServe(":" + port, h)
	log.Fatal(err)
}
