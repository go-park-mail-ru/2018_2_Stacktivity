package game_server

import (
	"flag"
	"time"
)

type flags struct {
	Name           string
	Port           string
	WriteTimeout   time.Duration
	ReadTimeout    time.Duration
	DB             string
	AllowedIP      string
	AllowedMethods string
	SessionAddr    string
}

var config flags

func init() {
	flag.StringVar(&config.Name, "project name", "game-server", "set name of project")
	flag.StringVar(&config.Port, "port", ":8083", "service port")
	flag.DurationVar(&config.WriteTimeout, "write timeout", 15*time.Second, "timeout for write")
	flag.DurationVar(&config.ReadTimeout, "read timeout", 15*time.Second, "timeout for read")
	flag.StringVar(&config.DB, "database DSN", "host=postgres user=docker password=docker dbname=docker sslmode=disable", "DSN for database")
	flag.StringVar(&config.AllowedIP, "allowed IP", "http://blep.me", "IP for CORS")
	flag.StringVar(&config.AllowedMethods, "allowed HTTP methods", "POST, GET, PUT, DELETE, OPTIONS", "HTTP methids for CORS")
	flag.StringVar(&config.SessionAddr, "session addres", "session:8081", "addres of session microservice")
}
