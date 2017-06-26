package main

import (
	"encoding/gob"
	"log"
	"net/http"
	"time"

	"go-teach-me/controllers"
	"go-teach-me/database"
	"go-teach-me/database/users"
	"go-teach-me/views"

	"github.com/gorilla/mux"
)

func index(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/retswerk/", 307)
}

func main() {
	gob.Register(users.User{})
	err := database.Initialize()
	if err != nil {
		panic(err)
	}

	r := mux.NewRouter()
	r.HandleFunc("/", index)
	views.MountViewsRouter(r.PathPrefix("/retswerk").Subrouter())
	controllers.MountControllersRouter(r.PathPrefix("/api").Subrouter())

	srv := &http.Server{
		Handler:      r,
		Addr:         "127.0.0.1:8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
