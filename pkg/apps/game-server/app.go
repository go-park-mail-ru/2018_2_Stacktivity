package game_server

import (
	. "2018_2_Stacktivity/pkg/apps/game-server/game"
	"2018_2_Stacktivity/pkg/session"
	"2018_2_Stacktivity/storage"
	"context"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type Server struct {
	httpSrv *http.Server
	sm      session.SessionManagerClient
	users   storage.UserStorageI
	game    *Game
	log     *log.Logger
}

func newServer(logger *log.Logger, sessionConn *grpc.ClientConn) *Server {
	return &Server{
		httpSrv: &http.Server{
			Addr:         config.Port,
			WriteTimeout: config.WriteTimeout,
			ReadTimeout:  config.ReadTimeout,
		},
		game:  NewGame(logger),
		sm:    session.NewSessionManagerClient(sessionConn),
		users: storage.GetUserStorage(),
		log:   logger,
	}
}

func (srv *Server) createRoute() {
	r := mux.NewRouter()
	r.Use(srv.logginigMiddleware)
	r.Use(corsMiddleware)
	r.Use(srv.authMiddleware)

	gameRouter := r.PathPrefix("/game").Subrouter()
	// Create/Get/Delete Game
	gameRouter.Use(srv.checkAuthorization)

	gameRouter.HandleFunc("/singleplayer", promhttp.InstrumentHandlerCounter(
		singleplayerHitsMetric,
		http.HandlerFunc(srv.CreateSinglePlayer),
	))

	gameRouter.HandleFunc("/multiplayer", promhttp.InstrumentHandlerCounter(
		multiplayerHitsMetric,
		http.HandlerFunc(srv.CreatePlayer),
	))

	gameRouter.HandleFunc("/game/{id:[0-9]+}", promhttp.InstrumentHandlerCounter(
		gameHitsMetric,
		http.HandlerFunc(GetRoom),
	))

	r.Handle("/metrics", promhttp.Handler())

	srv.httpSrv.Handler = r
}

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

	prometheus.MustRegister(PlayersLeftGameMetric, PlayersPendingRoomMetric, RoomCountMetric, singleplayerHitsMetric, multiplayerHitsMetric, gameHitsMetric)

	srv := newServer(logger, sessionConn)
	srv.createRoute()

	go func() {
		logger.Infof("Starting game-server on %s", config.Port)
		if err := srv.httpSrv.ListenAndServe(); err != nil {
			log.Warnln(err)
		}
	}()
	srv.game.Start()
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	<-c
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	srv.game.Stop()
	if err = srv.httpSrv.Shutdown(ctx); err != nil {
		logger.Warnln("can't shutdown game-server")
	}
	log.Infoln("Shutdown public-api-server...")
	os.Exit(0)
}
