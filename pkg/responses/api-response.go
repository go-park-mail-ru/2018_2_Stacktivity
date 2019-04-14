package responses

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/securecookie"
)

type Error struct {
	Message string `json:"message,omitempty"`
}

type UserID struct {
	ID int32 `json:"ID"`
}

var hashKey = []byte("hash-key") // TODO delete this from repository

// Block keys should be 16 bytes (AES-128) or 32 bytes (AES-256) long.
// Shorter keys may weaken the encryption used.
var blockKey = []byte("key-key-key-key-") // TODO delete this from repository
var s = securecookie.New(hashKey, blockKey)

func WriteCookie(w http.ResponseWriter, name string, value string, expires time.Time) {
	log.Println("Write cookie name:", name)
	if encoded, err := s.Encode(name, value); err == nil {
		cookie := http.Cookie{
			Name:     name,
			Value:    encoded,
			Expires:  time.Now().Add(365 * 24 * time.Hour),
			Secure:   true,
			HttpOnly: true,
		}
		http.SetCookie(w, &cookie)
	} else {
		log.Println(err.Error())
	}
}

func GetValueFromCookie(r *http.Request, name string) (string, error) {
	cookie, err := r.Cookie(name)
	if err != nil {
		return "", err
	}
	var value string
	err = s.Decode(name, cookie.Value, &value)
	if err != nil {
		return "", err
	}

	return value, nil
}

func Write(w http.ResponseWriter, statusCode int, response interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(statusCode)
	resp, err := json.Marshal(response)
	if err != nil {
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	_, err = w.Write(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
