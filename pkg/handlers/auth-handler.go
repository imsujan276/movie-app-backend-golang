package handlers

import (
	"backend/models"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/pascaldekloe/jwt"
	"golang.org/x/crypto/bcrypt"
)

var validUser = models.User{
	ID:       1,
	Email:    "test@gmail.com",
	Password: func() string { return getHasedPassword("password123") }(),
}

func getHasedPassword(password string) string {
	hashedPasswoed, _ := bcrypt.GenerateFromPassword([]byte(password), 12)
	return string(hashedPasswoed)
}

func comparePasswords(p1, p2 string) error {
	return bcrypt.CompareHashAndPassword([]byte(p1), []byte(p2))
}

type Credentials struct {
	UserName string `json:"email"`
	Password string `json:"password"`
}

func (m *Repository) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var creds Credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		m.errorJSON(w, errors.New("Unauthorized"))
		return
	}

	if strings.ToLower(creds.UserName) != strings.ToLower(validUser.Email) {
		m.errorJSON(w, errors.New("Unauthorized User"))
		return
	}

	hasedPassword := validUser.Password
	err = comparePasswords(hasedPassword, creds.Password)
	if err != nil {
		m.errorJSON(w, errors.New("Unauthorized"))
		return
	}

	var claims jwt.Claims
	claims.Subject = fmt.Sprint(validUser.ID)
	claims.Issued = jwt.NewNumericTime(time.Now())
	claims.NotBefore = jwt.NewNumericTime(time.Now())
	claims.Expires = jwt.NewNumericTime(time.Now().Add(24 * time.Hour))
	claims.Issuer = "mydomain.com"
	claims.Audiences = []string{"mydomain.com"}

	jwtBytes, err := claims.HMACSign(jwt.HS256, []byte(m.App.Config.Jwt.Secret))
	if err != nil {
		m.errorJSON(w, errors.New("Error Signing in"))
		return
	}

	m.writeJSON(w, http.StatusOK, string(jwtBytes), "response")
}
