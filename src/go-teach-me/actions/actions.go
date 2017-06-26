package actions

import (
	"net/http"

	"go-teach-me/models/files"
	"go-teach-me/models/users"
	"go-teach-me/sessionStore"

	"github.com/gorilla/mux"
)

func login(w http.ResponseWriter, r *http.Request) {
	session := sessionStore.Get(r)
	r.ParseForm()
	user, err := users.GetUser(r.Form["email"][0], r.Form["password"][0])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	session.Values["authenticatedUser"] = user
	session.Save(r, w)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func logout(w http.ResponseWriter, r *http.Request) {
	session := sessionStore.Get(r)
	session.Options.MaxAge = -1
	session.Save(r, w)
	http.Redirect(w, r, "/", 303)
}

func register(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	err := users.InsertUser(r.Form["first-name"][0], r.Form["last-name"][0], r.Form["email"][0], r.Form["password"][0], 2)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", 303)
}

func upload(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(32 << 20)
	file, handler, err := r.FormFile("uploadfile")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()
	files.InsertFile(file, handler)
	http.Redirect(w, r, "/retswerk/upload", 303)
}

func MountActionsRouter(s *mux.Router) {
	s.HandleFunc("/login", login)
	s.HandleFunc("/logout", logout)
	s.HandleFunc("/upload", upload)
	s.HandleFunc("/register", register)
}
