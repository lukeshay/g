package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/lukeshay/g/auth"
	adaptors "github.com/lukeshay/g/auth/adapters"
	"github.com/lukeshay/g/auth/encrypters"
	"github.com/lukeshay/g/auth/generators"
)

func main() {
	encrypter, err := encrypters.NewAesEncrypter("thisissupersecret")
	if err != nil {
		panic(err)
	}

	authManager := auth.New(auth.NewOptions{
		Adapter:   adaptors.NewInMemoryAdapter(),
		Encrypter: encrypter,
		Generator: generators.NewBase32Generator(15),
		CookieOptions: auth.CookieOptions{
			Name:   "session_id",
			Path:   "/",
			Secure: false,
		},
	})

	http.HandleFunc("/signin", func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			w.Write([]byte(fmt.Sprintf("Error reading body: %s", err.Error())))

			return
		}

		_, session, err := authManager.CreateNewSession(r.Context(), w, &adaptors.Session{
			UserID:       string(body),
			ExpiresAt:    time.Now().Add(time.Hour),
			RefreshUntil: time.Now().Add(time.Hour * 24),
		})
		if err != nil {
			w.Write([]byte(fmt.Sprintf("Error creating session: %s", err.Error())))
			return
		}

		res, err := json.Marshal(map[string]string{
			"sessionId":    session.GetSessionID(),
			"userId":       session.GetUserID(),
			"expiresAt":    session.GetExpiresAt().Format(time.RFC3339),
			"refreshUntil": session.GetRefreshUntil().Format(time.RFC3339),
		})
		if err != nil {
			w.Write([]byte(fmt.Sprintf("Error marshalling session: %s", err.Error())))
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(res)
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, session, err := authManager.GetSessionAndRefresh(r.Context(), w, r, time.Now().Add(time.Hour))
		if err != nil {
			w.Write([]byte(fmt.Sprintf("Error getting session: %s", err.Error())))
			return
		}

		w.Write([]byte(fmt.Sprintf("Session: %#v", session)))
	})

	http.HandleFunc("/signout", func(w http.ResponseWriter, r *http.Request) {
		_, err := authManager.InvalidateSession(r.Context(), w, r)
		if err != nil {
			w.Write([]byte(fmt.Sprintf("Error invalidating session: %s", err.Error())))
			return
		}

		w.Write([]byte("Session invalidated"))
	})

	fmt.Println("Listening on :8080")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error starting server: %s", err.Error())
	}
}
