package utils

import (
	"testing"
)

func TestValidateStruct_ValidData(t *testing.T) {
	type TestStruct struct {
		Email    string `validate:"required,email"`
		Username string `validate:"required,min=3,max=50"`
		Age      int    `validate:"required,min=1,max=150"`
	}

	valid := TestStruct{
		Email:    "test@example.com",
		Username: "testuser",
		Age:      25,
	}

	err := ValidateStruct(&valid)
	if err != nil {
		t.Errorf("Expected no error for valid data, got: %v", err)
	}
}

func TestValidateStruct_InvalidEmail(t *testing.T) {
	type TestStruct struct {
		Email string `validate:"required,email"`
	}

	invalid := TestStruct{
		Email: "not-an-email",
	}

	err := ValidateStruct(&invalid)
	if err == nil {
		t.Error("Expected error for invalid email, got nil")
	}
}

func TestValidateStruct_RequiredField(t *testing.T) {
	type TestStruct struct {
		Name string `validate:"required"`
	}

	invalid := TestStruct{
		Name: "",
	}

	err := ValidateStruct(&invalid)
	if err == nil {
		t.Error("Expected error for empty required field, got nil")
	}
}

func TestValidateStruct_MinLength(t *testing.T) {
	type TestStruct struct {
		Password string `validate:"required,min=6"`
	}

	invalid := TestStruct{
		Password: "123",
	}

	err := ValidateStruct(&invalid)
	if err == nil {
		t.Error("Expected error for password too short, got nil")
	}

	valid := TestStruct{
		Password: "123456",
	}

	err = ValidateStruct(&valid)
	if err != nil {
		t.Errorf("Expected no error for valid password length, got: %v", err)
	}
}

func TestValidateStruct_MaxLength(t *testing.T) {
	type TestStruct struct {
		Username string `validate:"required,max=10"`
	}

	invalid := TestStruct{
		Username: "verylongusername",
	}

	err := ValidateStruct(&invalid)
	if err == nil {
		t.Error("Expected error for username too long, got nil")
	}
}
