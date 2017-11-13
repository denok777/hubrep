package main

import (
	"io/ioutil"
	"time"
)

type Report struct {
	name    string
	modtime time.Time
}

func (r *Report) Name() string {
	return r.name
}

func (r *Report) Time() string {
	return r.modtime.Format(time.RFC3339)
}

func FileList(dir string) ([]Report, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return []Report{}, err
	}

	list := []Report{}
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		list = append(list, Report{
			file.Name(),
			file.ModTime(),
		})
	}

	return list, nil
}
