package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
)

var ReportsDir string

func init() {
	flag.StringVar(&ReportsDir, "p", "path", "reports directory path")
	flag.Parse()
}

func main() {
	files, err := ioutil.ReadDir(ReportsDir)
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		fmt.Println(file.Name())
	}

	fmt.Println()
	fmt.Println("ok")
}
