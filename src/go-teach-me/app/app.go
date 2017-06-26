package main

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
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
	session, _ := store.Get(r, "session")
	session.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7,
		Secure:   false,
		HttpOnly: true,
	}
	r.ParseForm()
	user, err := users.GetUser(r.Form["email"][0], r.Form["password"][0])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	session.Values["authenticatedUser"] = user
	session.Values["userID"] = user.UserID
	session.Save(r, w)

	fmt.Fprintf(w, "Signed in!")
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
	gob.Register(users.User{})
	err := database.Initialize()
	if err != nil {
		panic(err)
	}

	r := mux.NewRouter()
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/static/", 301)
	})
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("public"))))
	r.HandleFunc("/ping", ping)
	r.HandleFunc("/upload", upload)
	r.HandleFunc("/download/{file_id}/{filename}", download)
	r.HandleFunc("/register", createUserHandler)
	r.HandleFunc("/login", login)
	r.HandleFunc("/self", self)

	srv := &http.Server{
		Handler:      r,
		Addr:         "127.0.0.1:8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
