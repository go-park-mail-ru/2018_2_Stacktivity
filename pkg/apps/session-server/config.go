package session_server

import (
	"flag"
	"time"
)

type flags struct {
	Name         string
	Port         string
	WriteTimeout time.Duration
	ReadTimeout  time.Duration
	RedisAddr    string
}

var config flags

func init() {
	flag.StringVar(&config.Name, "project name", "session-server", "set name of project")
	flag.StringVar(&config.Port, "port", ":8081", "service port")
	flag.DurationVar(&config.WriteTimeout, "write timeout", 15*time.Second, "timeout for write")
	flag.DurationVar(&config.ReadTimeout, "read timeout", 15*time.Second, "timeout for read")
	flag.StringVar(&config.RedisAddr, "redis addres", "redis://user:@redis:6379/0", "redis addr")
}
