package game_server

import (
	"2018_2_Stacktivity/models"
	"2018_2_Stacktivity/pkg/responses"
	"2018_2_Stacktivity/pkg/session"
	"context"
	"log"
	"net/http"

	"github.com/google/uuid"
)

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("Access-Control-Allow-Origin", config.AllowedIP)
			w.Header().Set("Access-Control-Allow-Methods", config.AllowedMethods)
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusOK)
				return
			}
			next.ServeHTTP(w, r)
		})
}

func (srv *Server) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			var isAuth bool
			var sess *session.Session
			var id uuid.UUID
			ctx := r.Context()
			value, err := responses.GetValueFromCookie(r, "sessionID")
			if err == http.ErrNoCookie {
				isAuth = false
			} else {
				id, err = uuid.Parse(value)
				if err != nil {
					srv.log.Warnln("can't parse ")
					ctx = context.WithValue(ctx, "isAuth", false)
					next.ServeHTTP(w, r.WithContext(ctx))
				}
				sess, err = srv.sm.Check(ctx, &session.SessionID{ID: id.String()})
				if err != nil {
					srv.log.Println("can't check session ID: ", err)
				}
				if sess != nil {
					isAuth = true
					ctx = context.WithValue(ctx, "userID", sess.ID)
				}
			}
			ctx = context.WithValue(ctx, "sessionID", id.String())
			ctx = context.WithValue(ctx, "isAuth", isAuth)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
}

func (srv *Server) logginigMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			srv.log.Infoln(r.Method, r.URL.Path)
			next.ServeHTTP(w, r)
		})
}

func (srv *Server) checkAuthorization(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			log.Println("checkAuthorization")
			if !getIsAuth(r) {
				log.Println("user is not auth in game-server")
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			id := getUserID(r)
			user, has, err := srv.users.GetByID(id)
			if err != nil {
				log.Println("can't get user by ID: ", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			if !has {
				log.Printf("user not found by ID: %d\n", id)
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			ctx := context.WithValue(r.Context(), "user", user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
}

func getIsAuth(r *http.Request) bool {
	return r.Context().Value("isAuth").(bool)
}

func getUserID(r *http.Request) int32 {
	return r.Context().Value("userID").(int32)
}

func getUser(r *http.Request) models.User {
	return r.Context().Value("user").(models.User)
}
