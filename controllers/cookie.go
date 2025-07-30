package controllers

import (
	"fmt"
	"net/http"
	"time"
)

const (
	CookieSession = "session"
)

func setCookie(w http.ResponseWriter, name, value string, expire ...bool) {
	cookie := http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		HttpOnly: true,
	}

	//invalidate the cookie
	if len(expire) > 0 && expire[0] {
		cookie.MaxAge = -1
		cookie.Expires = time.Unix(0, 0)
		cookie.Value = ""
	}
	http.SetCookie(w, &cookie)
}

func readCookie(r *http.Request, name string) (string, error) {
	cookie, err := r.Cookie(name)
	if err != nil {
		return "", fmt.Errorf("read cookie: %w", err)
	}
	return cookie.Value, nil
}

func deleteCookie(w http.ResponseWriter, name string) {
	setCookie(w, name, "", true)
}
