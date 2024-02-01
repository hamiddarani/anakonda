package helper

import "net/http"

type ServiceError struct {
	EndUserMessage   string `json:"endUserMessage"`
	TechnicalMessage string `json:"technicalMessage"`
	Err              error
}

func (s *ServiceError) Error() string {
	return s.EndUserMessage
}

const (

	// Token
	UnExpectedError = "unexpected error"
	ClaimsNotFound  = "claims not found"
	TokenRequired   = "token required"
	TokenExpired    = "token expired"
	TokenInvalid    = "token invalid"

	// User
	PermissionDenied = "Permission denied"

	// DB
	RecordNotFound = "record not found"
)

var StatusCodeMapping = map[string]int{
	// DB
	RecordNotFound: 400,
}

func TranslateErrorToStatusCode(err error) int {
	value, ok := StatusCodeMapping[err.Error()]
	if !ok {
		return http.StatusInternalServerError
	}
	return value
}
