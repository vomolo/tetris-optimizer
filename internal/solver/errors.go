package solver

type validationError struct {
    message string
}

func (e *validationError) Error() string {
    return e.message
}

func newValidationError(message string) error {
    return &validationError{message: message}
}