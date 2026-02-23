package ce

import (
	"errors"

	"github.com/jackc/pgx/v5"
)

// Internal Errors
var (
	ErrDBAffectNoRows error = errors.New("no rows affected")
	ErrDBQueryNoRows  error = pgx.ErrNoRows
)

// Internal Error Codes
const (
	CodeAuthNotFound            errCode = "ERR_AUTH_NOT_FOUND"
	CodeDBQueryExec             errCode = "ERR_DB_QUERY_EXECUTION"
	CodeDBTransaction           errCode = "ERR_DB_TRANSACTION"
	CodeIdentifierNotRegistered errCode = "ERR_IDENTIFIER_NOT_REGISTERED"
	CodeInvalidIdentifier       errCode = "ERR_INVALID_IDENTIFIER"
	CodeInvalidPassword         errCode = "ERR_INVALID_PASSWORD"
	CodeInvalidRequestMeta      errCode = "ERR_INVALID_REQUEST_META"
	CodeJWTGenerationFailed     errCode = "ERR_JWT_GENERATION_FAILED"
	CodeUUIDGenerationFailed    errCode = "ERR_UUID_GENERATION_FAILED"
	CodeWrongPassword           errCode = "ERR_WRONG_PASSWORD"
)

// External Error Messages
const (
	MsgInternalServer     string = "Internal server error"
	MsgInvalidCredentials string = "Invalid credentials"
	MsgResourceNotFound   string = "Resource not found"
)
