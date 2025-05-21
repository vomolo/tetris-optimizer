package solver

import (
	"errors"
	"testing"
)

func TestValidationErrorImplementsError(t *testing.T) {
	err := newValidationError("invalid input")

	var e *validationError
	if !errors.As(err, &e) {
		t.Errorf("expected error to be of type *validationError, got %T", err)
	}
}

func TestValidationErrorMessage(t *testing.T) {
	msg := "this is a validation error"
	err := newValidationError(msg)

	if err.Error() != msg {
		t.Errorf("expected error message %q, got %q", msg, err.Error())
	}
}
