package main

import (
	"bufio"
	"fmt"
	"github.com/gorilla/sessions"
	"html/template"
	"io"
	"net/http"
	"os"
	"path"
	"strings"
)

var (
	key   = []byte("todo-any-secret-key")
	store = sessions.NewCookieStore(key)
)

const (
	SessionName = "reports-app"

	UrlSignIn   = "/signin"
	UrlAuth     = "/auth"
	UrlReports  = "/reports"
	UrlDownload = "/download"
)

func auth(w http.ResponseWriter, r *http.Request) {
	s, _ := store.Get(r, SessionName)

	r.ParseForm()
	publisher := r.Form.Get("publisher")
	if len(publisher) == 0 {
		http.Error(w, "Username is required", http.StatusBadRequest)
		return
	}

	publishers := []string{}
	file, err := os.Open(PublishersList)
	if err != nil {
		http.Error(w, "Critical error", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		publishers = append(publishers, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		http.Error(w, "Critical error", http.StatusInternalServerError)
		return
	}

	var found bool
	for _, name := range publishers {
		if found = (name == publisher); found {
			break
		}
	}

	if !found {
		http.Error(w, "Unknown publisher", http.StatusForbidden)
		return
	}

	s.Values["auth"] = true
	s.Values["publisher"] = publisher
	s.Save(r, w)
	http.Redirect(w, r, UrlReports, http.StatusFound)
}

func signin(w http.ResponseWriter, r *http.Request) {
	s, _ := store.Get(r, SessionName)
	if auth, ok := s.Values["auth"].(bool); ok && auth {
		http.Redirect(w, r, UrlReports, http.StatusSeeOther)
		return
	}

	t, _ := template.ParseFiles("login-form.html")
	t.Execute(w, "")
}

func reports(w http.ResponseWriter, r *http.Request) {
	s, _ := store.Get(r, SessionName)
	var publisher string
	var ok bool
	if publisher, ok = s.Values["publisher"].(string); !ok {
		// TODO expressive error message
		http.Error(w, "Critical error", http.StatusInternalServerError)
		return
	}

	files, err := FileList(ReportsDir, publisher)
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}

	t, _ := template.ParseFiles("list.html")
	t.Execute(w, files)
}

func download(w http.ResponseWriter, r *http.Request) {
	fname := strings.TrimSpace(r.URL.Query().Get("report"))
	if len(fname) == 0 {
		http.Error(w, "Report filename is required.", http.StatusBadRequest)
		return
	}

	fname = path.Clean(fname)
	fname = strings.Replace(fname, "/", "", -1)

	s, _ := store.Get(r, SessionName)
	var publisher string
	var ok bool
	if publisher, ok = s.Values["publisher"].(string); !ok {
		// TODO expressive error message
		http.Error(w, "Critical error", http.StatusInternalServerError)
		return
	}

	idx := len(publisher)
	if len(fname) <= idx || fname[:idx] != publisher {
		http.Error(w, "File not found.", http.StatusNotFound)
		return
	}

	fpath := ReportsDir + "/" + fname

	file, err := os.Open(fpath)
	defer file.Close()
	if err != nil {
		http.Error(w, "File not found.", http.StatusNotFound)
		return
	}

	var fi os.FileInfo
	if fi, err = file.Stat(); err != nil {
		http.Error(w, "Critical error", http.StatusInternalServerError)
		return
	}

	w.Header().Set(
		"Content-Disposition",
		fmt.Sprintf("attachment; filename=%s", fname),
	)
	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", fi.Size()))

	io.Copy(w, file)
}

func verifyUser(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s, _ := store.Get(r, SessionName)
		if auth, ok := s.Values["auth"].(bool); !ok || !auth {
			http.Redirect(w, r, UrlSignIn, http.StatusSeeOther)
			return
		}

		h.ServeHTTP(w, r)
	})
}

func verify(fn http.HandlerFunc) http.Handler {
	return verifyUser(http.HandlerFunc(fn))
}
