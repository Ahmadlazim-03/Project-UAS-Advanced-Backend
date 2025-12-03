package utils

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestHashPassword(t *testing.T) {
	password := "testpassword123"
	
	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	if hash == "" {
		t.Error("Hash should not be empty")
	}

	if hash == password {
		t.Error("Hash should not equal plain password")
	}
}

func TestCheckPassword(t *testing.T) {
	password := "testpassword123"
	wrongPassword := "wrongpassword"
	
	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	// Test correct password
	if !CheckPassword(password, hash) {
		t.Error("CheckPassword should return true for correct password")
	}

	// Test wrong password
	if CheckPassword(wrongPassword, hash) {
		t.Error("CheckPassword should return false for wrong password")
	}
}

func TestGenerateToken(t *testing.T) {
	userID := uuid.New()
	username := "testuser"
	email := "test@example.com"
	roleID := uuid.New()
	roleName := "Admin"
	permissions := []string{"user:read", "user:create"}
	secret := "test-secret-key"
	expiry := 1 * time.Hour

	token, err := GenerateToken(userID, username, email, roleID, roleName, permissions, secret, expiry)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	if token == "" {
		t.Error("Token should not be empty")
	}

	// Token should be a JWT (has 3 parts separated by dots)
	parts := len(token)
	if parts < 50 {
		t.Error("Token seems too short to be a valid JWT")
	}
}
