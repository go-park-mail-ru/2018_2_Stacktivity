package public_api_server

import (
	"2018_2_Stacktivity/models"
	"2018_2_Stacktivity/pkg/responses"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func (srv *Server) createUser(w http.ResponseWriter, r *http.Request) {
	if getIsAuth(r) {
		responses.Write(w, http.StatusBadRequest, "User alredy signed in")
		return
	}
	var registrationReq models.Registration
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		srv.log.Warnln("can't read request from body", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err = registrationReq.UnmarshalJSON(body); err != nil {
		srv.log.Warnln("can't unmarshal request", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err = srv.validate.Struct(&registrationReq); err != nil {
		responses.Write(w, http.StatusBadRequest, err.Error())
		return
	}
	newUser := models.NewUser(registrationReq.Username, registrationReq.Email, registrationReq.Password)
	usernameExist, emailExist, err := srv.users.CheckExists(newUser)
	if err != nil {
		srv.log.Warnln("can't check user exist", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if usernameExist || emailExist {
		responses.Write(w, http.StatusConflict, "User alredy exists")
		return
	}
	if err = srv.users.Add(&newUser); err != nil {
		srv.log.Warnln("can't insert user into db", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	responses.Write(w, http.StatusCreated, newUser)
}

func (srv *Server) getUser(w http.ResponseWriter, r *http.Request) {
	if !getIsAuth(r) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		responses.Write(w, http.StatusBadRequest, "Incorrect path")
		return
	}
	if int32(id) != getUserID(r) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	user, has, err := srv.users.GetByID(int32(id))
	if err != nil {
		srv.log.Warnln("can't get user", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if !has {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	responses.Write(w, http.StatusOK, user)
}

func (srv *Server) updateUser(w http.ResponseWriter, r *http.Request) {
	if !getIsAuth(r) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	var updateReq models.UserUpdate
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		srv.log.Warnln("can't read request from body", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err = updateReq.UnmarshalJSON(body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	updatedUser, err := srv.users.UpdateUser(getUserID(r), updateReq)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	responses.Write(w, http.StatusOK, updatedUser)
}

func (srv *Server) getUsers(w http.ResponseWriter, r *http.Request) {
	users, err := srv.users.GetAll()
	if err != nil {
		srv.log.Warnln("can't get all users", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	responses.Write(w, http.StatusOK, users)
}

func (srv *Server) GetUsersWithOptions(w http.ResponseWriter, r *http.Request) {
	limit, err := strconv.Atoi(r.FormValue("limit"))
	if err != nil {
		responses.Write(w, http.StatusBadRequest, "Incorrect query")
		return
	}
	offset, err := strconv.Atoi(r.FormValue("offset"))
	if err != nil {
		responses.Write(w, http.StatusBadRequest, "Incorrect query")
		return
	}
	users, err := srv.users.GetAllWithOptions(limit, offset)
	if err != nil {
		srv.log.Warnln("can't get users with options")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	responses.Write(w, http.StatusOK, users)
}
