package errors

type AppError struct {
	Code    string
	Message string
	Err     error
}

// Gera um novo AppError, sem o campo err
func New(code string, msg string) *AppError {
	{
		return &AppError{
			Code:    code,
			Message: msg,
		}
	}
}

// Associa o erro ao AppError
func Wrap(code string, msg string, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: msg,
		Err:     err,
	}
}

func (ap *AppError) Error() string {
	return ap.Message
}
