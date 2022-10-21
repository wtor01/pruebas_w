package apperrors

type AppError struct {
	err         string
	shouldRetry bool
}

func (o AppError) Error() string {
	return o.err
}

func (o AppError) ShouldRetry() bool {
	return o.shouldRetry
}

func NewAppError(str string, shouldRetry bool) AppError {
	return AppError{err: str, shouldRetry: shouldRetry}
}
