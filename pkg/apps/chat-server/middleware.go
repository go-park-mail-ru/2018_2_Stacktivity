package chat_server

import (
	"2018_2_Stacktivity/models"
	"2018_2_Stacktivity/pkg/session"
	"context"
	"net/http"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
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
			s, err := r.Cookie("session-server-id")
			if err == http.ErrNoCookie {
				isAuth = false
			} else {
				id, err = uuid.Parse(s.Value)
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
			ctx = context.WithValue(ctx, "sessionID", id)
			ctx = context.WithValue(ctx, "isAuth", isAuth)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
}

func (srv *Server) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			srv.log.Infoln(r.Method, r.URL.Path)
			next.ServeHTTP(w, r)
		})
}

func (srv *Server) checkAuthorization(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			var user models.User
			var err error
			log.Println("checkAuthorization")
			if !getIsAuth(r) {
				user, _, err = srv.users.GetByUsername("anonimous")
				if err != nil {
					log.Println("can't get user by username", err)
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
			} else {
				var has bool
				id := getUserID(r)
				user, has, err = srv.users.GetByID(id)
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

func getSessionID(r *http.Request) string {
	return r.Context().Value("sessionID").(string)
}

func getUser(r *http.Request) models.User {
	return r.Context().Value("user").(models.User)
}
