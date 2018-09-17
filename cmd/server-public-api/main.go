package main

import (
	"flag"
	"net/http"
	"os"

	"github.com/sirupsen/logrus"
)

func main() {
	flag.Parse()
	log := logrus.New()
	log.SetLevel(logrus.InfoLevel)
	log.SetOutput(os.Stdout)

	log.Infof("starting server listening on %s", config.Port)
	http.ListenAndServe(config.Port, Route())
}

func Route() http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir(config.DirWithStatic)))
	return mux
}
