package main

import (
	"flag"
)

type flags struct {
	Name          string
	Port          string
	DB            string
	APIPrefix     string
	AllowedIP     string
	PathToOpenAPI string

	PageSize int
}

var config flags

func init() {
	flag.StringVar(&config.Name, "project name", "server-public-api", "set name of project")
	flag.StringVar(&config.Port, "port", ":3001", "service port")
	flag.StringVar(&config.DB, "database DSN", "user=postgres password=postgres dbname=db-tech sslmode=disable", "DSN for database")
	flag.StringVar(&config.APIPrefix, "prefix URL for API", "/api/v1", "URL for requests for API")
	flag.StringVar(&config.AllowedIP, "allowed IP", "http://212.109.223.57:10002", "IP for CORS")
	flag.StringVar(&config.PathToOpenAPI, "path to OpenAPI", "./swagger/html-client", "path for returning OpenAPI")

	flag.IntVar(&config.PageSize, "size of scroboard page", 5, "set size for user scoreboard pages")
}
