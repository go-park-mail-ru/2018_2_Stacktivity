package server

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
	RedisAddr      string
	APIPrefix      string
	AllowedIP      string
	AllowedMethods string
	PathToOpenAPI  string
}

var config flags

func init() {
	flag.StringVar(&config.Name, "project name", "server", "set name of project")
	flag.StringVar(&config.Port, "port", ":3000", "service port")
	flag.DurationVar(&config.WriteTimeout, "write timeout", 15*time.Second, "timeout for write")
	flag.DurationVar(&config.ReadTimeout, "read timeout", 15*time.Second, "timeout for read")
	flag.StringVar(&config.DB, "database DSN", "user=postgres password=postgres dbname=postgres sslmode=disable", "DSN for database")
	flag.StringVar(&config.RedisAddr, "redis addres", "redis://user:@localhost:6379/0", "redis addr")
	flag.StringVar(&config.APIPrefix, "prefix URL for API", "/api/v1", "URL for requests for API")
	flag.StringVar(&config.AllowedIP, "allowed IP", "http://212.109.223.57:10002", "IP for CORS")
	flag.StringVar(&config.AllowedMethods, "allowed HTTP methods", "POST, GET, PUT, DELETE, OPTIONS", "HTTP methids for CORS")
	flag.StringVar(&config.PathToOpenAPI, "path to OpenAPI", "$GOPATH/src/2018_2_Stacktivity/docs/swagger/html-client", "path for returning OpenAPI")
}
