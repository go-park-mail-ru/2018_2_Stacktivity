package storage

import (
	"crypto/md5"
	"database/sql"
	"fmt"
	"log"
	"strings"

	_ "github.com/lib/pq"
	"github.com/pkg/errors"
)

type UserStorage struct {
	DB    *sql.DB
	stmts *statements
}

type statements struct {
	addUser        *sql.Stmt
	updateUser     *sql.Stmt
	getByEmail     *sql.Stmt
	getByUsername  *sql.Stmt
	getByID        *sql.Stmt
	getAll         *sql.Stmt
	getWithOptions *sql.Stmt
}

func NewUserStorage(db *sql.DB) *UserStorage {
	storage := &UserStorage{
		DB:    db,
		stmts: &statements{},
	}
	return storage
}

func (s *UserStorage) Prepare() error {
	var err error
	s.stmts.addUser, err = s.DB.Prepare(addUser)
	if err != nil {
		return errors.Wrap(err, "can't prepare statements")
	}
	s.stmts.updateUser, err = s.DB.Prepare(updateUser)
	if err != nil {
		return errors.Wrap(err, "can't prepare statements")
	}
	s.stmts.getByID, err = s.DB.Prepare(getByID)
	if err != nil {
		return errors.Wrap(err, "can't prepare statements")
	}
	s.stmts.getByUsername, err = s.DB.Prepare(getByUsername)
	if err != nil {
		return errors.Wrap(err, "can't prepare statements")
	}
	s.stmts.getByEmail, err = s.DB.Prepare(getByEmail)
	if err != nil {
		return errors.Wrap(err, "can't prepare statements")
	}
	s.stmts.getAll, err = s.DB.Prepare(getAll)
	if err != nil {
		return errors.Wrap(err, "can't prepare statements")
	}
	s.stmts.getWithOptions, err = s.DB.Prepare(getWithOptions)
	if err != nil {
		return errors.Wrap(err, "can't prepare statements")
	}
	return nil
}

var addUser = `INSERT INTO users (username, email, pass) VALUES ($1, $2, $3) RETURNING uID;`

func (s *UserStorage) Add(user User) (uid int, err error) {
	err = s.stmts.addUser.QueryRow(user.Username, user.Email, createPassword(user.Password)).Scan(&uid)
	if err != nil {
		err = errors.Wrap(err, "failed to query database")
		return 0, err
	}
	return uid, nil
}

var updateUser = `UPDATE users SET username = $2, email = $3, pass = $4 WHERE uID = $1;`

func (s *UserStorage) Update(user User) error {
	_, err := s.stmts.updateUser.Exec(user.ID, user.Username, user.Email, createPassword(user.Password))
	if err != nil {
		err = errors.Wrap(err, "failed to query database")
		return err
	}
	return nil
}

var getByEmail = `SELECT uid, username, email, pass, score FROM users WHERE users.email = $1 LIMIT 1;`

func (s *UserStorage) GetByEmail(username string) (User, bool, error) {
	var uid, score int
	var name, email, password string
	err := s.stmts.getByEmail.QueryRow(username).Scan(&uid, &name, &email, &password, &score)
	if err == sql.ErrNoRows {
		return User{}, false, nil
	}
	if err != nil {
		return User{}, false, errors.Wrap(err, "failed to query database")
	}
	resUser := User{
		ID:       uid,
		Username: name,
		Email:    email,
		Password: password,
		Score:    score,
	}
	return resUser, true, nil
}

var getByUsername = `SELECT uid, username, email, pass, score FROM users WHERE users.username = $1 LIMIT 1;`

func (s *UserStorage) GetByUsername(username string) (User, bool, error) {
	var uid, score int
	var name, email, password string
	err := s.stmts.getByUsername.QueryRow(username).Scan(&uid, &name, &email, &password, &score)
	if err == sql.ErrNoRows {
		return User{}, false, nil
	}
	if err != nil {
		return User{}, false, errors.Wrap(err, "failed to query database")
	}
	resUser := User{
		ID:       uid,
		Username: name,
		Email:    email,
		Password: password,
		Score:    score,
	}
	return resUser, true, nil
}

