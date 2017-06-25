package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"go-teach-me/database"
	"go-teach-me/database/fileIO"
	"go-teach-me/database/users"
)

func ping(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "PONG")
}

func createUserHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	err := users.InsertUser(r.Form["first-name"][0], r.Form["last-name"][0], r.Form["email"][0], r.Form["password"][0], 2)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "PONG")
}

func login(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	user, err := users.GetUser(r.Form["email"][0], r.Form["password"][0])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
	fileID := strings.Split(r.URL.Path, "/")[2]
	if fileID == "" {
		http.Error(w, "No file ID provided", http.StatusInternalServerError)
		return
	}
	filename, data := fileIO.GetFile(fileID)
	http.ServeContent(w, r, filename, time.Now(), bytes.NewReader(data))
}

func upload(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		http.Redirect(w, r, "/upload.html", 301)
	} else {
		r.ParseMultipartForm(32 << 20)
		file, handler, err := r.FormFile("uploadfile")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer file.Close()
		http.Redirect(w, r, "/upload.html", 301)
		fileIO.InsertFile(file, handler)
	}
}

type customHandler struct {
	defaultHandler http.Handler
}

func (ch customHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if strings.HasPrefix(r.URL.Path, "/download/") {
		download(w, r)
		return
	}

	ch.defaultHandler.ServeHTTP(w, r)
}

func main() {
	err := database.Initialize()
	if err != nil {
		panic(err)
	}

	http.Handle("/", customHandler{http.FileServer(http.Dir("public"))})
	http.HandleFunc("/ping", ping)

	http.HandleFunc("/upload", upload)
	http.HandleFunc("/register", createUserHandler)
	http.HandleFunc("/login", login)

	http.ListenAndServe(":8080", nil)
}
