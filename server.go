package main

import (
	"fmt"
	"html/template"
	"net/http"
)

func listHandler(w http.ResponseWriter, r *http.Request) {
	files, err := FileList(ReportsDir)
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}

	t, _ := template.ParseFiles("list.html")
	t.Execute(w, files)
}
