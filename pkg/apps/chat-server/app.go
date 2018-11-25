package chat_server

import (
	"2018_2_Stacktivity/pkg/session"
	"2018_2_Stacktivity/storage"
	"context"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type Server struct {
	httpSrv *http.Server
	sm      session.SessionManagerClient
	users   storage.UserStorageI
	cm      *ChatManager
	log     *log.Logger
}

func newServer(logger *log.Logger, sessionConn *grpc.ClientConn) *Server {
	return &Server{
		httpSrv: &http.Server{
			Addr:         config.Port,
			WriteTimeout: config.WriteTimeout,
			ReadTimeout:  config.ReadTimeout,
		},
		cm:    NewChatManager(logger),
		sm:    session.NewSessionManagerClient(sessionConn),
		users: storage.GetUserStorage(),
		log:   logger,
	}
}

func (srv *Server) createRoute() {
	r := mux.NewRouter()
	r.Use(srv.loggingMiddleware)
	r.Use(corsMiddleware)
	r.Use(srv.authMiddleware)
	r.Use(srv.checkAuthorization)

	r.HandleFunc("/chat", srv.AddConnection)
	srv.httpSrv.Handler = r
}

func StartApp() {
	flag.Parse()
	logger := log.New()
	logger.SetLevel(log.InfoLevel)
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

	srv := newServer(logger, sessionConn)
	srv.createRoute()
	go srv.cm.Run()
	go func() {
		logger.Infof("Starting chat-server on %s", config.Port)
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
