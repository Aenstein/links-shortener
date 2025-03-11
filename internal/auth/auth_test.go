package auth_test

import (
	"bytes"
	"encoding/json"
	"linkshorter/configs"
	"linkshorter/internal/auth"
	"linkshorter/internal/user"
	"linkshorter/pkg/db"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func bootstrap() (*auth.AuthHandler, sqlmock.Sqlmock, error) {
	database, mock, err := sqlmock.New()
	if err != nil {
		return nil, nil, err
	}

	gormDb, err := gorm.Open(postgres.New(postgres.Config{
		Conn: database,
	}))
	if err != nil {
		return nil, nil, err
	}

	userRepository := user.NewUserRepository(&db.Db{
		DB: gormDb,
	})

	handler := auth.AuthHandler{
		Config: &configs.Config{
			Auth: configs.AuthConfig{
				Secret: "secret",
			},
		},
		AuthService: auth.NewAuthService(userRepository),
	}

	return &handler, mock, nil
}

func TestHandlerLoginSuccess(t *testing.T) {
	handler, mock, err := bootstrap()
	if err != nil {
		t.Fatal(err)
		return
	}
	rows := sqlmock.NewRows([]string{"email", "password"}).
		AddRow("l2@mail.com", "$2a$10$LHphcxNMgkvcGfZcHJWNoeJay7bMzlYpdNhwalOXWD5uz9cuZ2MUC")
	mock.ExpectQuery("SELECT").WillReturnRows(rows)

	data, _ := json.Marshal(&auth.LoginRequest{
		Email: "l2@mail.com",
		Password: "l",
	})

	reader := bytes.NewReader(data)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/auth/login", reader)
	
	handler.Login()(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("got %d, expeted %d", w.Code, http.StatusOK)
	}
}

func TestHandlerLoginFailed(t *testing.T) {
	handler, mock, err := bootstrap()
	if err != nil {
		t.Fatal(err)
		return
	}
	rows := sqlmock.NewRows([]string{"email", "password"}).
		AddRow("l2@mail.com", "$2a$10$LHphcxNMgkvcGfZcHJWNoeJay7bMzlYpdNhwalOXWD5uz9cuZ2MUC")
	mock.ExpectQuery("SELECT").WillReturnRows(rows)

	data, _ := json.Marshal(&auth.LoginRequest{
		Email: "l2@mail.com",
		Password: "k",
	})

	reader := bytes.NewReader(data)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/auth/login", reader)
	
	handler.Register()(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Errorf("got %d, expeted %d", w.Code, http.StatusUnauthorized)
	}
}

func TestHandlerRegisterSuccess(t *testing.T) {
	handler, mock, err := bootstrap()
	if err != nil {
		t.Fatal(err)
		return
	}
	rows := sqlmock.NewRows([]string{"email", "password", "name"})
	mock.ExpectQuery("SELECT").WillReturnRows(rows)
	mock.ExpectBegin()
	mock.ExpectQuery("INSERT").WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	data, _ := json.Marshal(&auth.RegisterRequest{
		Email: "l2@mail.com",
		Password: "l",
		Name: "Leonid",
	})

	reader := bytes.NewReader(data)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/auth/register", reader)
	
	handler.Register()(w, req)
	if w.Code != http.StatusCreated {
		t.Errorf("got %d, expeted %d", w.Code, http.StatusCreated)
	}
}