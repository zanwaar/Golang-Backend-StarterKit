package utils_test

import (
	"golang-backend/utils"
	"testing"
)

func TestHashPassword(t *testing.T) {
	password := "secret123"
	hashed, err := utils.HashPassword(password)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	if len(hashed) == 0 {
		t.Fatal("Hashed password is empty")
	}

	if hashed == password {
		t.Fatal("Hashed password should not be equal to original password")
	}
}

func TestCheckPassword(t *testing.T) {
	password := "secret123"
	hashed, _ := utils.HashPassword(password)

	// Test correct password
	if err := utils.CheckPassword(password, hashed); err != nil {
		t.Errorf("Password check failed for correct password: %v", err)
	}

	// Test incorrect password
	if err := utils.CheckPassword("wrongpassword", hashed); err == nil {
		t.Error("Password check should have failed for incorrect password")
	}
}

func TestGenerateRandomCode(t *testing.T) {
	length := 6
	code := utils.GenerateRandomCode(length)

	if len(code) != length {
		t.Errorf("Expected code length %d, got %d", length, len(code))
	}

	// Range check
	for _, char := range code {
		if char < '0' || char > '9' {
			t.Errorf("Code contains non-digit character: %c", char)
		}
	}
}
