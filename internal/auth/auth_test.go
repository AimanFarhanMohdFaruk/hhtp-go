package auth_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/AimanFarhanMohdFaruk/hhtp-go.git/internal/auth"
	"github.com/google/uuid"
)

func TestMakeJWT(t *testing.T) {
	userId := uuid.New()
	secret := "mysecretkey"
	testJWT, err := auth.MakeJWT(userId, secret)
	if err != nil {
		t.Fatalf("failed to create JWT: %v", err)
	}
	
	t.Logf("Generated JWT: %s", testJWT)
}

func TestMakeRefreshToken(t *testing.T) {
	generateRefreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		t.Fatalf("failed to create refresh token: %v", err)
	}
	t.Logf("Generated refresh token: %s", generateRefreshToken)
}

func TestGetBearerToken(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Add("Authorization", "Bearer test")

	authHeader, err := auth.GetBearerToken(req.Header)

	if err != nil {
		t.Fatalf("failed to retrieve bearer token")
	}

	want := authHeader == "test"
	if !want {
		t.Fatalf("failed to retrieve bearer token")
	}
}

func TestValidateJWT(t *testing.T) {
	userId := uuid.New()
	secret := "mysecretkey"
	testJWT, err := auth.MakeJWT(userId, secret)
	if err != nil {
		t.Fatalf("failed to create JWT: %v", err)
	}

	parsedId, err := auth.ValidateJWT(testJWT, secret)
	if err != nil {
		t.Fatalf("failed to validate JWT: %v", err)
	}

	want := userId == parsedId
	if !want {
		t.Fatalf("failed to validate correct JWT")
	}
}

// Make sure that you can create and validate JWTs, and that expired tokens are rejected and JWTs signed with the wrong secret are rejected.

