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

	// Transaction errors (0300-0399)
	ErrCodeInvalidCurrency       = "0300"
	ErrCodeInvalidAmount         = "0301"
	ErrCodeInsufficientBalance   = "0302"
	ErrCodeFailedToUpdateBalance = "0303"

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
	ErrMsgInvalidCurrency       = "invalid currency. Must be 'wata' or 'usdt'"
	ErrMsgInvalidAmount         = "invalid amount"
	ErrMsgInsufficientBalance   = "insufficient balance"
	ErrMsgFailedToUpdateBalance = "failed to update balance"
	ErrMsgInternalServerError   = "internal server error"
)
