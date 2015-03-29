package persona

import (
	"net/http"

	"github.com/gorilla/sessions"
)

type Store interface {
	Set(email string, w http.ResponseWriter, r *http.Request)
	Get(r *http.Request) string
}

type emailStore struct {
	store sessions.Store
}

func NewStore(secret string) Store {
	return &emailStore{sessions.NewCookieStore([]byte(secret))}
}

func (s emailStore) Get(r *http.Request) string {
	session, _ := s.store.Get(r, "session-name")

	if v, ok := session.Values["email"].(string); ok {
		return v
	}

	return ""
}

func (s emailStore) Set(email string, w http.ResponseWriter, r *http.Request) {
	session, _ := s.store.Get(r, "session-name")
	session.Values["email"] = email
	session.Save(r, w)
}
