package main

import (
	"flag"
	"log"
	"net/http"
)

var ReportsDir string

func init() {
	flag.StringVar(&ReportsDir, "p", "path", "reports directory path")
	flag.Parse()
}

func main() {
	http.HandleFunc("/list", listHandler)

	err := http.ListenAndServe(":4242", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
