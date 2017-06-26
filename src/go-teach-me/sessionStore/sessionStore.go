package sessionStore

import (
	"net/http"

	"go-teach-me/database/users"

	"github.com/gorilla/sessions"
)

var (
	secret = "2qbyjHaYmdQgQNjJ"
	store  = sessions.NewCookieStore([]byte(secret))
)

func GetSessionUser(r *http.Request) *users.User {
	session, _ := store.Get(r, "session")
	val := session.Values["authenticatedUser"]
	if user, ok := val.(users.User); !ok {
		return nil
	} else {
		return &user
	}
}

func Get(r *http.Request) *sessions.Session {
	session, _ := store.Get(r, "session")
	return session
}
