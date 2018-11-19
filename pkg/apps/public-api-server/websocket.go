package public_api_server

import (
	"2018_2_Stacktivity/models"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
)

func CreateConnection(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
	u := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	conn, err := u.Upgrade(w, r, nil)
	if err != nil {
		return nil, errors.Wrap(err, "can't upgrade connections")
	}
	m := models.Message{
		Event: models.CreateConn,
	}
	if err := conn.WriteJSON(m); err != nil {
		return nil, errors.Wrap(err, "can't send message to player")
	}
	return conn, nil
}
