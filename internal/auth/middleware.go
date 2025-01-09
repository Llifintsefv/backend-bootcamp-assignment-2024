package auth

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type JWTClaims struct {
	exp  int64  `json:"exp"`
	Type string `json:"user_type"`
	jwt.RegisteredClaims
}

func AuthMiddleware(secretKey []byte) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := r.Header.Get("Authorization")
			if token == "" {
				http.Error(w, "Unauthorized dsd", http.StatusUnauthorized)
				return
			}
			headerPart := strings.Split(token, " ")
			if len(headerPart) != 2 || headerPart[0] != "Bearer" {
				http.Error(w, "Unauthorized vav", http.StatusUnauthorized)
				return
			}
			tokenString := headerPart[1]

			claims, err := validateToken(tokenString, secretKey)
			if err != nil {
				http.Error(w, "Unauthorized huz", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), "user_id", claims.exp)
			ctx = context.WithValue(ctx, "user_type", claims.Type)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func validateToken(tokenString string, secretKey []byte) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("uxpecred signing method: %v", t.Header["alg"])
		}
		return secretKey, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token claims")
	}

	return claims, nil
}

func ModeratorMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userType, ok := r.Context().Value("user_type").(string)
		if !ok {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		if userType != "moderator" {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}
