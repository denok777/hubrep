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

const SessionName = "reports-app"

func login(w http.ResponseWriter, r *http.Request) {
	s, _ := store.Get(r, SessionName)

	r.ParseForm()
	user := r.Form.Get("user")
	if len(user) == 0 {
		http.Error(w, "Username is required", http.StatusBadRequest)
		return
	}

	// TODO: list of available users
	users := []string{
		"fry",
		"lila",
	}

	var found bool
	for _, name := range users {
		if name == user {
			found = true
			break
		}
	}

	if !found {
		http.Error(w, "Unknown user", http.StatusForbidden)
		return
	}

	s.Values["auth"] = true
	s.Values["publisher"] = user
	s.Save(r, w)
	http.Redirect(w, r, "/list", http.StatusFound)
}

func signin(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("login-form.html")
	t.Execute(w, "")
}

func listHandler(w http.ResponseWriter, r *http.Request) {
	files, err := FileList(ReportsDir)
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
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		h.ServeHTTP(w, r)
	})
}

func verify(fn http.HandlerFunc) http.Handler {
	return verifyUser(http.HandlerFunc(fn))
}
