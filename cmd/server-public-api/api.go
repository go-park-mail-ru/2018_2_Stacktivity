package main

import (
	"2018_2_Stacktivity/cmd/server-public-api/requests"
	"2018_2_Stacktivity/cmd/server-public-api/responses"
	"2018_2_Stacktivity/cmd/server-public-api/session"
	"2018_2_Stacktivity/cmd/server-public-api/storage"
	"2018_2_Stacktivity/cmd/server-public-api/validate"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/google/uuid"
)

func (srv *Server) createUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		err := responses.WriteResponse(w, http.StatusBadRequest, &responses.ResponseForm{
			ValidateSuccess: false,
			Error:           responses.NewError("Bad method"),
		})
		if err != nil {
			srv.log.Warnf("error in %s: %s", r.URL.Path, err.Error())
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	var registrationReq requests.Registration
	defer r.Body.Close()
	req, err := ioutil.ReadAll(r.Body)
	if err != nil {
		srv.log.Warnf("error in %s: %s", r.URL.Path, err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = json.Unmarshal(req, &registrationReq)
	if err != nil {
		srv.log.Warnf("error in %s: %s", r.URL.Path, err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	response := validate.RegistrationValidate(&registrationReq)
	if !response.ValidateSuccess {
		err = responses.WriteResponse(w, http.StatusBadRequest, response)
		if err != nil {
			srv.log.Warnf("error in /registration: %s", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		return
	}

	newUser := storage.NewUser(registrationReq.Username, registrationReq.Email, registrationReq.Password1)
	if srv.users.Has(newUser.Username) {
		response.ValidateSuccess = false
		response.Error = &responses.Error{
			Message: "Username already exists",
		}
		err = responses.WriteResponse(w, http.StatusBadRequest, response)
		if err != nil {
			srv.log.Warnf("error in %d: %s", r.URL.Path, err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		return
	}

	_, hasEmail, err := srv.users.GetByEmail(newUser.Email)
	if err != nil {
		srv.log.Warnf("error in %d: %s", r.URL.Path, err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if hasEmail {
		response.ValidateSuccess = false
		response.Error = &responses.Error{
			Message: "Email already exists",
		}
		err = responses.WriteResponse(w, http.StatusBadRequest, response)
		if err != nil {
			srv.log.Warnf("error in %d: %s", r.URL.Path, err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		return
	}

	id, err := srv.users.Add(newUser)
	if err != nil {
		srv.log.Warnln("Can't add user into db", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	newUser.ID = id

	sess, err := srv.sm.Create(&session.Session{
		Username:  registrationReq.Username,
		Useragent: r.UserAgent(),
	})
	if err != nil {
		srv.log.Warnln("Can't create session", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "session-id",
		Value:   sess.ID.String(),
		Expires: time.Now().Add(365 * 24 * time.Hour),
	})

	err = responses.WriteResponse(w, http.StatusOK, &responses.ResponseForm{
		ValidateSuccess: true,
		User:            &newUser,
	})
	if err != nil {
		srv.log.Warnf("error in %s: %s", r.URL.Path, err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (srv *Server) createSession(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		err := responses.WriteResponse(w, http.StatusBadRequest, &responses.ResponseForm{
			ValidateSuccess: false,
			Error:           responses.NewError("Bad method"),
		})
		if err != nil {
			srv.log.Warnf("error in %s: %s", r.URL.Path, err.Error())
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}
	var loginReq requests.Login
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		srv.log.Warnf("error in %s: %s", r.URL.Path, err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err = json.Unmarshal(body, &loginReq); err != nil {
		srv.log.Warnf("error in %s: %s", r.URL.Path, err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	response := validate.LoginValidate(&loginReq)
	if !response.ValidateSuccess {
		if err = responses.WriteResponse(w, http.StatusBadRequest, response); err != nil {
			srv.log.Warnf("error in %s: %s", r.URL.Path, err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		return
	}

	user, find, err := srv.users.GetByUsername(loginReq.Username)
	if !find || err != nil {
		response.ValidateSuccess = false
		response.Error = &responses.Error{
			Message: "Incorrect login or password",
		}
		if err = responses.WriteResponse(w, http.StatusBadRequest, response); err != nil {
			srv.log.Warnf("error in %s: %s", r.URL.Path, err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		return
	}
	if !storage.CheckPassword(loginReq.Password, user.Password) {
		response.ValidateSuccess = false
		response.Error = &responses.Error{
			Message: "Incorrect login or password",
		}
		if err = responses.WriteResponse(w, http.StatusBadRequest, response); err != nil {
			srv.log.Warnf("error in %s: %s", r.URL.Path, err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		return
	}
	sess, err := srv.sm.Create(&session.Session{
		Username: loginReq.Username,
	})
	if err != nil {
		srv.log.Infoln("Can't create session")
		srv.log.Warnf("error in %s: %s", r.URL.Path, err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	cookie := http.Cookie{
		Name:    "session-id",
		Value:   sess.ID.String(),
		Expires: time.Now().Add(365 * 24 * time.Hour),
	}
	http.SetCookie(w, &cookie)
	response.User = &user
	if err = responses.WriteResponse(w, http.StatusOK, response); err != nil {
		srv.log.Warnf("error in %s: %s", r.URL.Path, err.Error())
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (srv *Server) deleteSession(w http.ResponseWriter, r *http.Request) {
	s, err := r.Cookie("session-id")
	if err == http.ErrNoCookie {
		w.WriteHeader(http.StatusOK)
		return
	}
	if err != nil {
		srv.log.Warnln(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	id, err := uuid.Parse(s.Value)
	if err != nil {
		srv.log.Warnln(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	sessionID := &session.SessionID{
		ID: id,
	}
	srv.sm.Delete(sessionID)
	w.WriteHeader(http.StatusOK)
}

func (srv *Server) getSession(w http.ResponseWriter, r *http.Request) {
	s, err := r.Cookie("session-id")
	if err == http.ErrNoCookie {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	id, err := uuid.Parse(s.Value)
	if err != nil {
		srv.log.Warnf("error in %s: %s", r.URL.Path, err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	sess := srv.sm.Check(&session.SessionID{
		ID: id,
	})
	if sess == nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	user, find, err := srv.users.GetByUsername(sess.Username)
	if err != nil {
		srv.log.Warnf("error in %s: %s", r.URL.Path, err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if !find {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if err = responses.WriteResponse(w, http.StatusOK, user); err != nil {
		srv.log.Warnf("error in %s: %s", r.URL.Path, err.Error())
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (srv *Server) getUser(w http.ResponseWriter, r *http.Request) {
	var head string
	head, r.URL.Path = ShiftPath(r.URL.Path)
	id, err := strconv.Atoi(head)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	user, find, err := srv.users.GetByID(id)
	if err != nil {
		srv.log.Warnf("error in %s: %s", r.URL.Path, err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if !find {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if err = responses.WriteResponse(w, http.StatusOK, user); err != nil {
		srv.log.Warnf("error in %s: %s", r.URL.Path, err.Error())
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (srv *Server) updateUser(w http.ResponseWriter, r *http.Request) {
	var head string
	head, r.URL.Path = ShiftPath(r.URL.Path)
	ID, err := strconv.Atoi(head)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	s, err := r.Cookie("session-id")
	if err == http.ErrNoCookie {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if err != nil {
		srv.log.Warnln(err.Error())
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	sID, err := uuid.Parse(s.Value)
	if err != nil {
		srv.log.Warnln(err.Error())
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	sess := srv.sm.Check(&session.SessionID{
		ID: sID,
	})
	if sess == nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	userFromCookie, has, err := srv.users.GetByUsername(sess.Username)
	if err != nil {
		srv.log.Warnf("error in %s: %s", r.URL.Path, err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if !has {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	userFromQuery, has, err := srv.users.GetByID(ID)
	if err != nil {
		srv.log.Warnf("error in %s: %s", r.URL.Path, err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if !has {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if userFromCookie.ID != userFromQuery.ID {
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	var updateReq requests.UserUpdate
	defer r.Body.Close()
	req, err := ioutil.ReadAll(r.Body)
	if err != nil {
		srv.log.Warnf("error in %s: %s", r.URL.Path, err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = json.Unmarshal(req, &updateReq)
	if err != nil {
		srv.log.Warnf("error in %s: %s", r.URL.Path, err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	response := validate.UpdateValidate(&updateReq)
	if !response.ValidateSuccess {
		err = responses.WriteResponse(w, http.StatusBadRequest, response)
		if err != nil {
			srv.log.Warnf("error: %s", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		return
	}
	has = srv.users.Has(updateReq.Username)
	if has {
		response.ValidateSuccess = false
		response.Error = responses.NewError("Username alredy exists")
		err = responses.WriteResponse(w, http.StatusBadRequest, response)
		if err != nil {
			srv.log.Warnf("error: %s", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		return
	}
	_, has, err = srv.users.GetByEmail(updateReq.Email)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if has {
		response.ValidateSuccess = false
		response.Error = responses.NewError("Email alredy exists")
		err = responses.WriteResponse(w, http.StatusBadRequest, response)
		if err != nil {
			srv.log.Warnf("error: %s", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		return
	}
	resUser := storage.User{
		ID:       userFromQuery.ID,
		Username: updateReq.Username,
		Email:    updateReq.Email,
		Password: userFromQuery.Password,
	}
	err = srv.users.Update(resUser)
	if err != nil {
		srv.log.Warnln("Can't update user into db", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
	}
	response.User = &resUser
	err = responses.WriteResponse(w, http.StatusOK, response)
	if err != nil {
		srv.log.Warnf("error: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
	}
}
func (srv *Server) getUsers(w http.ResponseWriter, r *http.Request) {
	var err error
	var users []storage.User
	pageNum := -1
	page, ok := r.URL.Query()["page"]
	if ok {
		pageNum, err = strconv.Atoi(page[0])
		if err != nil {
			pageNum = -1
		}
	}
	if pageNum <= 0 {
		users, err = srv.users.GetAll()
	} else {
		users, err = srv.users.GetWithOptions(config.PageSize, (pageNum-1)*config.PageSize)
	}

	if err != nil {
		srv.log.Warnf("error in %s: %s", r.URL.Path, err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	sort.Sort(storage.Users(users))
	err = responses.WriteResponse(w, http.StatusOK, users)
	if err != nil {
		srv.log.Warnf("error in %s: %s", r.URL.Path, err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
