package public_api_server

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
	APIPrefix      string
	AllowedIP      string
	AllowedMethods string
	PathToOpenAPI  string
	SessionAddr    string
}

var config flags

func init() {
	flag.StringVar(&config.Name, "project name", "public-api-server", "set name of project")
	flag.StringVar(&config.Port, "port", ":3001", "service port")
	flag.DurationVar(&config.WriteTimeout, "write timeout", 15*time.Second, "timeout for write")
	flag.DurationVar(&config.ReadTimeout, "read timeout", 15*time.Second, "timeout for read")
	flag.StringVar(&config.DB, "database DSN", "user=docker password=docker dbname=docker sslmode=disable", "DSN for database")
	flag.StringVar(&config.APIPrefix, "prefix URL for API", "/api/v1", "URL for requests for API")
	flag.StringVar(&config.AllowedIP, "allowed IP", "http://blep.me", "IP for CORS")
	flag.StringVar(&config.AllowedMethods, "allowed HTTP methods", "POST, GET, PUT, DELETE, OPTIONS", "HTTP methids for CORS")
	flag.StringVar(&config.PathToOpenAPI, "path to OpenAPI", "$GOPATH/src/2018_2_Stacktivity/docs/swagger/html-client", "path for returning OpenAPI")
	flag.StringVar(&config.SessionAddr, "session addres", "127.0.0.1:8081", "addres of session microservice")
}
