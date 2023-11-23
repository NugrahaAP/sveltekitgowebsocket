package misc

import (
	"net/http"
	"time"

	"github.com/google/uuid"
)

func ValidateUUID(u string) bool {
	_, err := uuid.Parse(u)
	return err == nil
}

func CreateCookie(name string, value string) *http.Cookie {
	return &http.Cookie{
		Name:     name,
		Value:    value,
		Expires:  time.Now().Add(time.Minute * 15),
		HttpOnly: true,
		Path:     "/",
		Secure:   false,
		SameSite: http.SameSiteStrictMode,
	}
}
