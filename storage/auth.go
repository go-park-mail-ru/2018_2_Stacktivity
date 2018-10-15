package storage

import (
	"2018_2_Stacktivity/models"
	"crypto/md5"
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

var (
	ErrNotFound          = errors.New("Username not found")
	ErrIncorrectPassword = errors.New("Incorrect password")
)

func (s *UserStorage) Login(username string, password string) (models.User, error) {
	user, has, err := s.GetByUsername(username)
	if err != nil {
		return user, errors.Wrap(err, "can't get user by username")
	}
	if !has {
		return user, ErrNotFound
	}
	if !checkPassword(password, user.Password) {
		return user, ErrIncorrectPassword
	}
	return user, nil
}

func createPassword(password string) string {
	hash1 := fmt.Sprintf("%x", md5.Sum([]byte("key"+password)))
	hash2 := fmt.Sprint("%x", md5.Sum([]byte(hash1)))
	return fmt.Sprint("%x", md5.Sum([]byte(hash2)))
}

func checkPassword(password string, correct string) bool {
	hash1 := fmt.Sprintf("%x", md5.Sum([]byte("key"+password)))
	hash2 := fmt.Sprint("%x", md5.Sum([]byte(hash1)))
	hash3 := fmt.Sprint("%x", md5.Sum([]byte(hash2)))
	return strings.EqualFold(hash3, correct)
}
