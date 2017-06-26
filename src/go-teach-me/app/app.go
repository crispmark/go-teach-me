package main

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"go-teach-me/database"
	"go-teach-me/database/fileIO"
	"go-teach-me/database/users"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

var (
	secret = "2qbyjHaYmdQgQNjJ"
	store  = sessions.NewCookieStore([]byte("something-very-secret"))
)

func getSessionUser(r *http.Request) *users.User {
	session, _ := store.Get(r, "session")
	val := session.Values["authenticatedUser"]
	if user, ok := val.(users.User); !ok {
		return nil
	} else {
		return &user
	}
}

func ping(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "PONG")
}

func register(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		t, _ := template.ParseFiles("templates/register.html")
		t.Execute(w, nil)
		return
	}
	r.ParseForm()
	err := users.InsertUser(r.Form["first-name"][0], r.Form["last-name"][0], r.Form["email"][0], r.Form["password"][0], 2)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", 307)
}

func login(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		http.Redirect(w, r, "/", 307)
		return
	}
	session, _ := store.Get(r, "session")
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
	session, _ := store.Get(r, "session")
	session.Options.MaxAge = -1
	session.Save(r, w)
	t, _ := template.ParseFiles("templates/login.html")
	t.Execute(w, nil)
}

func self(w http.ResponseWriter, r *http.Request) {
	user := getSessionUser(r)
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

func download(w http.ResponseWriter, r *http.Request) {
	fileID := mux.Vars(r)["file_id"]
	filename, data := fileIO.GetFile(fileID)
	if filename == "" {
		http.NotFound(w, r)
		return
	}
	http.ServeContent(w, r, filename, time.Now(), bytes.NewReader(data))
}

func upload(w http.ResponseWriter, r *http.Request, user *users.User) {
	if r.Method == "GET" {
		t, _ := template.ParseFiles("templates/upload.html", "templates/nav.html")
		t.Execute(w, user)
	} else {
		r.ParseMultipartForm(32 << 20)
		file, handler, err := r.FormFile("uploadfile")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer file.Close()
		t, _ := template.ParseFiles("templates/upload.html", "templates/nav.html")
		t.Execute(w, user)
		fileIO.InsertFile(file, handler)
	}
}

func authHandleFunc(r *mux.Router, pattern string, handler func(http.ResponseWriter, *http.Request, *users.User)) {
	r.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		user := getSessionUser(r)
		if user == nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		handler(w, r, user)
	})
}

func index(w http.ResponseWriter, r *http.Request) {
	user := getSessionUser(r)
	if user == nil {
		t, _ := template.ParseFiles("templates/login.html")
		t.Execute(w, "")
		return
	}
	files, _ := fileIO.GetAllFileInfo()
	t, _ := template.ParseFiles("templates/index.html", "templates/nav.html")
	t.Execute(w, userFiles{User: user, Files: files})
}

type userFiles struct {
	User  *users.User
	Files *[]fileIO.File
}

func main() {
	gob.Register(users.User{})
	err := database.Initialize()
	if err != nil {
		panic(err)
	}

	r := mux.NewRouter()
	r.HandleFunc("/", index)
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("public"))))
	r.HandleFunc("/ping", ping)
	authHandleFunc(r, "/upload", upload)
	r.HandleFunc("/download/{file_id}/{filename}", download)
	r.HandleFunc("/register", register)
	r.HandleFunc("/login", login)
	r.HandleFunc("/logout", logout)
	r.HandleFunc("/self", self)

	srv := &http.Server{
		Handler:      r,
		Addr:         "127.0.0.1:8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
