package jwt_test

import (
	"linkshorter/pkg/jwt"
	"testing"
)

func TestJWTCreate(t *testing.T) {
	const email = "l2@mail.com"
	jwtService := jwt.NewJWT("/2+XnmJGz1j3ehIVI/5P9k1+CghrE3DcS7rnT+qar5w=")

	token, err := jwtService.Create(jwt.JWTData{
		Email: "l2@mail.com",
	})
	if err != nil {
		t.Fatal(err)
	}

	ok, data := jwtService.Parse(token)
	if !ok {
		t.Fatal("Invalide token")
	}

	if data.Email != email {
		t.Fatalf("Email %s mot equal %s", data.Email, email)
	}
}