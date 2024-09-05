package server

type ApplicationError struct {
    // ErrorCode should have its 3 first digits represent http status code.
    // For example, 50001 will return 500, 40401 will return 404...
    ErrorCode int
    Message   string
}

func NewApplicationError(errorCode int, message string) ApplicationError {
    return ApplicationError{
        ErrorCode: errorCode,
        Message:   message,
    }
}

func (err ApplicationError) Error() string {
    return err.Message
}
