package auth

import (
	"errors"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func TestMakeJWT_ValidateJWT_RoundTrip(t *testing.T) {
	userID := uuid.New()
	secret := "my-super-secret-key-for-testing-12345"
	expiresIn := 1 * time.Hour

	token, err := MakeJWT(userID, secret, expiresIn)
	if err != nil {
		t.Fatalf("MakeJWT failed: %v", err)
	}

	validatedID, err := ValidateJWT(token, secret)
	if err != nil {
		t.Fatalf("ValidateJWT failed: %v", err)
	}

	if validatedID != userID {
		t.Errorf("UUID round-trip failed\n   got:  %s\n   want: %s", validatedID, userID)
	}
}
func TestValidateJWT_RejectsExpiredToken(t *testing.T) {
	userID := uuid.New()
	secret := "my-super-secret-key-for-testing-12345"
	expiresIn := -1 * time.Hour

	token, err := MakeJWT(userID, secret, expiresIn)
	if err != nil {
		t.Fatalf("MakeJWT failed: %v", err)
	}

	_, err = ValidateJWT(token, secret)

	if err == nil {
		t.Fatal("expected ValidateJWT to reject an expired token, but it returned nil error")
	}

	if !errors.Is(err, jwt.ErrTokenExpired) {
		t.Errorf("expected jwt.ErrTokenExpired, got: %v", err)
	}
}
func TestValidateJWT_RejectsTokenSignedWithWrongSecret(t *testing.T) {
	userID := uuid.New()
	correctSecret := "my-super-secret-key-for-testing-12345"
	wrongSecret := "this-is-definitely-the-wrong-secret-99999"
	expiresIn := 1 * time.Hour

	token, err := MakeJWT(userID, correctSecret, expiresIn)
	if err != nil {
		t.Fatalf("MakeJWT failed: %v", err)
	}

	_, err = ValidateJWT(token, wrongSecret)

	if err == nil {
		t.Fatal("expected ValidateJWT to reject token signed with wrong secret, but got nil error")
	}

	if !errors.Is(err, jwt.ErrTokenSignatureInvalid) {
		t.Errorf("expected jwt.ErrTokenSignatureInvalid, got: %v", err)
	}
}
