package public_api_server

import (
	"2018_2_Stacktivity/models"
	"2018_2_Stacktivity/pkg/session"
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

type SessionManagerClientMock struct{}

func (c *SessionManagerClientMock) Create(ctx context.Context, in *session.Session, opts ...grpc.CallOption) (*session.SessionID, error) {
	return &session.SessionID{ID: "42"}, nil
}

func (c *SessionManagerClientMock) Check(ctx context.Context, in *session.SessionID, opts ...grpc.CallOption) (*session.Session, error) {
	return &session.Session{}, nil
}

func (c *SessionManagerClientMock) Delete(ctx context.Context, in *session.SessionID, opts ...grpc.CallOption) (*session.Nothing, error) {
	return &session.Nothing{}, nil
}

type UserStorageMock struct{}

func (s *UserStorageMock) Add(*models.User) error {
	return nil
}

func (s *UserStorageMock) GetAll() ([]models.User, error) {
	return []models.User{}, nil
}

func (s *UserStorageMock) GetAllWithOptions(limit int, offset int) ([]models.User, error) {
	return []models.User{}, nil
}

func (s *UserStorageMock) GetByID(id int32) (models.User, bool, error) {
	return models.User{}, true, nil
}

func (s *UserStorageMock) GetByEmail(email string) (models.User, bool, error) {
	return models.User{}, true, nil
}

func (s *UserStorageMock) GetByUsername(username string) (models.User, bool, error) {
	return models.User{}, true, nil
}

func (s *UserStorageMock) GetLevelByNumber(number int) (models.LevelInStorage, error) {
	return models.LevelInStorage{}, nil
}

func (s *UserStorageMock) UpdateUser(uid int32, update models.UserUpdate) (models.User, error) {
	return models.User{}, nil
}

func (s *UserStorageMock) UpdateScore(uid int, newScore int) error {
	return nil
}

func (s *UserStorageMock) UpdateLevel(uid int, newLevel int) error {
	return nil
}

func (s *UserStorageMock) AddScore(id int32, score int) error {
	return nil
}

func (s *UserStorageMock) CheckExists(models.User) (usernameExist bool, emailExist bool, err error) {
	return false, false, nil
}

func (s *UserStorageMock) Login(username string, password string) (models.User, error) {
	return models.User{}, nil
}

func CreateMockServer() *Server {
	logger := logrus.New()
	logger.Out = ioutil.Discard
	return &Server{
		&http.Server{},
		&SessionManagerClientMock{},
		&UserStorageMock{},
		models.InitValidator(),
		logger,
	}
}

func AuthTestHandler() http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if getIsAuth(r) {
			w.Header().Set("IsAuth", "true")
		} else {
			w.Header().Set("IsAuth", "false")
		}
	}
	return http.HandlerFunc(fn)
}

func CORSTestHandler() http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
	}
	return http.HandlerFunc(fn)
}

func TestAuthMiddlewareWithoutCookie(t *testing.T) {
	assert := assert.New(t)
	ts := httptest.NewServer(CreateMockServer().authMiddleware(AuthTestHandler()))

	resp, err := http.Get(ts.URL)
	assert.NoError(err)

	assert.Equal("false", resp.Header.Get("IsAuth"))
}

var secureCookie = "MTU1MTk3MzkwOXxDeHlDcEJYTWk1dkthdmVHSnFYWlF0ZWhqZmxsXzFEZnZaLVRoNTRyZkM2d1phSV9weGRRQVpnc2d4cW1SaGRuM2pKQmowYjdGWTQ9fKLJ-qma-d3u2WPzuxNZ9H5f1hEeCySxGGKpkYordprW"

func TestAuthMiddlewareWithCookie(t *testing.T) {
	assert := assert.New(t)
	ts := httptest.NewServer(CreateMockServer().authMiddleware(AuthTestHandler()))

	req, err := http.NewRequest(http.MethodGet, ts.URL, nil)
	assert.NoError(err)

	req.AddCookie(&http.Cookie{Name: "sessionID", Value: secureCookie})

	client := &http.Client{}
	resp, err := client.Do(req)
	assert.NoError(err)

	assert.Equal("true", resp.Header.Get("IsAuth"))
}

func TestCORSMiddleware(t *testing.T) {
	assert := assert.New(t)
	ts := httptest.NewServer(CreateMockServer().authMiddleware(CORSTestHandler()))

	req, err := http.NewRequest(http.MethodGet, ts.URL, nil)
	assert.NoError(err)

	client := &http.Client{}
	resp, err := client.Do(req)
	assert.NoError(err)

	assert.Equal(resp.StatusCode, http.StatusOK)
}
