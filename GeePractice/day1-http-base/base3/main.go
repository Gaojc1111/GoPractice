package main

import (
	"fmt"
	"gee"
	"log"
	"net/http"
)

func main() {
	r := gee.New()
	r.Get("/", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "URL.PATH = %q", req.URL.Path)
	})

	r.Get("/hello", func(w http.ResponseWriter, req *http.Request) {
		for k, v := range req.Header {
			fmt.Fprintf(w, "Header[%q] = %q\n", k, v)
		}
	})
	log.Fatal(r.Run(":9999"))
}
