package middleware

import (
	"net/http"
	"strings"

	"github.com/drofloh/lenslocked.com/context"
	"github.com/drofloh/lenslocked.com/models"
)

// User ...
type User struct {
	models.UserService
}

// Apply ...
func (mw *User) Apply(next http.Handler) http.HandlerFunc {
	return mw.ApplyFn(next.ServeHTTP)
}

// ApplyFn ...
func (mw *User) ApplyFn(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		// if the user is requesting a static asset or image we will not need to
		// look up the currect user.
		if strings.HasPrefix(path, "/assets/") || strings.HasPrefix(path, "/images/") {
			next(w, r)
			return
		}

		cookie, err := r.Cookie("remember_token")
		if err != nil {
			next(w, r)
			return
		}
		user, err := mw.UserService.ByRemember(cookie.Value)
		if err != nil {
			next(w, r)
			return
		}
		ctx := r.Context()
		ctx = context.WithUser(ctx, user)
		r = r.WithContext(ctx)
		next(w, r)
	})
}

// RequireUser assumes that User has already been run otherwise it wil not work
// correctly.
type RequireUser struct {
	User
}

// Apply ...
func (mw *RequireUser) Apply(next http.Handler) http.HandlerFunc {
	return mw.ApplyFn(next.ServeHTTP)
}

// ApplyFn ...
func (mw *RequireUser) ApplyFn(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := context.User(r.Context())
		if user == nil {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		next(w, r)
	})
}
