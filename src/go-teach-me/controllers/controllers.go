package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"go-teach-me/models/files"
	"go-teach-me/sessionStore"

	"github.com/gorilla/mux"
)

func ping(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "PONG")
}

func download(w http.ResponseWriter, r *http.Request) {
	fileID := mux.Vars(r)["file_id"]
	filename, data := files.GetFile(fileID)
	if filename == "" {
		http.NotFound(w, r)
		return
	}
	http.ServeContent(w, r, filename, time.Now(), bytes.NewReader(data))
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

func MountControllersRouter(s *mux.Router) {
	s.HandleFunc("/", ping)
	s.HandleFunc("/download/{file_id}/{filename}", download)
	s.HandleFunc("/self", self)
}
