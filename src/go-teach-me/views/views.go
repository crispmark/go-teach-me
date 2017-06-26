package views

import (
	"html/template"
	"net/http"

	"go-teach-me/models/files"
	"go-teach-me/models/users"
	"go-teach-me/sessionStore"

	"github.com/gorilla/mux"
)

func renderIndex(w http.ResponseWriter, user *users.User, fileArray []*files.File) {
	data := struct {
		User  *users.User
		Files []*files.File
	}{
		user,
		fileArray,
	}
	t, _ := template.ParseFiles("templates/index.html", "templates/nav.html")
	t.Execute(w, data)
}

func renderLogin(w http.ResponseWriter) {
	t, _ := template.ParseFiles("templates/login.html")
	t.Execute(w, nil)
}

func renderRegister(w http.ResponseWriter) {
	t, _ := template.ParseFiles("templates/register.html")
	t.Execute(w, nil)
}

func renderUpload(w http.ResponseWriter, user *users.User) {
	t, _ := template.ParseFiles("templates/upload.html", "templates/nav.html")
	t.Execute(w, nil)
}

func index(w http.ResponseWriter, r *http.Request) {
	user := sessionStore.GetSessionUser(r)
	if user == nil {
		renderLogin(w)
	} else {
		files, _ := files.GetAllFileInfo()
		renderIndex(w, user, files)
	}
}

func login(w http.ResponseWriter, r *http.Request) {
	user := sessionStore.GetSessionUser(r)
	if user != nil {
		http.Redirect(w, r, "/retswerk/", 307)
	} else {
		renderLogin(w)
	}
}

func register(w http.ResponseWriter, r *http.Request) {
	user := sessionStore.GetSessionUser(r)
	if user != nil {
		http.Redirect(w, r, "/retswerk/", 307)
	} else {
		renderRegister(w)
	}
}

func upload(w http.ResponseWriter, r *http.Request) {
	user := sessionStore.GetSessionUser(r)
	if user == nil {
		http.Redirect(w, r, "/retswerk/", 307)
	} else {
		renderUpload(w, user)
	}
}

func MountViewsRouter(s *mux.Router) {
	s.PathPrefix("/static").Handler(http.StripPrefix("/retswerk/static/", http.FileServer(http.Dir("public"))))
	s.HandleFunc("/", index)
	s.HandleFunc("/login", login)
	s.HandleFunc("/register", register)
	s.HandleFunc("/upload", upload)
}
