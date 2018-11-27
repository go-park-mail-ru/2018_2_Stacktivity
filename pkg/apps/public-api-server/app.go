package public_api_server

import (
	"2018_2_Stacktivity/models"
	"2018_2_Stacktivity/pkg/session"
	"2018_2_Stacktivity/storage"
	"2018_2_Stacktivity/storage/migrations"
	"context"
	"flag"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"gopkg.in/go-playground/validator.v9"
)

type Server struct {
	httpSrv *http.Server
	sm      session.SessionManagerClient
	users   storage.UserStorageI

	validate *validator.Validate
	log      *log.Logger
}

func newServer(logger *log.Logger, sessionConn *grpc.ClientConn) *Server {
	return &Server{
		httpSrv: &http.Server{
			Addr:         config.Port,
			WriteTimeout: config.WriteTimeout,
			ReadTimeout:  config.ReadTimeout,
		},
		sm:       session.NewSessionManagerClient(sessionConn),
		users:    storage.GetUserStorage(),
		validate: models.InitValidator(),
		log:      logger,
	}
}

var (
	scoreboardHitMetric = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "public_api_scoreboard_handler_requests_total",
			Help: "Total number of requests by HTTP status code and method.",
		},
		[]string{"code", "method"}, // "method" label is for OPTIONS
	)

	createUserHitMetric = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "public_api_create_user_handler_requests_total",
			Help: "Total number of requests by HTTP status code and method",
		},
		[]string{"code", "method"}, // "method" label is for OPTIONS
	)

	getUserHitMetric = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "public_api_get_user_handler_requests_total",
			Help: "Total number of requests by HTTP status code and method",
		},
		[]string{"code", "method"}, // "method" label is for OPTIONS
	)

	updateUserHitMetric = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "public_api_update_user_handler_requests_total",
			Help: "Total number of requests by HTTP status code and method.",
		},
		[]string{"code", "method"}, // "method" label is for OPTIONS
	)

	getAllUsersHitMetric = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "public_api_get_all_users_handler_requests_total",
			Help: "Total number of requests by HTTP status code and method.",
		},
		[]string{"code", "method"}, // "method" label is for OPTIONS
	)

	sessionHitMetrics = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "public_api_session_handler_requests_total",
			Help: "Total number of requests by HTTP status code and method.",
		},
		[]string{"code", "method"},
	)
)

func getBoundCounterMetricMiddleware(counter *prometheus.CounterVec) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return promhttp.InstrumentHandlerCounter(counter, next)
	}
}

func (srv *Server) createRoute() {
	prometheus.MustRegister(scoreboardHitMetric, createUserHitMetric, getUserHitMetric, updateUserHitMetric, getAllUsersHitMetric, sessionHitMetrics)

	r := mux.NewRouter()
	r.Use(srv.logginigMiddleware)
	r.Use(corsMiddleware)
	r.Use(srv.authMiddleware)

	// route for OPTIONS
	//r.HandleFunc("/", srv.getSession).Methods(http.MethodOptions)

	userRouter := r.PathPrefix("/user").Subrouter()

	// GetScoreboard
	userRouter.HandleFunc("", promhttp.InstrumentHandlerCounter(
		scoreboardHitMetric,
		http.HandlerFunc(srv.GetUsersWithOptions),
	)).Methods(http.MethodGet).
		Queries("limit", "{limit:[0-9]*?}", "offset", "{offset:[0-9]*?}")

	// Create/Get User
	userRouter.HandleFunc("", promhttp.InstrumentHandlerCounter(
		createUserHitMetric,
		http.HandlerFunc(srv.createUser),
	)).Methods(http.MethodPost, http.MethodOptions)

	userRouter.HandleFunc("/{id:[0-9]+}", promhttp.InstrumentHandlerCounter(
		getUserHitMetric,
		http.HandlerFunc(srv.getUser),
	)).Methods(http.MethodGet)

	// UpdateUser
	userRouter.HandleFunc("/{id:[0-9]+}", promhttp.InstrumentHandlerCounter(
		updateUserHitMetric,
		http.HandlerFunc(srv.updateUser),
	)).Methods(http.MethodPatch, http.MethodOptions)

	// GetAllUsers
	userRouter.HandleFunc("", promhttp.InstrumentHandlerCounter(
		getAllUsersHitMetric,
		http.HandlerFunc(srv.getUsers),
	)).Methods(http.MethodGet)

	sessionRouter := r.PathPrefix("/session").Subrouter()
	sessionRouter.Use(getBoundCounterMetricMiddleware(sessionHitMetrics))

	// Create/Get/Delete Session
	sessionRouter.HandleFunc("", srv.createSession).Methods(http.MethodPost, http.MethodOptions)
	sessionRouter.HandleFunc("", srv.getSession).Methods(http.MethodGet)
	sessionRouter.HandleFunc("", srv.deleteSession).Methods(http.MethodDelete, http.MethodOptions)

	// Prometheus handler
	r.Handle("/metrics", promhttp.Handler())

	srv.httpSrv.Handler = r
}

var ()

func StartApp() {

	flag.Parse()
	logger := log.New()
	logger.SetLevel(log.InfoLevel)
	// TODO add hook for logrus
	logger.SetOutput(os.Stdout)
	err := storage.InitDB(config.DB)
	if err != nil {
		log.Warnln("can't init database", err.Error())
		return
	}
	sessionConn, err := grpc.Dial(config.SessionAddr, grpc.WithInsecure())
	if err != nil {
		log.Warnln("can't connect to grpc")
		return
	}
	defer sessionConn.Close()

	migrations.InitMigration()
	srv := newServer(logger, sessionConn)
	srv.createRoute()

	go func() {
		logger.Infof("Starting public-api-server on %s", config.Port)
		if err := srv.httpSrv.ListenAndServe(); err != nil {
			log.Warnln(err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	<-c
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	srv.httpSrv.Shutdown(ctx)
	log.Infoln("Shutdown public-api-server...")
	os.Exit(0)
}