var getByID = `SELECT uid, username, email, pass, score FROM users WHERE users.uid = $1 LIMIT 1;`

func (s *UserStorage) GetByID(id int) (User, bool, error) {
	var uID, score int
	var name, email, password string
	err := s.stmts.getByID.QueryRow(id).Scan(&uID, &name, &email, &password, &score)
	if err == sql.ErrNoRows {
		return User{}, false, nil
	}
	if err != nil {
		return User{}, false, errors.Wrap(err, "failed to query getById")
	}
	resUser := User{
		ID:       uID,
		Username: name,
		Email:    email,
		Password: password,
		Score:    score,
	}
	return resUser, true, nil
}

func (s *UserStorage) Has(username string) bool {
	_, find, _ := s.GetByUsername(username)
	return find
}

var getAll = `SELECT uid, username, email, score FROM users;`

func (s *UserStorage) GetAll() ([]User, error) {
	users := make([]User, 0)
	rows, err := s.stmts.getAll.Query()
	if err != nil {
		return users, errors.Wrap(err, "failed to query database")
	}
	for rows.Next() {
		var id, score int
		var username, email string
		if err := rows.Scan(&id, &username, &email, &score); err != nil {
			log.Fatal(err)
		}
		users = append(users, User{
			ID:       id,
			Username: username,
			Email:    email,
			Score:    score,
		})
	}
	return users, nil
}

var getWithOptions = `SELECT uid, username, email, score FROM users LIMIT $1 OFFSET $2;`

func (s *UserStorage) GetWithOptions(limit int, offset int) ([]User, error) {
	users := make([]User, 0)
	rows, err := s.stmts.getWithOptions.Query(limit, offset)
	if err != nil {
		return users, errors.Wrap(err, "failed to query database")
	}
	for rows.Next() {
		var id, score int
		var username, email string
		if err := rows.Scan(&id, &username, &email, &score); err != nil {
			log.Fatal(err)
		}
		users = append(users, User{
			ID:       id,
			Username: username,
			Email:    email,
			Score:    score,
		})
	}
	return users, nil
}

type UserStorageI interface {
	Add(User) (int, error)
	Update(User) error
	GetByEmail(string) (User, bool, error)
	GetByUsername(string) (User, bool, error)
	GetByID(int) (User, bool, error)
	Has(string) bool
	GetAll() ([]User, error)
	GetWithOptions(int, int) ([]User, error)
	Prepare() error
}

func NewUser(username string, email string, pswd string) User {
	return User{
		Username: username,
		Email:    email,
		Password: pswd,
		Score:    0,
	}
}

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"-"`
	Avatar   string `json:"avatar,omitempty"`
	Score    int    `json:"score"`
}

type Users []User

func (u Users) Len() int {
	return len(u)
}

func (u Users) Less(i, j int) bool {
	return u[i].Score < u[j].Score
}

func (u Users) Swap(i, j int) {
	u[i], u[j] = u[j], u[i]
}

func createPassword(password string) string {
	hash1 := fmt.Sprintf("%x", md5.Sum([]byte("key"+password)))
	hash2 := fmt.Sprint("%x", md5.Sum([]byte(hash1)))
	return fmt.Sprint("%x", md5.Sum([]byte(hash2)))
}

func CheckPassword(password string, correct string) bool {
	hash1 := fmt.Sprintf("%x", md5.Sum([]byte("key"+password)))
	hash2 := fmt.Sprint("%x", md5.Sum([]byte(hash1)))
	hash3 := fmt.Sprint("%x", md5.Sum([]byte(hash2)))
	return strings.EqualFold(hash3, correct)
}
