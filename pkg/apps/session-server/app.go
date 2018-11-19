package session_server

import (
	"2018_2_Stacktivity/pkg/session"
	"flag"
	"fmt"
	"net"
	"net/http"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type Server struct {
	httpSrv *http.Server
	log     *log.Logger
}

func StartApp() {
	flag.Parse()
	logger := log.New()
	logger.SetLevel(log.InfoLevel)
	lis, err := net.Listen("tcp", config.Port)
	if err != nil {
		logger.Warnln("can't listen port", err)
		return
	}

	server := grpc.NewServer()
	err = InitRedis(config.RedisAddr)
	if err != nil {
		logger.Warnln("can't listen port: ", err)
		return
	}

	session.RegisterSessionManagerServer(server, NewSessionManager(GetInstanse()))
	fmt.Println("starting session-server at :8081")
	server.Serve(lis)
}
