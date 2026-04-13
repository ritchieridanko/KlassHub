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
	CodeAlreadyExists         errCode = "ERR_ALREADY_EXISTS"
	CodeAuthNotVerified       errCode = "ERR_AUTH_NOT_VERIFIED"
	CodeDBQueryExec           errCode = "ERR_DB_QUERY_EXECUTION"
	CodeDBTransaction         errCode = "ERR_DB_TRANSACTION"
	CodeFailedPrecondition    errCode = "ERR_FAILED_PRECONDITION"
	CodeInternal              errCode = "ERR_INTERNAL"
	CodeInvalidArgument       errCode = "ERR_INVALID_ARGUMENT"
	CodeInvalidPayload        errCode = "ERR_INVALID_PAYLOAD"
	CodeMissingContextValue   errCode = "ERR_MISSING_CONTEXT_VALUE"
	CodeMissingMetadata       errCode = "ERR_MISSING_METADATA"
	CodeNotFound              errCode = "ERR_NOT_FOUND"
	CodePermissionDenied      errCode = "ERR_PERMISSION_DENIED"
	CodeSchoolNotRegistered   errCode = "ERR_SCHOOL_NOT_REGISTERED"
	CodeTypeConversionFailed  errCode = "ERR_TYPE_CONVERSION_FAILED"
	CodeUnauthenticated       errCode = "ERR_UNAUTHENTICATED"
	CodeUnauthorizedRole      errCode = "ERR_UNAUTHORIZED_ROLE"
	CodeUnauthorizedSubdomain errCode = "ERR_UNAUTHORIZED_SUBDOMAIN"
	CodeUnknown               errCode = "ERR_UNKNOWN"
	CodeUUIDGenerationFailed  errCode = "ERR_UUID_GENERATION_FAILED"
)

// External Error Messages
const (
	MsgAuthNotVerified    string = "Require account verification"
	MsgInternalServer     string = "Internal server error"
	MsgInvalidCredentials string = "Invalid credentials"
	MsgUnauthenticated    string = "Unauthenticated"
	MsgUnauthorized       string = "Unauthorized"
)
