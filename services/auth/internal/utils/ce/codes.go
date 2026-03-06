package ce

import (
	"errors"

	"github.com/jackc/pgx/v5"
)

// Internal Errors
var (
	ErrDBAffectNoRows error = errors.New("no rows affected")
	ErrDBQueryNoRows  error = pgx.ErrNoRows
	ErrWrongSubdomain error = errors.New("wrong subdomain")
)

// Internal Error Codes
const (
	CodeAuthNotFound            errCode = "ERR_AUTH_NOT_FOUND"
	CodeDBQueryExec             errCode = "ERR_DB_QUERY_EXECUTION"
	CodeDBTransaction           errCode = "ERR_DB_TRANSACTION"
	CodeIdentifierNotRegistered errCode = "ERR_IDENTIFIER_NOT_REGISTERED"
	CodeInternal                errCode = "ERR_INTERNAL"
	CodeInvalidPayload          errCode = "ERR_INVALID_PAYLOAD"
	CodeJWTGenerationFailed     errCode = "ERR_JWT_GENERATION_FAILED"
	CodeUnknown                 errCode = "ERR_UNKNOWN"
	CodeUUIDGenerationFailed    errCode = "ERR_UUID_GENERATION_FAILED"
	CodeWrongPassword           errCode = "ERR_WRONG_PASSWORD"
	CodeWrongSubdomain          errCode = "ERR_WRONG_SUBDOMAIN"
)

// External Error Messages
const (
	MsgAuthNotFound       string = "Auth not found"
	MsgInternalServer     string = "Internal server error"
	MsgInvalidCredentials string = "Invalid credentials"
)
