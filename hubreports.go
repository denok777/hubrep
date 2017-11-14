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
	http.Handle(UrlSignIn, http.HandlerFunc(signin))
	http.Handle(UrlAuth, http.HandlerFunc(auth))
	http.Handle(UrlReports, verify(reports))
	http.Handle(UrlDownload, verify(download))

	err := http.ListenAndServe(
		fmt.Sprintf(":%d", Port),
		context.ClearHandler(http.DefaultServeMux),
	)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
