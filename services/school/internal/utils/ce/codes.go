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
	CodeEventFetchingFailed   errCode = "ERR_EVENT_FETCHING_FAILED"
	CodeEventCommittingFailed errCode = "ERR_EVENT_COMMITTING_FAILED"
	CodeInvalidContextValue   errCode = "ERR_INVALID_CONTEXT_VALUE"
	CodeInvalidPayload        errCode = "ERR_INVALID_PAYLOAD"
	CodeMissingContextValue   errCode = "ERR_MISSING_CONTEXT_VALUE"
	CodeMissingMetadata       errCode = "ERR_MISSING_METADATA"
	CodePanicOccurred         errCode = "ERR_PANIC_OCCURRED"
	CodeProtobufParsingFailed errCode = "ERR_PROTOBUF_PARSING_FAILED"
	CodeSchoolNotFound        errCode = "ERR_SCHOOL_NOT_FOUND"
	CodeTypeConversionFailed  errCode = "ERR_TYPE_CONVERSION_FAILED"
	CodeUnauthenticated       errCode = "ERR_UNAUTHENTICATED"
	CodeUnauthorizedRole      errCode = "ERR_UNAUTHORIZED_ROLE"
	CodeUnauthorizedSubdomain errCode = "ERR_UNAUTHORIZED_SUBDOMAIN"
	CodeUnknown               errCode = "ERR_UNKNOWN"
)

// External Error Messages
const (
	MsgAuthNotVerified string = "Require account verification"
	MsgInternalServer  string = "Internal server error"
	MsgSchoolNotFound  string = "School not found"
	MsgUnauthenticated string = "Unauthenticated"
	MsgUnauthorized    string = "Unauthorized"
)
