package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"go-teach-me/database/fileIO"
	"go-teach-me/database/users"
	"go-teach-me/sessionStore"

	"github.com/gorilla/mux"
)

func ping(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "PONG")
}

func download(w http.ResponseWriter, r *http.Request) {
	fileID := mux.Vars(r)["file_id"]
	filename, data := fileIO.GetFile(fileID)
	if filename == "" {
		http.NotFound(w, r)
		return
	}
	http.ServeContent(w, r, filename, time.Now(), bytes.NewReader(data))
}

func login(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		http.Redirect(w, r, "/", 307)
		return
	}
	session := sessionStore.Get(r)
	r.ParseForm()
	user, err := users.GetUser(r.Form["email"][0], r.Form["password"][0])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	session.Values["authenticatedUser"] = user
	session.Save(r, w)
	http.Redirect(w, r, "/", 307)
}

func logout(w http.ResponseWriter, r *http.Request) {
	session := sessionStore.Get(r)
	session.Options.MaxAge = -1
	session.Save(r, w)
	http.Redirect(w, r, "/", 307)
}

func register(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	err := users.InsertUser(r.Form["first-name"][0], r.Form["last-name"][0], r.Form["email"][0], r.Form["password"][0], 2)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", 307)
}

func self(w http.ResponseWriter, r *http.Request) {
	user := sessionStore.GetSessionUser(r)
	if user == nil {
		http.NotFound(w, r)
		return
	}
	js, err := json.Marshal(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func upload(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(32 << 20)
	file, handler, err := r.FormFile("uploadfile")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()
	fileIO.InsertFile(file, handler)
	http.Redirect(w, r, "/retswerk/upload", 307)
}

func MountControllersRouter(s *mux.Router) {
	s.HandleFunc("/", ping)
	s.HandleFunc("/download/{file_id}/{filename}", download)
	s.HandleFunc("/login", login)
	s.HandleFunc("/logout", logout)
	s.HandleFunc("/upload", upload)
	s.HandleFunc("/register", register)
	s.HandleFunc("/self", self)
}
