package main

import (
	"bytes"
	"encoding/json"
	"io"
	"linkshorter/internal/auth"
	"linkshorter/internal/user"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB() *gorm.DB {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}
	db, err := gorm.Open(postgres.Open(os.Getenv("DSN")), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	return db
}

func TestLoginFailed(t *testing.T) {
	db := InitDB()
	initData(db)
	defer removeData(db)
	
	ts := httptest.NewServer(App())
	defer ts.Close()

	data, _ := json.Marshal(&auth.LoginRequest{
		Email: "l2@mail.com",
		Password: "a",
	})

	res, err := http.Post(ts.URL + "/auth/login", "application/json", bytes.NewReader(data))
	if err != nil {
		t.Fatal(err)
	}
	if res.StatusCode != 401 {
		t.Fatalf("Expected %d got %d", 401, res.StatusCode)
	}
}

func initData(db *gorm.DB) {
	db.Create(&user.User{
		Email: "l2@mail.com",
		Password: "$2a$10$LHphcxNMgkvcGfZcHJWNoeJay7bMzlYpdNhwalOXWD5uz9cuZ2MUC",
		Name: "Leonid",
	})
}

func removeData(db *gorm.DB) {
	db.Unscoped().
		Where("email = ?", "l2@mail.com").
		Delete(&user.User{ })
}

func TestLoginSuccess(t *testing.T) {
	db := InitDB()
	initData(db)
	defer removeData(db)

	ts := httptest.NewServer(App())
	defer ts.Close()

	data, _ := json.Marshal(&auth.LoginRequest{
		Email: "l2@mail.com",
		Password: "l",
	})

	res, err := http.Post(ts.URL + "/auth/login", "application/json", bytes.NewReader(data))
	if err != nil {
		t.Fatal(err)
	}
	if res.StatusCode != 200 {
		t.Fatalf("Expected %d got %d", 200, res.StatusCode)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}

	var resData auth.LoginResponse

	err = json.Unmarshal(body, &resData)
	if err != nil {
		t.Fatal(err)
	}

	if resData.Token == "" {
		t.Fatal("Token empty")
	}
}