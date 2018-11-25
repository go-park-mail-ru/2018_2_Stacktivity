package storage

import (
	"2018_2_Stacktivity/models"
	"database/sql"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
)

type UserStorageI interface {
	Add(*models.User) error

	GetAll() ([]models.User, error)
	GetAllWithOptions(limit int, offset int) ([]models.User, error)

	GetByID(id int32) (models.User, bool, error)
	GetByEmail(email string) (models.User, bool, error)
	GetByUsername(username string) (models.User, bool, error)

	UpdateUser(uid int32, update models.UserUpdate) (models.User, error)
	UpdateScore(uid int, newScore int) error

	CheckExists(models.User) (usernameExist bool, emailExist bool, err error)
	Login(username string, password string) (models.User, error)
}

type UserStorage struct {
	DB *sqlx.DB
}

func GetUserStorage() *UserStorage {
	storage := &UserStorage{
		DB: db,
	}
	return storage
}

var addUser = `INSERT INTO "user" (username, email, pass) VALUES ($1, $2, $3) RETURNING uID;`

func (s *UserStorage) Add(user *models.User) error {
	if err := s.DB.QueryRow(addUser, user.Username, user.Email, createPassword(user.Password)).Scan(&user.ID); err != nil {
		return errors.Wrap(err, "can't insert user into db")
	}
	return nil
}

var getAll = `SELECT * FROM "user";`

func (s *UserStorage) GetAll() (users []models.User, err error) { //
	users = []models.User{}
	if err = s.DB.Select(&users, getAll); err != nil {
		return nil, errors.Wrap(err, "failed to query database")
	}
	return users, nil
}

var getAllWithOptions = `SELECT uid, username, email, score FROM "user" ORDER BY score DESC LIMIT $1 OFFSET $2;`

func (s *UserStorage) GetAllWithOptions(limit int, offset int) ([]models.User, error) { //
	users := []models.User{}
	err := s.DB.Select(&users, getAllWithOptions, limit, offset)
	if err != nil {
		return users, errors.Wrap(err, "failed to query database")
	}
	return users, nil
}

var getByID = `SELECT * FROM "user" WHERE uID = $1 LIMIT 1;`

func (s *UserStorage) GetByID(uid int32) (user models.User, has bool, err error) {
	user = models.User{}
	if err = s.DB.Get(&user, getByID, uid); err != nil {
		if err == sql.ErrNoRows {
			return user, false, nil
		}
		return user, false, errors.Wrap(err, "failed to query database")
	}
	return user, true, nil
}

var getByEmail = `SELECT * FROM "user" WHERE email = $1 LIMIT 1;`

func (s *UserStorage) GetByEmail(email string) (user models.User, has bool, err error) {
	user = models.User{}
	if err = s.DB.Get(&user, getByEmail, email); err != nil {
		if err == sql.ErrNoRows {
			return user, false, nil
		}
		return user, false, errors.Wrap(err, "can't select from DB")
	}
	return user, true, nil
}

var getByUsername = `SELECT * FROM "user" WHERE username = $1 LIMIT 1;`

func (s *UserStorage) GetByUsername(username string) (user models.User, has bool, err error) {
	user = models.User{}
	log.Println(username)
	if err = s.DB.Get(&user, getByUsername, username); err != nil {
		if err == sql.ErrNoRows {
			return user, false, nil
		}
		return user, false, errors.Wrap(err, "can't select from DB")
	}
	return user, true, nil
}

var updateUser = `UPDATE "user" SET  = coalesce(coalesce(nullif($2, ''), username)), 
			email = coalesce(coalesce(nullif($3, ''), email)), 
			pass = coalesce(coalesce(nullif($4, ''), pass)) WHERE nickname = $1 RETURNING fullname, email, about;`

func (s *UserStorage) UpdateUser(uid int32, user models.UserUpdate) (models.User, error) {
	var newUser models.User
	return newUser, nil
}

func (s *UserStorage) UpdateUsername(id int, newUsername string) error {
	return s.updateRow(id, "username", newUsername)
}

func (s *UserStorage) UpdateEmail(id int, newEmail string) error {
	return s.updateRow(id, "email", newEmail)
}

func (s *UserStorage) UpdatePassword(id int, newPassword string) error {
	return s.updateRow(id, "password", createPassword(newPassword))
}

func (s *UserStorage) UpdateScore(id int, newScore int) error {
	return s.updateRow(id, "score", newScore)
}

var update = `UPDATE "user" SET $2 = $3 WHERE uID = $1;`

func (s *UserStorage) updateRow(id int, field string, value interface{}) error {
	_, err := s.DB.Exec(update, id, field, value)
	return err
}

func (s *UserStorage) CheckExists(user models.User) (usernameExist bool, emailExist bool, err error) {
	usernameExist, emailExist, err = false, false, nil
	_, usernameExist, err = s.GetByUsername(user.Username)
	if err != nil {
		err = errors.Wrap(err, "can't get user by username")
		return
	}
	_, emailExist, err = s.GetByEmail(user.Email)
	if err != nil {
		err = errors.Wrap(err, "can't get user by email")
	}
	return
}
