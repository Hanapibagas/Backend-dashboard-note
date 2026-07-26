package errors

import (
	"errors"
	"testing"
)

func TestDomainError_Error(t *testing.T) {
	tests := []struct {
		name    string
		err     *DomainError
		wantMsg string
	}{
		{
			name:    "error without underlying error",
			err:     ErrEmailAlreadyExists,
			wantMsg: "EMAIL_ALREADY_EXISTS: email already exists",
		},
		{
			name:    "error with underlying error",
			err:     NewDomainError("TEST_CODE", "test message", errors.New("underlying")),
			wantMsg: "TEST_CODE: test message (caused by: underlying)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.Error(); got != tt.wantMsg {
				t.Errorf("DomainError.Error() = %v, want %v", got, tt.wantMsg)
			}
		})
	}
}

func TestDomainError_Is(t *testing.T) {
	err1 := ErrEmailAlreadyExists
	err2 := &DomainError{Code: "EMAIL_ALREADY_EXISTS", Message: "email already exists"}
	err3 := ErrUserNotFound

	if !errors.Is(err1, err2) {
		t.Error("Expected errors to match by code")
	}

	if errors.Is(err1, err3) {
		t.Error("Expected errors not to match")
	}
}

func TestDomainError_Unwrap(t *testing.T) {
	underlying := errors.New("underlying error")
	err := NewDomainError("TEST", "test", underlying)

	if unwrapped := errors.Unwrap(err); unwrapped != underlying {
		t.Errorf("Unwrap() = %v, want %v", unwrapped, underlying)
	}
}

func TestWrapError(t *testing.T) {
	underlying := errors.New("underlying error")
	wrapped := WrapError(ErrFailedToCheckEmail, underlying)

	if wrapped == nil {
		t.Fatal("WrapError() should not return nil when err is not nil")
	}

	if !errors.Is(wrapped, ErrFailedToCheckEmail) {
		t.Error("Wrapped error should match the domain error")
	}

	if errors.Unwrap(wrapped) != underlying {
		t.Error("Wrapped error should have the underlying error")
	}
}

func TestPredefinedErrors(t *testing.T) {
	tests := []struct {
		name string
		err  *DomainError
		code string
		msg  string
	}{
		{"ErrEmailAlreadyExists", ErrEmailAlreadyExists, "EMAIL_ALREADY_EXISTS", "email already exists"},
		{"ErrUserNotFound", ErrUserNotFound, "USER_NOT_FOUND", "user not found"},
		{"ErrInvalidCredentials", ErrInvalidCredentials, "INVALID_CREDENTIALS", "invalid credentials"},
		{"ErrInvalidEmailFormat", ErrInvalidEmailFormat, "INVALID_EMAIL_FORMAT", "invalid email format"},
		{"ErrPasswordTooShort", ErrPasswordTooShort, "PASSWORD_TOO_SHORT", "password must be at least 8 characters"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.err.Code != tt.code {
				t.Errorf("Error code = %v, want %v", tt.err.Code, tt.code)
			}
			if tt.err.Message != tt.msg {
				t.Errorf("Error message = %v, want %v", tt.err.Message, tt.msg)
			}
		})
	}
}
