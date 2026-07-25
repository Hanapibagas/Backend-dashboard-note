package entity_test

import (
	"auth/domain/entity"
	"auth/domain/valueobject"
	"testing"
	"time"

	"github.com/google/uuid"
)

// Helper function to create a test user with value objects
func createTestUser(email, password, fullName string) (*entity.User, error) {
	emailVO, err := valueobject.NewEmail(email)
	if err != nil {
		return nil, err
	}

	passwordVO, err := valueobject.NewPassword(password)
	if err != nil {
		return nil, err
	}

	hashedPassword, err := passwordVO.Hash()
	if err != nil {
		return nil, err
	}

	return entity.NewUser(emailVO, hashedPassword, fullName)
}

// Test NewUser - Positive Cases
func TestNewUser_Success(t *testing.T) {
	tests := []struct {
		name     string
		email    string
		password string
		fullName string
	}{
		{
			name:     "valid user with all fields",
			email:    "user@example.com",
			password: "Password123",
			fullName: "John Doe",
		},
		{
			name:     "valid user with minimum password length",
			email:    "test@example.com",
			password: "Pass1234",
			fullName: "Jane Smith",
		},
		{
			name:     "valid user with long password",
			email:    "admin@example.com",
			password: "MyVerySecurePassword123!",
			fullName: "Admin User",
		},
		{
			name:     "valid user with special characters in name",
			email:    "user2@example.com",
			password: "Password123",
			fullName: "José María García-López",
		},
		{
			name:     "valid user with numbers in name",
			email:    "user3@example.com",
			password: "Test1234",
			fullName: "User123",
		},
		{
			name:     "valid user with single character name",
			email:    "user4@example.com",
			password: "Pass1234",
			fullName: "A",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create user with value objects
			user, err := createTestUser(tt.email, tt.password, tt.fullName)

			if err != nil {
				t.Errorf("NewUser() error = %v, want nil", err)
				return
			}

			if user == nil {
				t.Error("NewUser() returned nil, want valid user")
				return
			}

			// Verify ID is not empty
			if user.ID == uuid.Nil {
				t.Error("NewUser() ID is empty UUID")
			}

			// Verify email is set correctly
			if user.GetEmail() != tt.email {
				t.Errorf("NewUser() Email = %v, want %v", user.GetEmail(), tt.email)
			}

			// Verify password hash is set and not empty
			if user.GetPasswordHash() == "" {
				t.Error("NewUser() PasswordHash is empty")
			}

			// Verify password hash is not the same as plain password
			if user.GetPasswordHash() == tt.password {
				t.Error("NewUser() PasswordHash should not be plain password")
			}

			// Verify full name is set correctly
			if user.FullName != tt.fullName {
				t.Errorf("NewUser() FullName = %v, want %v", user.FullName, tt.fullName)
			}

			// Verify CreatedAt is set and is recent (within last minute)
			if user.CreatedAt.IsZero() {
				t.Error("NewUser() CreatedAt is zero")
			}

			if time.Since(user.CreatedAt) > time.Minute {
				t.Error("NewUser() CreatedAt is too far in the past")
			}

			// Verify UpdatedAt is set and is recent
			if user.UpdatedAt.IsZero() {
				t.Error("NewUser() UpdatedAt is zero")
			}

			if time.Since(user.UpdatedAt) > time.Minute {
				t.Error("NewUser() UpdatedAt is too far in the past")
			}

			// Verify CreatedAt and UpdatedAt are approximately equal
			if user.CreatedAt.Sub(user.UpdatedAt) > time.Second {
				t.Error("NewUser() CreatedAt and UpdatedAt should be approximately equal")
			}
		})
	}
}

// Test GetID Method
func TestUser_GetID(t *testing.T) {
	emailVO, _ := valueobject.NewEmail("test@example.com")
	passwordVO, _ := valueobject.NewPassword("Password123")
	hashedPassword, _ := passwordVO.Hash()

	user, err := entity.NewUser(emailVO, hashedPassword, "Test User")
	if err != nil {
		t.Fatalf("NewUser() error = %v", err)
	}

	id := user.GetID()
	if id == uuid.Nil {
		t.Error("GetID() returned empty UUID")
	}

	if id != user.ID {
		t.Errorf("GetID() = %v, want %v", id, user.ID)
	}
}

// Test GetEmail Method
func TestUser_GetEmail(t *testing.T) {
	tests := []struct {
		name  string
		email string
	}{
		{
			name:  "simple email",
			email: "user@example.com",
		},
		{
			name:  "email with numbers",
			email: "user123@example.com",
		},
		{
			name:  "email with dots",
			email: "first.last@example.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := createTestUser(tt.email, "Password123", "Test User")
			if err != nil {
				t.Fatalf("NewUser() error = %v", err)
			}

			email := user.GetEmail()
			if email != tt.email {
				t.Errorf("GetEmail() = %v, want %v", email, tt.email)
			}
		})
	}
}

