package public_api_server

import (
	"2018_2_Stacktivity/models"
	"2018_2_Stacktivity/pkg/responses"
	"2018_2_Stacktivity/pkg/session"
	"2018_2_Stacktivity/storage"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
)

func (srv *Server) createSession(w http.ResponseWriter, r *http.Request) {
	if getIsAuth(r) {
		responses.Write(w, http.StatusBadRequest, "User alredy signed in")
		return
	}
	var loginReq models.Login
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		srv.log.Warnln("can't read request from body", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err = json.Unmarshal(body, &loginReq); err != nil {
		srv.log.Warnln("can't unmarshal request", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err = srv.validate.Struct(loginReq); err != nil {
		responses.Write(w, http.StatusBadRequest, err.Error())
		return
	}
	user, err := srv.users.Login(loginReq.Username, loginReq.Password)
	if err != nil {
		if err == storage.ErrNotFound || err == storage.ErrIncorrectPassword {
			responses.Write(w, http.StatusBadRequest, "Incorrect login or password")
		} else {
			srv.log.Warnln("can't login user", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}
	sess, err := srv.sm.Create(r.Context(), &session.Session{
		ID: user.ID,
	})
	if err != nil {
		srv.log.Warnln("can't create session-server", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	responses.WriteCookie(w, "session-server-id", sess.ID, time.Now().Add(7*24*time.Hour))
	responses.Write(w, http.StatusCreated, user)
}

func (srv *Server) getSession(w http.ResponseWriter, r *http.Request) {
	if !getIsAuth(r) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	responses.Write(w, http.StatusOK, responses.UserID{
		ID: getUserID(r),
	})
}

func (srv *Server) deleteSession(w http.ResponseWriter, r *http.Request) {
	if !getIsAuth(r) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	srv.sm.Delete(r.Context(), &session.SessionID{
		ID: getSessionID(r),
	})
	w.WriteHeader(http.StatusOK)
}
