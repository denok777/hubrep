package main

import (
	"flag"
	"fmt"
	"github.com/gorilla/context"
	"log"
	"net/http"
	"os"
)

var ReportsDir string
var PublishersList string
var Port int

func init() {
	flag.StringVar(&ReportsDir, "d", "", "reports directory")
	flag.StringVar(&PublishersList, "u", "", "publishers list file")
	flag.IntVar(&Port, "p", 4242, "application port")
	flag.Parse()

	if len(PublishersList) == 0 {
		log.Fatal("publishers list file is required")
	}

	file, err := os.Open(PublishersList)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}
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
