package handler

import (
	"context"
	"errors"
	"net/http"

	"wata-bot-BE/internal/model"
	"wata-bot-BE/internal/types"
	"wata-bot-BE/internal/utils"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// ErrorHandler wraps httpx.ErrorCtx to log errors and return formatted error response
func ErrorHandler(ctx context.Context, w http.ResponseWriter, err error) {
	// Log all errors to file
	context := map[string]interface{}{
		"error": err.Error(),
	}
	utils.WriteErrorLogWithContext("Handler error", err, context)
	logx.Errorf("Handler error: %v", err)

	// Check if error is APIError with error code
	var apiErr *model.APIError
	if errors.As(err, &apiErr) {
		// Return formatted error response with error code
		// Use BadRequest for client errors, InternalServerError for server errors
		statusCode := http.StatusBadRequest
		if apiErr.Code == model.ErrCodeInternalServerError {
			statusCode = http.StatusInternalServerError
		}
		httpx.WriteJsonCtx(ctx, w, statusCode, types.ErrorResp{
			ErrorCode: apiErr.Code,
			Message:   apiErr.Message,
		})
		return
	}

	// Default error response for unknown errors
	httpx.WriteJsonCtx(ctx, w, http.StatusInternalServerError, types.ErrorResp{
		ErrorCode: model.ErrCodeInternalServerError,
		Message:   err.Error(),
	})
}
