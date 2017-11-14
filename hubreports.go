package main

import (
	"flag"
	"fmt"
	"github.com/gorilla/context"
	"log"
	"net/http"
)

var ReportsDir string
var Port int

func init() {
	flag.StringVar(&ReportsDir, "d", "", "reports directory")
	flag.IntVar(&Port, "p", 4242, "application port")
	flag.Parse()
}

func main() {
	http.Handle("/login", http.HandlerFunc(login))
	http.Handle("/list", verify(listHandler))

	err := http.ListenAndServe(
		fmt.Sprintf(":%d", Port),
		context.ClearHandler(http.DefaultServeMux),
	)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