// Test GetFullName Method
func TestUser_GetFullName(t *testing.T) {
	tests := []struct {
		name     string
		fullName string
	}{
		{
			name:     "simple name",
			fullName: "John Doe",
		},
		{
			name:     "name with special characters",
			fullName: "José María",
		},
		{
			name:     "single character name",
			fullName: "A",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := createTestUser("test@example.com", "Password123", tt.fullName)
			if err != nil {
				t.Fatalf("NewUser() error = %v", err)
			}

			fullName := user.GetFullName()
			if fullName != tt.fullName {
				t.Errorf("GetFullName() = %v, want %v", fullName, tt.fullName)
			}
		})
	}
}

// Test GetCreatedAt Method
func TestUser_GetCreatedAt(t *testing.T) {
	user, err := createTestUser("test@example.com", "Password123", "Test User")
	if err != nil {
		t.Fatalf("NewUser() error = %v", err)
	}

	createdAt := user.GetCreatedAt()
	if createdAt.IsZero() {
		t.Error("GetCreatedAt() returned zero time")
	}

	if createdAt != user.CreatedAt {
		t.Errorf("GetCreatedAt() = %v, want %v", createdAt, user.CreatedAt)
	}

	// Verify it's approximately now
	if time.Since(createdAt) > time.Minute {
		t.Error("GetCreatedAt() time is too far in the past")
	}
}

// Test UpdateFullName Method
func TestUser_UpdateFullName(t *testing.T) {
	tests := []struct {
		name        string
		initialName string
		newName     string
	}{
		{
			name:        "update to different name",
			initialName: "John Doe",
			newName:     "Jane Smith",
		},
		{
			name:        "update to name with special characters",
			initialName: "Test User",
			newName:     "José María García-López",
		},
		{
			name:        "update to single character",
			initialName: "Long Name Here",
			newName:     "A",
		},
		{
			name:        "update to name with numbers",
			initialName: "Old Name",
			newName:     "User123",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := createTestUser("test@example.com", "Password123", tt.initialName)
			if err != nil {
				t.Fatalf("NewUser() error = %v", err)
			}

			initialUpdatedAt := user.UpdatedAt
			// Wait a tiny bit to ensure timestamp difference
			time.Sleep(time.Millisecond)

			user.UpdateFullName(tt.newName)

			if user.FullName != tt.newName {
				t.Errorf("UpdateFullName() FullName = %v, want %v", user.FullName, tt.newName)
			}

			// Verify UpdatedAt was updated
			if !user.UpdatedAt.After(initialUpdatedAt) {
				t.Error("UpdateFullName() UpdatedAt should be after initial UpdatedAt")
			}
		})
	}
}

// Test UpdatePassword Method - Positive Cases
func TestUser_UpdatePassword_Success(t *testing.T) {
	tests := []struct {
		name         string
		initialPass  string
		newPassword  string
	}{
		{
			name:         "update to valid password",
			initialPass:  "OldPassword123",
			newPassword:  "NewPassword123",
		},
		{
			name:         "update to password with special chars",
			initialPass:  "Password123",
			newPassword:  "N3wP@ssw0rd!",
		},
		{
			name:         "update to longer password",
			initialPass:  "Pass1234",
			newPassword:  "ThisIsAVeryLongPassword123!",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := createTestUser("test@example.com", tt.initialPass, "Test User")
			if err != nil {
				t.Fatalf("NewUser() error = %v", err)
			}

			// Create new hashed password
			newPasswordVO, err := valueobject.NewPassword(tt.newPassword)
			if err != nil {
				t.Fatalf("NewPassword() error = %v", err)
			}
			newHashedPassword, err := newPasswordVO.Hash()
			if err != nil {
				t.Fatalf("Hash() error = %v", err)
			}

			oldPasswordHash := user.GetPasswordHash()
			initialUpdatedAt := user.UpdatedAt
			// Wait a tiny bit to ensure timestamp difference
			time.Sleep(time.Millisecond)

			user.UpdatePassword(newHashedPassword)

			// Verify password hash changed
			if user.GetPasswordHash() == oldPasswordHash {
				t.Error("UpdatePassword() PasswordHash should have changed")
			}

			// Verify password hash is not the same as plain password
			if user.GetPasswordHash() == tt.newPassword {
				t.Error("UpdatePassword() PasswordHash should not be plain password")
			}

			// Verify UpdatedAt was updated
			if !user.UpdatedAt.After(initialUpdatedAt) {
				t.Error("UpdatePassword() UpdatedAt should be after initial UpdatedAt")
			}
		})
	}
}

