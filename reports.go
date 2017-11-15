package main

import (
	"io/ioutil"
	"sort"
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

type Reports []Report

func (r Reports) Len() int {
	return len(r)
}

func (r Reports) Less(i, j int) bool {
	return r[i].modtime.After(r[j].modtime)
}

func (r Reports) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}

func FileList(dir, publisher string) (Reports, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return []Report{}, err
	}

	idx := len(publisher)
	list := make(Reports, 0)
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

	sort.Sort(list)

	return list, nil
}
