package server

import (
	"2018_2_Stacktivity/models"
	"2018_2_Stacktivity/pkg/apps/game"
	"2018_2_Stacktivity/session"
	"2018_2_Stacktivity/storage"
	"2018_2_Stacktivity/storage/migrations"
	"context"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"gopkg.in/go-playground/validator.v9"
)

type Server struct {
	httpSrv  *http.Server
	sm       session.SessionManagerI
	users    storage.UserStorageI
	game     *game.Game
	validate *validator.Validate
	log      *log.Logger
}

func newServer(logger *log.Logger) *Server {
	return &Server{
		httpSrv: &http.Server{
			Addr:         config.Port,
			WriteTimeout: config.WriteTimeout,
			ReadTimeout:  config.ReadTimeout,
		},
		sm:       session.NewSessionManager(),
		users:    storage.GetUserStorage(),
		validate: models.InitValidator(),
		game:     game.NewGame(logger),
		log:      logger,
	}
}

func (srv *Server) createRoute() {
	r := mux.NewRouter()
	r.Use(srv.logginigMiddleware)
	r.Use(corsMiddleware)
	r.Use(srv.authMiddleware)

	// route for OPTIONS
	r.HandleFunc("/", srv.getSession).Methods(http.MethodOptions)

	userRouter := r.PathPrefix("/user").Subrouter()

	// GetScoreboard
	userRouter.HandleFunc("", srv.GetUsersWithOptions).Methods(http.MethodGet).
		Queries("limit", "{limit:[0-9]*?}", "offset", "{offset:[0-9]*?}")

	// Create/Get User
	userRouter.HandleFunc("", srv.createUser).Methods(http.MethodPost)
	userRouter.HandleFunc("/{id:[0-9]+}", srv.getUser).Methods(http.MethodGet)

	// UpdateUser
	userRouter.HandleFunc("/{id:[0-9]+}", srv.updateUser).Methods(http.MethodPatch)

	// GetAllUsers
	userRouter.HandleFunc("", srv.getUsers).Methods(http.MethodGet)

	sessionRouter := r.PathPrefix("/session").Subrouter()
	// Create/Get/Delete Session
	sessionRouter.HandleFunc("", srv.createSession).Methods(http.MethodPost)
	sessionRouter.HandleFunc("", srv.getSession).Methods(http.MethodGet)
	sessionRouter.HandleFunc("", srv.deleteSession).Methods(http.MethodDelete)
	srv.httpSrv.Handler = r

	gameRouter := r.PathPrefix("/game").Subrouter()
	// Create/Get/Delete Game
	gameRouter.Use(srv.checkAuthorization)
	gameRouter.HandleFunc("/singleplayer", CreateRoom)
	gameRouter.HandleFunc("/multiplayer", srv.CreatePlayer)
	gameRouter.HandleFunc("/room/{id:[0-9]+}", GetRoom)
}

func StartApp() {
	flag.Parse()
	logger := log.New()
	logger.SetLevel(log.InfoLevel)
	// TODO add hook to Graylog
	logger.SetOutput(os.Stdout)
	err := storage.InitDB(config.DB)
	if err != nil {
		log.Warnln("can't init database", err.Error())
	}
	migrations.InitMigration()
	srv := newServer(logger)
	srv.createRoute()
	srv.game.Start()
	go func() {
		logger.Infof("Starting server on %s", config.Port)
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
	log.Infoln("Shutdown server...")
	os.Exit(0)
}
