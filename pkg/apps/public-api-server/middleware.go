package public_api_server

import (
	"2018_2_Stacktivity/pkg/responses"
	"2018_2_Stacktivity/pkg/session"
	"context"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/prometheus/client_golang/prometheus"
)

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("Access-Control-Allow-Origin", config.AllowedIP)
			w.Header().Set("Access-Control-Allow-Methods", config.AllowedMethods)
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			w.Header().Set("X-Vasily", "48")

			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusOK)
				return
			}
			next.ServeHTTP(w, r)
		})
}

// re-implementation of ResponseWriter to access status code
type statusWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *statusWriter) WriteHeader(code int) {
	if code != 0 {
		w.statusCode = code
	} else {
		w.statusCode = 200
	}

	w.ResponseWriter.WriteHeader(code)
}

func (srv *Server) metricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writter := statusWriter{w, 200}
		next.ServeHTTP(&writter, r)
		//log.Println("metric increment")
		ApiMetric.With(prometheus.Labels{"code": strconv.Itoa(writter.statusCode), "path": r.URL.Path, "method": r.Method}).Inc()
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
			println("value: " + value)
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

func getIsAuth(r *http.Request) bool {
	return r.Context().Value("isAuth").(bool)
}

func getUserID(r *http.Request) int32 {
	return r.Context().Value("userID").(int32)
}

func getSessionID(r *http.Request) string {
	return r.Context().Value("sessionID").(string)
}
