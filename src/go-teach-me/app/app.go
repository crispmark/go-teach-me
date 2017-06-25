package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
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
	filename, data := fileIO.GetFile(2)
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

func main() {
	err := database.Initialize()
	if err != nil {
		panic(err)
	}

	http.Handle("/", http.FileServer(http.Dir("public")))
	http.HandleFunc("/ping", ping)

	http.HandleFunc("/upload", upload)
	http.HandleFunc("/download", download)
	http.HandleFunc("/register", createUserHandler)
	http.HandleFunc("/login", login)
	http.ListenAndServe(":8080", nil)
}
