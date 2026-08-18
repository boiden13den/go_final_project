package api

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("todo-secret-key")

func passHash(pass string) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(pass)))
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJson(w, http.StatusBadRequest, map[string]any{"error": err.Error()})
		return
	}
	if req.Password != appPassword {
		writeJson(w, http.StatusUnauthorized, map[string]any{"error": "Неверный пароль"})
		return
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"hash": passHash(req.Password),
	})
	signed, err := token.SignedString(jwtSecret)
	if err != nil {
		writeJson(w, http.StatusInternalServerError, map[string]any{"error": err.Error()})
		return
	}
	writeJson(w, http.StatusOK, map[string]any{"token": signed})
}

func auth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if appPassword == "" {
			next(w, r)
			return
		}
		cookie, err := r.Cookie("token")
		if err != nil {
			http.Error(w, "Authentification required", http.StatusUnauthorized)
			return
		}
		token, err := jwt.Parse(cookie.Value, func(t *jwt.Token) (any, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method")
			}
			return jwtSecret, nil
		})
		if err != nil || !token.Valid {
			http.Error(w, "Authentification required", http.StatusUnauthorized)
			return
		}
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || claims["hash"] != passHash(appPassword) {
			http.Error(w, "Authentification required", http.StatusUnauthorized)
			return
		}
		next(w, r)
	}
}
