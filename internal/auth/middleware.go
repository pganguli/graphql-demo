package auth

import (
	"context"
	"net/http"
    "log"
	"strings"

	"github.com/pganguli/hnews/internal/users"
	"github.com/pganguli/hnews/pkg/jwt"
)

var userCtxKey = &contextKey{"user"}

type contextKey struct {
	name string
}

func BearerTokenExtract(authHeader string) string {
	parts := strings.Split(authHeader, "Bearer")
	if len(parts) != 2 {
		return ""
	}

	token := strings.TrimSpace(parts[1])
	if len(token) < 1 {
		return ""
	}

	return token
}

func Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		token := BearerTokenExtract(authHeader)

		// Allow unauthenticated users in
		if token == "" {
			next.ServeHTTP(w, r)
			return
		}

		// Validate jwt token
		username, err := jwt.ParseToken(token)
		if err != nil {
            log.Print(err)
			http.Error(w, "Invalid token", http.StatusForbidden)
			return
		}

		// Create user and check if user exists in db
		id, err := users.GetUserIdByUsername(username)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}
		user := users.User{ID: id, Username: username}

		// Put it in context
		ctx := context.WithValue(r.Context(), userCtxKey, &user)

		// and call the next with our new context
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

// ForContext finds the user from the context. REQUIRES Middleware to have run.
func ForContext(ctx context.Context) *users.User {
	raw, _ := ctx.Value(userCtxKey).(*users.User)
	return raw
}
