package main

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

const sessionTimeoutMinutes = 30

type SessionState struct {
	lock     sync.Mutex
	sessions map[string]time.Time
}

func InitSessionState() (s SessionState) {
	s.sessions = make(map[string]time.Time)
	go s.gc()
	return
}

func (s SessionState) Make(w http.ResponseWriter) error {
	id := make([]byte, 32, 32)
	_, err := rand.Read(id)
	if err != nil {
		return errors.New(fmt.Sprintf("Cannot generate random numbers for session id: %s", err))
	}
	idString := hex.EncodeToString(id)
	http.SetCookie(w, &http.Cookie{Name: "session", Value: idString, MaxAge: sessionTimeoutMinutes * 60, HttpOnly: true, Path: "/"})
	s.lock.Lock()
	s.sessions[idString] = time.Now().Add(time.Duration(sessionTimeoutMinutes) * time.Minute)
	s.lock.Unlock()
	log.Printf("Created session, id: %s\n", idString)
	return nil
}

func (s SessionState) Delete(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session")
	if err != nil {
		return
	}
	http.SetCookie(w, &http.Cookie{Name: "session", Value: "", MaxAge: 0, HttpOnly: true, Path: "/"})
	id := cookie.Value
	s.lock.Lock()
	delete(s.sessions, id)
	s.lock.Unlock()
}

func (s SessionState) IsLoggedIn(r *http.Request) bool {
	cookie, err := r.Cookie("session")
	if err != nil {
		log.Println("Missing session cookie")
		return false
	}
	id := cookie.Value
	s.lock.Lock()
	defer s.lock.Unlock()
	expirationTime, exists := s.sessions[id]
	if !exists {
		log.Println("Given session does not exist")
		return false
	}
	stillValid := time.Now().Before(expirationTime)
	if !stillValid {
		log.Printf("Session expired: %s\n", id)
		delete(s.sessions, id)
		return false
	}
	log.Printf("Valid session: %s\n", id)
	return true
}

func (s SessionState) gc() {
	for {
		time.Sleep(time.Duration(sessionTimeoutMinutes) * time.Minute)
		now := time.Now()
		s.lock.Lock()
		for session, expirationTime := range s.sessions {
			stillValid := now.Before(expirationTime)
			if !stillValid {
				delete(s.sessions, session)
				log.Printf("Deleted expired session: %s\n", session)
			}
		}
		s.lock.Unlock()
	}
}
