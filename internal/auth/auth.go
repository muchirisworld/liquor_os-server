package auth

import "net/http"

type AuthProvider interface {
	RequireAuth() func(http.Handler) http.Handler
}