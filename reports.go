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

func FileList(dir, publisher string) ([]Report, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return []Report{}, err
	}

	idx := len(publisher)
	list := []Report{}
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		fname := file.Name()
		if len(fname) <= idx {
			continue
		}
		if fname[:idx] != publisher {
			continue
		}
		list = append(list, Report{
			fname,
			file.ModTime(),
		})
	}

	return list, nil
}
