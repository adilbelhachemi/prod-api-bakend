package server

import (
	"context"
	"net/http"
	"pratbacknd/internal/types"
	"strings"
)

const (
	username = "adil"
	password = "password"
)

func (s *Server) enableCORS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", s.allowedOrigins)
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == "OPTIONS" {
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type")
			return
		} else {
			h.ServeHTTP(w, r)
		}
	})
}

func (s *Server) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		un, pass, ok := r.BasicAuth()
		if !ok || un != username || pass != password {
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), "user", types.User{ID: un})
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (s *Server) AuthenticateV2(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		splits := strings.Split(authHeader, " ")
		if len(splits) != 2 {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		token, err := s.firebaseAuthClient.VerifyIDToken(r.Context(), splits[1])
		if err != nil {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		ctx := context.WithValue(r.Context(), "user", types.User{ID: token.UID})
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
