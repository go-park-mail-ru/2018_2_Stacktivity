package server

import (
	"net/http"

	"log"

	"github.com/gorilla/websocket"
)

func CreateRoom(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (srv *Server) CreatePlayer(w http.ResponseWriter, r *http.Request) {
	println("func CreatePlayer")
	user := getUser(r)
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("can't upgrade connection: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	srv.game.AddPlayer(&user, conn)

}

func GetRoom(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
