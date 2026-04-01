package ce

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ritchieridanko/klasshub/services/gateway/internal/infra/logger"
)

type errCode string

type Error struct {
	code    errCode
	message string
	err     error
	fields  []logger.Field
}

func NewError(ec errCode, message string, err error, fields ...logger.Field) *Error {
	return &Error{
		code:    ec,
		message: message,
		err:     err,
		fields:  fields,
	}
}

func (e *Error) Code() errCode {
	return e.code
}

func (e *Error) Message() string {
	return e.message
}

func (e *Error) Error() string {
	if e.err != nil {
		return e.message + ": " + e.err.Error()
	}
	return e.message
}

func (e *Error) Fields() []logger.Field {
	return e.fields
}

func (e *Error) Unwrap() error {
	return e.err
}

func (e *Error) Append(fields ...logger.Field) *Error {
	e.fields = append(e.fields, fields...)
	return e
}

func (e *Error) Bind(ctx *gin.Context) {
	ctx.Error(e)
}

func (e *Error) ToHTTPStatus() int {
	switch e.code {
	case CodeInvalidParams, CodeInvalidPayload, CodeInvalidRequestMetadata:
		return http.StatusBadRequest
	case CodeRefreshTokenNotFound, CodeUnauthenticated:
		return http.StatusUnauthorized
	case CodeUnauthorized:
		return http.StatusForbidden
	case CodeNotFound:
		return http.StatusNotFound
	case CodeAlreadyExists:
		return http.StatusConflict
	case CodeInternal, CodeMissingContextValue, CodeUnknown,
		CodeUUIDGenerationFailed:
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}
