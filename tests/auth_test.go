package tests

import (
	"testing"

	"github.com/Ahmadlazim-03/Project-UAS-Advanced-Backend/utils"
)

func TestHashPassword(t *testing.T) {
	password := "secret"
	hash, err := utils.HashPassword(password)
	if err != nil {
		t.Errorf("HashPassword returned error: %v", err)
	}

	if hash == "" {
		t.Error("HashPassword returned empty string")
	}

	if !utils.CheckPasswordHash(password, hash) {
		t.Error("CheckPasswordHash failed to verify correct password")
	}

	if utils.CheckPasswordHash("wrong", hash) {
		t.Error("CheckPasswordHash verified wrong password")
	}
}
