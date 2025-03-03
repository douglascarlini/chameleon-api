package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type Claims struct {
	UserID string `json:"user_id"`
	jwt.MapClaims
}

func generateToken(userID string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	basic := jwt.MapClaims{"exp": expirationTime}
	claims := &Claims{UserID: userID, MapClaims: basic}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(JWT_CODE))
}

func jwtAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing token", http.StatusUnauthorized)
			return
		}

		claims := &Claims{}
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
			return []byte(JWT_CODE), nil
		})
		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			fmt.Printf("ERROR: %s\n", err.Error())
			return
		}

		w.Header().Set("User-ID", claims.UserID)

		next.ServeHTTP(w, r)
	})
}
