package chat_server

import "net/http"

func (srv *Server) AddConnection(w http.ResponseWriter, r *http.Request) {
	user := getUser(r)
	conn, err := CreateConnection(w, r)
	if err != nil {
		srv.log.Println("can't create connection", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	srv.cm.CreatePlayer(&user, conn)
	w.WriteHeader(http.StatusOK)
}
