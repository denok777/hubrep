package main

import (
	"flag"
	"fmt"
	"log"
)

var ReportsDir string

func init() {
	flag.StringVar(&ReportsDir, "p", "path", "reports directory path")
	flag.Parse()
}

func main() {
	files, err := FileList(ReportsDir)
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {
		fmt.Printf("%s -- %s\n", file.Name(), file.Time())
	}

	fmt.Println()
	fmt.Println("ok")
}
