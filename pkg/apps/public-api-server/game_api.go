package public_api_server

import (
	"net/http"
)

func (srv *Server) CreateSinglePlayer(w http.ResponseWriter, r *http.Request) {
	user := getUser(r)
	conn, err := CreateConnection(w, r)
	if err != nil {
		srv.log.Println("can't create connection", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	srv.game.RunSinglePlayer(&user, conn)
}

func (srv *Server) CreatePlayer(w http.ResponseWriter, r *http.Request) {
	println("func CreatePlayer")
	user := getUser(r)
	conn, err := CreateConnection(w, r)
	if err != nil {
		srv.log.Println("can't create connection", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	srv.game.AddPlayer(&user, conn)
}

func GetRoom(w http.ResponseWriter, r *http.Request) {
	// TODO connect to game-server from room UID
	w.WriteHeader(http.StatusOK)
}
