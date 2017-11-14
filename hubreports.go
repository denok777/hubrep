package main

import (
	"flag"
	"fmt"
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
	http.Handle("/list", verify(listHandler))

	err := http.ListenAndServe(
		fmt.Sprintf(":%d", Port),
		nil,
	)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
