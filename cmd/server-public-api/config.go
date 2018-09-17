package main

import (
	"flag"
)

type flags struct {
	Name          string
	Port          string
	DirWithStatic string
}

var config flags

func init() {
	flag.StringVar(&config.Name, "project name", "server-public-api", "set name of project")
	flag.StringVar(&config.Port, "port", ":3000", "service port")
	flag.StringVar(&config.DirWithStatic, "public dir", "./public", "path to a folder with static files")

}
