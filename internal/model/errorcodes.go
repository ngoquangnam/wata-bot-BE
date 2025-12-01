package model

// Error codes for API responses
const (
	// Validation errors (0001-0099)
	ErrCodeInvalidAddressFormat = "0001"
	ErrCodeInvalidSignature     = "0002"
	ErrCodeInvalidMessage       = "0003"

	// Authentication errors (0100-0199)
	ErrCodeTokenGenerationFailed = "0100"

	// Database errors (0200-0299)
	ErrCodeDatabaseError      = "0200"
	ErrCodeFailedToCreateUser = "0201"
	ErrCodeFailedToFindUser   = "0202"

	// Server errors (0500-0599)
	ErrCodeInternalServerError = "0500"
)

// Error messages
const (
	ErrMsgInvalidAddressFormat  = "invalid address format"
	ErrMsgInvalidSignature      = "invalid signature"
	ErrMsgTokenGenerationFailed = "failed to generate tokens"
	ErrMsgDatabaseError         = "database error"
	ErrMsgFailedToCreateUser    = "failed to create user"
	ErrMsgFailedToFindUser      = "failed to find user"
	ErrMsgInternalServerError   = "internal server error"
)
