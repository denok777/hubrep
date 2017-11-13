package main

import (
	"fmt"
	"io/ioutil"
	"log"
)

func main() {
	dir := "/tmp/test/"

	files, err := ioutil.ReadDir(dir)
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
