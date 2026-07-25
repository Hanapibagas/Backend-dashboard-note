package valueobject_test

import (
	"errors"
	"testing"

	errorDomain "auth/domain/errors"
	"auth/domain/valueobject"
)

// Test NewEmail - Positive Cases
func TestNewEmail_Success(t *testing.T) {
	tests := []struct {
		name  string
		email string
	}{
		{
			name:  "valid simple email",
			email: "user@example.com",
		},
		{
			name:  "valid email with subdomain",
			email: "user@mail.example.com",
		},
		{
			name:  "valid email with numbers",
			email: "user123@example.com",
		},
		{
			name:  "valid email with dots",
			email: "user.name@example.com",
		},
		{
			name:  "valid email with plus sign",
			email: "user+tag@example.com",
		},
		{
			name:  "valid email with hyphen",
			email: "user-name@example.com",
		},
		{
			name:  "valid email with underscore",
			email: "user_name@example.com",
		},
		{
			name:  "valid email with percent",
			email: "user%name@example.com",
		},
		{
			name:  "valid email with uppercase (should be converted to lowercase)",
			email: "USER@EXAMPLE.COM",
		},
		{
			name:  "valid email with mixed case",
			email: "UsEr@ExAmPlE.com",
		},
		{
			name:  "valid email with leading/trailing whitespace",
			email: "  user@example.com  ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			email, err := valueobject.NewEmail(tt.email)

			if err != nil {
				t.Errorf("NewEmail() error = %v, want nil", err)
				return
			}

			if email == nil {
				t.Error("NewEmail() returned nil, want valid email")
				return
			}

			// Test String() method
			emailStr := email.String()
			if emailStr == "" {
				t.Error("String() returned empty string")
			}

			// Verify uppercase letters are converted to lowercase
			expectedLowercase := "user@example.com"
			if tt.email == "USER@EXAMPLE.COM" || tt.email == "UsEr@ExAmPlE.com" {
				if emailStr != expectedLowercase {
					t.Errorf("String() = %v, want %v (lowercase)", emailStr, expectedLowercase)
				}
			}

			// Verify whitespace is trimmed
			if tt.email == "  user@example.com  " {
				if emailStr != expectedLowercase {
					t.Errorf("String() = %v, want %v (trimmed)", emailStr, expectedLowercase)
				}
			}
		})
	}
}

// Test NewEmail - Negative Cases
func TestNewEmail_ValidationErrors(t *testing.T) {
	tests := []struct {
		name        string
		email       string
		expectedErr error
	}{
		{
			name:        "empty email",
			email:       "",
			expectedErr: errorDomain.ErrEmailEmpty,
		},
		{
			name:        "whitespace only email",
			email:       "   ",
			expectedErr: errorDomain.ErrEmailEmpty,
		},
		{
			name:        "missing @ sign",
			email:       "userexample.com",
			expectedErr: errorDomain.ErrInvalidEmailFormat,
		},
		{
			name:        "missing domain",
			email:       "user@",
			expectedErr: errorDomain.ErrInvalidEmailFormat,
		},
		{
			name:        "missing local part",
			email:       "@example.com",
			expectedErr: errorDomain.ErrInvalidEmailFormat,
		},
		{
			name:        "missing TLD",
			email:       "user@example",
			expectedErr: errorDomain.ErrInvalidEmailFormat,
		},
		{
			name:        "invalid characters in local part",
			email:       "user name@example.com",
			expectedErr: errorDomain.ErrInvalidEmailFormat,
		},
		{
			name:        "invalid TLD - single character",
			email:       "user@example.c",
			expectedErr: errorDomain.ErrInvalidEmailFormat,
		},
		{
			name:        "multiple @ signs",
			email:       "user@@example.com",
			expectedErr: errorDomain.ErrInvalidEmailFormat,
		},
		{
			name:        "no TLD at all",
			email:       "user@.",
			expectedErr: errorDomain.ErrInvalidEmailFormat,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			email, err := valueobject.NewEmail(tt.email)

			if err == nil {
				t.Error("NewEmail() error = nil, want error")
				return
			}

			if err != tt.expectedErr {
				t.Errorf("NewEmail() error = %v, want %v", err, tt.expectedErr)
			}

			if email != nil {
				t.Error("NewEmail() returned non-nil, want nil on error")
			}
		})
	}
}

