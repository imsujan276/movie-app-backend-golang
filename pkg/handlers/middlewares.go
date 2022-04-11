package handlers

import (
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/pascaldekloe/jwt"
)

// CheckToken is a middleware to check authorized access
func (m *Repository) CheckToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Add("Vary", "Authorization")
		authheader := r.Header.Get("Authorization")
		if authheader == "" {
			// could set an anonymous user
		}

		headerParts := strings.Split(authheader, " ")
		if len(headerParts) != 2 {
			m.errorJSON(w, errors.New("Invalid auth header"))
			return
		}
		if headerParts[0] != "Bearer" {
			m.errorJSON(w, errors.New("Unauthorized - no bearer"))
			return
		}

		token := headerParts[1]

		claims, err := jwt.HMACCheck([]byte(token), []byte(m.App.Config.Jwt.Secret))
		if err != nil {
			m.errorJSON(w, errors.New("Unauthorized - failed hmac check"), http.StatusForbidden)
			return
		}
		if !claims.Valid(time.Now()) {
			m.errorJSON(w, errors.New("Unauthorized - token expired"), http.StatusForbidden)
			return
		}
		if !claims.AcceptAudience("mydomain.com") {
			m.errorJSON(w, errors.New("Unauthorized - Invalid Audience"), http.StatusForbidden)
			return
		}
		if claims.Issuer != "mydomain.com" {
			m.errorJSON(w, errors.New("Unauthorized - Invalid Issuer"), http.StatusForbidden)
			return
		}

		userId, err := strconv.ParseInt(claims.Subject, 10, 64)
		if err != nil {
			m.errorJSON(w, errors.New("Unauthorized"), http.StatusForbidden)
			return
		}

		log.Println("Valid user:", userId)

		next.ServeHTTP(w, r)
	})
}
