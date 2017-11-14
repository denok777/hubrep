package main

import (
	"fmt"
	"github.com/gorilla/sessions"
	"html/template"
	"net/http"
)

var (
	key   = []byte("todo-any-secret-key")
	store = sessions.NewCookieStore(key)
)

const (
	SessionName = "reports-app"

	UrlSignIn  = "/signin"
	UrlReports = "/reports"
	UrlAuth    = "/auth"
)

func auth(w http.ResponseWriter, r *http.Request) {
	s, _ := store.Get(r, SessionName)

	r.ParseForm()
	publisher := r.Form.Get("publisher")
	if len(publisher) == 0 {
		http.Error(w, "Username is required", http.StatusBadRequest)
		return
	}

	// TODO: list of available publishers
	publishers := []string{
		"fry",
		"lila",
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
