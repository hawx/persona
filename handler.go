package persona

import (
	"log"
	"net/http"
)

func New(store Store, audience string, users []string) PersonaHandlers {
	return PersonaHandlers{
		SignIn:  signInHandler{store, audience},
		SignOut: signOutHandler{store},
		Protect: Protector(store, users),
		Switch:  Switcher(store, users),
	}
}

type Protect func(http.Handler) http.Handler
type Switch func(http.Handler, http.Handler) http.Handler

type PersonaHandlers struct {
	SignIn  http.Handler
	SignOut http.Handler
	Protect Protect
	Switch  Switch
}

func isSignedIn(toCheck string, users []string) bool {
	for _, user := range users {
		if user == toCheck {
			return true
		}
	}
	return false
}

func Switcher(store Store, users []string) Switch {
	return func(good, bad http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !isSignedIn(store.Get(r), users) {
				bad.ServeHTTP(w, r)
				return
			}

			good.ServeHTTP(w, r)
		})
	}
}

var forbidden = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "403 forbidden", http.StatusForbidden)
})

func Protector(store Store, users []string) Protect {
	return func(handler http.Handler) http.Handler {
		return Switcher(store, users)(handler, forbidden)
	}
}

type signInHandler struct {
	store    Store
	audience string
}

func (s signInHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	assertion := r.PostFormValue("assertion")
	email, err := assert(s.audience, assertion)

	if err != nil {
		log.Print("persona:", err)
		w.WriteHeader(403)
		return
	}

	s.store.Set(email, w, r)
	w.WriteHeader(200)
}

type signOutHandler struct {
	store Store
}

func (s signOutHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.store.Set("-", w, r)
	http.Redirect(w, r, "/", 307)
}
