package ce

import "github.com/ritchieridanko/klasshub/services/notification/internal/infra/logger"

type errCode string

type Error struct {
	code   errCode
	err    error
	fields []logger.Field
}

func NewError(ec errCode, err error, fields ...logger.Field) *Error {
	return &Error{
		code:   ec,
		err:    err,
		fields: fields,
	}
}

func (e *Error) Code() errCode {
	return e.code
}

func (e *Error) Error() string {
	return e.err.Error()
}

func (e *Error) Fields() []logger.Field {
	return e.fields
}

func (e *Error) Append(fields ...logger.Field) *Error {
	e.fields = append(e.fields, fields...)
	return e
}