// Test NewPassword - Positive Cases
func TestNewPassword_Success(t *testing.T) {
	tests := []struct {
		name     string
		password string
	}{
		{
			name:     "valid password with all requirements",
			password: "Password123",
		},
		{
			name:     "valid password exactly 8 characters",
			password: "Pass1234",
		},
		{
			name:     "valid password longer than 8 characters",
			password: "MySecurePassword123",
		},
		{
			name:     "valid password with special characters",
			password: "P@ssw0rd123!",
		},
		{
			name:     "valid password with numbers only as digits",
			password: "Password1",
		},
		{
			name:     "valid password with mixed letters and numbers",
			password: "Abcd1234EFGH",
		},
		{
			name:     "valid password starting with uppercase",
			password: "Abcdefg1",
		},
		{
			name:     "valid password starting with lowercase",
			password: "aBCDEFG1",
		},
		{
			name:     "valid password starting with number",
			password: "1AbcdefG",
		},
		{
			name:     "valid password with multiple uppercase",
			password: "PASSword123",
		},
		{
			name:     "valid password with multiple lowercase",
			password: "PASSwordWORD123",
		},
		{
			name:     "valid password with multiple numbers",
			password: "Password123456",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			password, err := valueobject.NewPassword(tt.password)

			if err != nil {
				t.Errorf("NewPassword() error = %v, want nil", err)
				return
			}

			if password == nil {
				t.Error("NewPassword() returned nil, want valid password")
			}
		})
	}
}

// Test NewPassword - Negative Cases
func TestNewPassword_ValidationerrorDomain(t *testing.T) {
	tests := []struct {
		name        string
		password    string
		expectedErr error
	}{
		{
			name:        "empty password",
			password:    "",
			expectedErr: errorDomain.ErrPasswordEmpty,
		},
		{
			name:        "password less than 8 characters - 7 chars",
			password:    "Pass123",
			expectedErr: errorDomain.ErrPasswordTooShort,
		},
		{
			name:        "password less than 8 characters - 1 char",
			password:    "A",
			expectedErr: errorDomain.ErrPasswordTooShort,
		},
		{
			name:        "password exactly 8 chars but no uppercase",
			password:    "password1",
			expectedErr: errorDomain.ErrPasswordMissingUppercase,
		},
		{
			name:        "password with lowercase and numbers only",
			password:    "password123",
			expectedErr: errorDomain.ErrPasswordMissingUppercase,
		},
		{
			name:        "password with numbers only",
			password:    "12345678",
			expectedErr: errorDomain.ErrPasswordMissingUppercase,
		},
		{
			name:        "password exactly 8 chars but no lowercase",
			password:    "PASSWORD1",
			expectedErr: errorDomain.ErrPasswordMissingLowercase,
		},
		{
			name:        "password with uppercase and numbers only",
			password:    "PASSWORD123",
			expectedErr: errorDomain.ErrPasswordMissingLowercase,
		},
		{
			name:        "password with letters only (mixed case)",
			password:    "Password",
			expectedErr: errorDomain.ErrPasswordMissingNumber,
		},
		{
			name:        "password with uppercase only (8+ chars)",
			password:    "PASSWORD",
			expectedErr: errorDomain.ErrPasswordMissingLowercase,
		},
		{
			name:        "password with lowercase only (8+ chars)",
			password:    "password",
			expectedErr: errorDomain.ErrPasswordMissingUppercase,
		},
		{
			name:        "password with uppercase and lowercase but no number",
			password:    "Password",
			expectedErr: errorDomain.ErrPasswordMissingNumber,
		},
		{
			name:        "password with special chars only",
			password:    "!@#$%^&*",
			expectedErr: errorDomain.ErrPasswordMissingUppercase,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			password, err := valueobject.NewPassword(tt.password)

			if err == nil {
				t.Error("NewPassword() error = nil, want error")
				return
			}

			if err != tt.expectedErr {
				t.Errorf("NewPassword() error = %v, want %v", err, tt.expectedErr)
			}

			if password != nil {
				t.Error("NewPassword() returned non-nil, want nil on error")
			}
		})
	}
}

