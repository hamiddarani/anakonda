package helper

type ResultCode int

const (
	Success             ResultCode = 0
	BadRequestError     ResultCode = 40000
	ValidationErrorCode ResultCode = 40001
	NotFoundError       ResultCode = 40401
	CustomRecovery      ResultCode = 50001
	InternalError       ResultCode = 50002
	DBConnectionError   ResultCode = 50003
)
