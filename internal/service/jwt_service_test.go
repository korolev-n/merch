package service_test

import (
	"os"
	"testing"

	"github.com/korolev-n/merch/internal/service"
)

func TestJWT_GenerateAndParse(t *testing.T) {
	_ = os.Setenv("JWT_SECRET", "testsecret")
	jwtSvc := service.NewJWTService()

	token, err := jwtSvc.GenerateToken(1, "user1")
	if err != nil {
		t.Fatalf("GenerateToken failed: %v", err)
	}

	claims, err := jwtSvc.ParseToken(token)
	if err != nil {
		t.Fatalf("ParseToken failed: %v", err)
	}

	if claims.UserID != 1 || claims.Username != "user1" {
		t.Fatalf("claims mismatch: got %v", claims)
	}
}

func TestJWT_ParseInvalid(t *testing.T) {
	_ = os.Setenv("JWT_SECRET", "testsecret")
	jwtSvc := service.NewJWTService()

	_, err := jwtSvc.ParseToken("bad.token.value")
	if err == nil {
		t.Fatal("expected error on invalid token")
	}
}
