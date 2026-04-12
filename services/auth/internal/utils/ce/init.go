package ce

import (
	"github.com/ritchieridanko/klasshub/services/auth/internal/infra/logger"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

func (e *Error) ToGRPCStatus() error {
	switch e.code {
	case CodeIdentifierNotProvided, CodeInvalidPayload, CodeTokenNotFound,
		CodeTokenNotOwned:
		return status.Error(codes.InvalidArgument, e.message)
	case CodeAuthNotFound, CodeSessionNotFound:
		return status.Error(codes.NotFound, e.message)
	case CodeEmailNotAvailable, CodeUsernameNotAvailable:
		return status.Error(codes.AlreadyExists, e.message)
	case CodeUnauthorizedRole, CodeUnauthorizedSubdomain:
		return status.Error(codes.PermissionDenied, e.message)
	case CodeAuthAlreadyVerified, CodeAuthNotVerified, CodeEmailNotRegistered:
		return status.Error(codes.FailedPrecondition, e.message)
	case CodeAuthNotRegistered, CodeIdentifierNotRegistered, CodeSessionExpired,
		CodeSessionNotOwned, CodeUnauthenticated, CodeWrongPassword, CodeWrongSubdomain:
		return status.Error(codes.Unauthenticated, e.message)
	case CodeBCryptHashingFailed, CodeCacheCommandExec, CodeCacheScriptExec,
		CodeDBQueryExec, CodeDBTransaction, CodeEventPublishingFailed, CodeInternal,
		CodeJWTGenerationFailed, CodeMissingContextValue, CodeMissingMetadata,
		CodeTypeConversionFailed, CodeUnknown, CodeUUIDGenerationFailed:
		return status.Error(codes.Internal, e.message)
	default:
		return status.Error(codes.Internal, e.message)
	}
}
