package password

import (
	"testing"
)

func TestHasher_HashPassword(t *testing.T) {
	hasher := DefaultHasher()
	password := "testPassword123!"

	hashedPassword, err := hasher.HashPassword(password)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	if hashedPassword == password {
		t.Error("Hashed password should not be the same as plain password")
	}

	if len(hashedPassword) == 0 {
		t.Error("Hashed password should not be empty")
	}
}

func TestHasher_VerifyPassword(t *testing.T) {
	hasher := DefaultHasher()
	password := "testPassword123!"

	hashedPassword, err := hasher.HashPassword(password)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	// Test correct password
	if !hasher.VerifyPassword(hashedPassword, password) {
		t.Error("Password verification should succeed for correct password")
	}

	// Test incorrect password
	if hasher.VerifyPassword(hashedPassword, "wrongPassword") {
		t.Error("Password verification should fail for incorrect password")
	}

	// Test empty password
	if hasher.VerifyPassword(hashedPassword, "") {
		t.Error("Password verification should fail for empty password")
	}
}

func TestHasher_DifferentCosts(t *testing.T) {
	password := "testPassword123!"

	// Test with different costs
	costs := []int{10, 12, 14}
	for _, cost := range costs {
		hasher := NewHasher(cost)
		hashedPassword, err := hasher.HashPassword(password)
		if err != nil {
			t.Fatalf("Failed to hash password with cost %d: %v", cost, err)
		}

		if !hasher.VerifyPassword(hashedPassword, password) {
			t.Errorf("Password verification should succeed for cost %d", cost)
		}
	}
}
