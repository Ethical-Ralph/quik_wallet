package utils

type AppError struct {
	status int
	err    string
}

var NOT_FOUND_ERROR string = "NOT_FOUND"

func NewError(status int, err string) AppError {
	return AppError{status, err}
}

func GetStatus(err error) int {
	switch err.Error() {
	case NOT_FOUND_ERROR:
		return 404

	default:
		return 400
	}
}