// Test Email String Method
func TestEmail_String(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "simple email",
			input:    "user@example.com",
			expected: "user@example.com",
		},
		{
			name:     "email with uppercase converted to lowercase",
			input:    "USER@EXAMPLE.COM",
			expected: "user@example.com",
		},
		{
			name:     "email with mixed case converted to lowercase",
			input:    "UsEr@ExAmPlE.com",
			expected: "user@example.com",
		},
		{
			name:     "email with whitespace trimmed",
			input:    "  user@example.com  ",
			expected: "user@example.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			email, err := valueobject.NewEmail(tt.input)
			if err != nil {
				t.Fatalf("NewEmail() error = %v", err)
			}

			result := email.String()
			if result != tt.expected {
				t.Errorf("String() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// Test Password Hash Method
func TestPassword_Hash(t *testing.T) {
	tests := []struct {
		name     string
		password string
	}{
		{
			name:     "hash valid password",
			password: "Password123",
		},
		{
			name:     "hash password with special chars",
			password: "P@ssw0rd123!",
		},
		{
			name:     "hash password with numbers",
			password: "Abc123456",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			password, err := valueobject.NewPassword(tt.password)
			if err != nil {
				t.Fatalf("NewPassword() error = %v", err)
			}

			hashedPassword, err := password.Hash()
			if err != nil {
				t.Errorf("Hash() error = %v, want nil", err)
				return
			}

			if hashedPassword == nil {
				t.Error("Hash() returned nil, want valid hashed password")
				return
			}

			// Verify hash is not empty
			hashStr := hashedPassword.String()
			if hashStr == "" {
				t.Error("HashedPassword.String() returned empty string")
			}

			// Verify hash is not equal to original password
			if hashStr == tt.password {
				t.Error("Hash should not be equal to plain password")
			}

			// Verify hash is in bcrypt format (starts with $2a$ or $2b$)
			if len(hashStr) < 4 || hashStr[0:3] != "$2a" && hashStr[0:3] != "$2b" {
				t.Errorf("Hash should be in bcrypt format, got: %s", hashStr[:4])
			}
		})
	}
}

// Test HashedPassword VerifyPassword Method
func TestHashedPassword_VerifyPassword(t *testing.T) {
	tests := []struct {
		name          string
		password      string
		wrongPassword string
	}{
		{
			name:          "verify correct password",
			password:      "Password123",
			wrongPassword: "WrongPassword123",
		},
		{
			name:          "verify password with special chars",
			password:      "P@ssw0rd123!",
			wrongPassword: "Wrong@123!",
		},
		{
			name:          "verify password with numbers",
			password:      "Abc123456",
			wrongPassword: "Wrong123456",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create password and hash it
			password, err := valueobject.NewPassword(tt.password)
			if err != nil {
				t.Fatalf("NewPassword() error = %v", err)
			}

			hashedPassword, err := password.Hash()
			if err != nil {
				t.Fatalf("Hash() error = %v", err)
			}

			// Verify correct password should succeed
			err = hashedPassword.VerifyPassword(tt.password)
			if err != nil {
				t.Errorf("VerifyPassword() with correct password returned error: %v", err)
			}

			// Verify wrong password should fail
			err = hashedPassword.VerifyPassword(tt.wrongPassword)
			if err == nil {
				t.Error("VerifyPassword() with wrong password should return error")
			}
			if !errors.Is(err, errorDomain.ErrPasswordMismatch) {
				t.Errorf("VerifyPassword() should return ErrPasswordMismatch, got: %v", err)
			}
		})
	}
}

// Test NewHashedPassword
func TestNewHashedPassword(t *testing.T) {
	tests := []struct {
		name string
		hash string
	}{
		{
			name: "valid bcrypt hash",
			hash: "$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy",
		},
		{
			name: "empty hash",
			hash: "",
		},
		{
			name: "invalid hash format",
			hash: "invalid-hash",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hashedPassword := valueobject.NewHashedPassword(tt.hash)

			if hashedPassword == nil {
				t.Error("NewHashedPassword() returned nil")
				return
			}

			result := hashedPassword.String()
			if result != tt.hash {
				t.Errorf("String() = %v, want %v", result, tt.hash)
			}
		})
	}
}

// Test HashedPassword Integration Test
func TestHashedPassword_Integration(t *testing.T) {
	// This test verifies the complete flow: Password → Hash → Verify
	plainPassword := "TestPassword123"

	// Step 1: Create Password VO
	password, err := valueobject.NewPassword(plainPassword)
	if err != nil {
		t.Fatalf("NewPassword() error = %v", err)
	}

	// Step 2: Hash the password
	hashedPassword, err := password.Hash()
	if err != nil {
		t.Fatalf("Hash() error = %v", err)
	}

	// Step 3: Verify correct password
	err = hashedPassword.VerifyPassword(plainPassword)
	if err != nil {
		t.Errorf("VerifyPassword() with correct password failed: %v", err)
	}

	// Step 4: Verify wrong password fails
	err = hashedPassword.VerifyPassword("WrongPassword123")
	if err == nil {
		t.Error("VerifyPassword() with wrong password should fail")
	}

	// Step 5: Create HashedPassword from hash string (simulating loading from DB)
	hashStr := hashedPassword.String()
	loadedHash := valueobject.NewHashedPassword(hashStr)

	// Step 6: Verify password with loaded hash
	err = loadedHash.VerifyPassword(plainPassword)
	if err != nil {
		t.Errorf("VerifyPassword() with loaded hash failed: %v", err)
	}
}
