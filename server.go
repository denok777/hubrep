package main

import (
	"fmt"
	"net/http"
)

func listHandler(w http.ResponseWriter, r *http.Request) {
	files, err := FileList(ReportsDir)
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}
	for _, file := range files {
		fmt.Fprintf(w, "%s -- %s\n", file.Name(), file.Time())
	}
}
