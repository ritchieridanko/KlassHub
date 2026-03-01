package ce

import (
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

func (e *Error) AppendFields(fields ...logger.Field) *Error {
	e.fields = append(e.fields, fields...)
	return e
}

func (e *Error) Bind(ctx *gin.Context) {
	ctx.Error(e)
}
