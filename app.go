package main

import (
	"encoding/gob"
	"net/http"

	"go-teach-me/actions"
	"go-teach-me/controllers"
	"go-teach-me/database"
	"go-teach-me/models/users"
	"go-teach-me/views"

	"github.com/gorilla/mux"
	"google.golang.org/appengine"
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
	actions.MountActionsRouter(r.PathPrefix("/actions").Subrouter())

	http.Handle("/", r)
	appengine.Main()
}
