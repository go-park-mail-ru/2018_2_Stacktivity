package session_server

import (
	"github.com/garyburd/redigo/redis"
	"github.com/pkg/errors"
)

var db *redis.Conn

func InitRedis(redisAddr string) error {
	newDb, err := redis.DialURL(redisAddr)
	if err != nil {
		return errors.Wrap(err, "can't init redis")
	}
	db = &newDb
	return nil
}

func GetInstanse() *redis.Conn {
	return db
}
