package responses

import (
	"encoding/json"
	"net/http"
	"time"
)

type Error struct {
	Message string `json:"message,omitempty"`
}

type UserID struct {
	ID int `json:"ID"`
}

func WriteCookie(w http.ResponseWriter, name string, value string, expires time.Time) {
	cookie := http.Cookie{
		Name:    name,
		Value:   value,
		Expires: time.Now().Add(365 * 24 * time.Hour),
	}
	http.SetCookie(w, &cookie)
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
