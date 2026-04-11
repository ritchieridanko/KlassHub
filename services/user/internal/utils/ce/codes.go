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
	CodeAuthNotVerified       errCode = "ERR_AUTH_NOT_VERIFIED"
	CodeDBQueryExec           errCode = "ERR_DB_QUERY_EXECUTION"
	CodeDBTransaction         errCode = "ERR_DB_TRANSACTION"
	CodeInvalidPayload        errCode = "ERR_INVALID_PAYLOAD"
	CodeMissingContextValue   errCode = "ERR_MISSING_CONTEXT_VALUE"
	CodeMissingMetadata       errCode = "ERR_MISSING_METADATA"
	CodeTypeConversionFailed  errCode = "ERR_TYPE_CONVERSION_FAILED"
	CodeUnauthenticated       errCode = "ERR_UNAUTHENTICATED"
	CodeUnauthorizedRole      errCode = "ERR_UNAUTHORIZED_ROLE"
	CodeUnauthorizedSubdomain errCode = "ERR_UNAUTHORIZED_SUBDOMAIN"
	CodeUnknown               errCode = "ERR_UNKNOWN"
	CodeUserNotFound          errCode = "ERR_USER_NOT_FOUND"
	CodeUserNotRegistered     errCode = "ERR_USER_NOT_REGISTERED"
	CodeUUIDGenerationFailed  errCode = "ERR_UUID_GENERATION_FAILED"
)

// External Error Messages
const (
	MsgAuthNotVerified    string = "Require account verification"
	MsgInternalServer     string = "Internal server error"
	MsgInvalidCredentials string = "Invalid credentials"
	MsgUnauthenticated    string = "Unauthenticated"
	MsgUnauthorized       string = "Unauthorized"
	MsgUserNotFound       string = "User not found"
)
