package main

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"sync"
)

type UserRole string

const (
	StudentUser UserRole = "student1"
	AdminUser   UserRole = "admin1"
)

var (
	sessionsMu sync.RWMutex
	sessions   = map[string]UserRole{}
)

func randomID() string {
	bytes := make([]byte, 32)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

func createSession(writer http.ResponseWriter, user UserRole) {
	sessionID := randomID()

	sessionsMu.Lock()
	sessions[sessionID] = user
	sessionsMu.Unlock()

	http.SetCookie(writer, &http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
}

func currentUser(request *http.Request) (UserRole, bool) {
	cookie, err := request.Cookie("session_id")
	if err != nil {
		return "", false
	}

	sessionsMu.RLock()
	user, ok := sessions[cookie.Value]
	sessionsMu.RUnlock()

	return user, ok
}

func clearSession(writer http.ResponseWriter, request *http.Request) {
	cookie, err := request.Cookie("session_id")
	if err == nil {
		sessionsMu.Lock()
		delete(sessions, cookie.Value)
		sessionsMu.Unlock()
	}

	http.SetCookie(writer, &http.Cookie{
		Name:   "session_id",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})
}