// Test VerifyPassword Method - Positive Cases
func TestUser_VerifyPassword_Success(t *testing.T) {
	tests := []struct {
		name     string
		password string
	}{
		{
			name:     "correct password",
			password: "Password123",
		},
		{
			name:     "password with special characters",
			password: "P@ssw0rd123!",
		},
		{
			name:     "password exactly 8 characters",
			password: "Pass1234",
		},
		{
			name:     "long password",
			password: "ThisIsAVeryLongPassword123!",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := createTestUser("test@example.com", tt.password, "Test User")
			if err != nil {
				t.Fatalf("NewUser() error = %v", err)
			}

			err = user.VerifyPassword(tt.password)
			if err != nil {
				t.Errorf("VerifyPassword() error = %v, want nil", err)
			}
		})
	}
}

// Test VerifyPassword Method - Negative Cases
func TestUser_VerifyPassword_Failure(t *testing.T) {
	user, err := createTestUser("test@example.com", "CorrectPassword123", "Test User")
	if err != nil {
		t.Fatalf("NewUser() error = %v", err)
	}

	tests := []struct {
		name     string
		password string
	}{
		{
			name:     "wrong password",
			password: "WrongPassword123",
		},
		{
			name:     "password with different case",
			password: "correctpassword123",
		},
		{
			name:     "password with one character difference",
			password: "CorrectPassword124",
		},
		{
			name:     "empty password",
			password: "",
		},
		{
			name:     "completely different password",
			password: "RandomPassword456",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err = user.VerifyPassword(tt.password)
			if err == nil {
				t.Error("VerifyPassword() error = nil, want error for incorrect password")
			}
		})
	}
}

// Test Password Hashing Uniqueness
func TestUser_PasswordHashingUniqueness(t *testing.T) {
	password := "Password123"

	// Create two users with the same password
	user1, err := createTestUser("user1@example.com", password, "User One")
	if err != nil {
		t.Fatalf("NewUser() error = %v", err)
	}

	user2, err := createTestUser("user2@example.com", password, "User Two")
	if err != nil {
		t.Fatalf("NewUser() error = %v", err)
	}

	// Password hashes should be different due to bcrypt salt
	if user1.GetPasswordHash() == user2.GetPasswordHash() {
		t.Error("Password hashes should be unique for the same password (due to salt)")
	}

	// But both should verify successfully
	err = user1.VerifyPassword(password)
	if err != nil {
		t.Errorf("User1 VerifyPassword() error = %v", err)
	}

	err = user2.VerifyPassword(password)
	if err != nil {
		t.Errorf("User2 VerifyPassword() error = %v", err)
	}
}

// Test User Fields Immutability after Creation
func TestUser_FieldsImmutability(t *testing.T) {
	email := "test@example.com"
	fullName := "Test User"

	user, err := createTestUser(email, "Password123", fullName)
	if err != nil {
		t.Fatalf("NewUser() error = %v", err)
	}

	// Verify email cannot be changed directly (it's exported though, so this is more of a documentation test)
	if user.GetEmail() != email {
		t.Errorf("Email field changed unexpectedly")
	}

	// Verify ID is set and doesn't change
	initialID := user.ID
	if user.GetID() != initialID {
		t.Error("ID should not change after creation")
	}

	// Verify CreatedAt doesn't change after creation
	initialCreatedAt := user.CreatedAt
	user.UpdateFullName("New Name")
	if user.CreatedAt != initialCreatedAt {
		t.Error("CreatedAt should not change after updates")
	}
}

// Test Multiple Updates
func TestUser_MultipleUpdates(t *testing.T) {
	user, err := createTestUser("test@example.com", "Password123", "Initial Name")
	if err != nil {
		t.Fatalf("NewUser() error = %v", err)
	}

	// Multiple full name updates
	user.UpdateFullName("Name 1")
	user.UpdateFullName("Name 2")
	user.UpdateFullName("Name 3")

	if user.FullName != "Name 3" {
		t.Errorf("After multiple updates, FullName = %v, want Name 3", user.FullName)
	}

	// Multiple password updates
	newPassword1, _ := valueobject.NewPassword("NewPassword1")
	newHashed1, _ := newPassword1.Hash()
	user.UpdatePassword(newHashed1)

	newPassword2, _ := valueobject.NewPassword("NewPassword2")
	newHashed2, _ := newPassword2.Hash()
	user.UpdatePassword(newHashed2)

	// Verify last password works
	err = user.VerifyPassword("NewPassword2")
	if err != nil {
		t.Errorf("VerifyPassword() after updates error = %v", err)
	}

	// Verify old passwords don't work
	err = user.VerifyPassword("NewPassword1")
	if err == nil {
		t.Error("VerifyPassword() with old password should fail")
	}
}
