package ce

import (
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/redis/go-redis/v9"
)

// Internal Errors
var (
	ErrCacheNoResult  error = redis.Nil
	ErrDBAffectNoRows error = errors.New("no rows affected")
	ErrDBQueryNoRows  error = pgx.ErrNoRows
	ErrWrongSubdomain error = errors.New("wrong subdomain")
)

// Internal Error Codes
const (
	CodeAuthAlreadyVerified     errCode = "ERR_AUTH_ALREADY_VERIFIED"
	CodeAuthNotFound            errCode = "ERR_AUTH_NOT_FOUND"
	CodeAuthNotRegistered       errCode = "ERR_AUTH_NOT_REGISTERED"
	CodeAuthNotVerified         errCode = "ERR_AUTH_NOT_VERIFIED"
	CodeBCryptHashingFailed     errCode = "ERR_BCRYPT_HASHING_FAILED"
	CodeCacheCommandExec        errCode = "ERR_CACHE_COMMAND_EXECUTION"
	CodeCacheScriptExec         errCode = "ERR_CACHE_SCRIPT_EXECUTION"
	CodeDBQueryExec             errCode = "ERR_DB_QUERY_EXECUTION"
	CodeDBTransaction           errCode = "ERR_DB_TRANSACTION"
	CodeEmailNotAvailable       errCode = "ERR_EMAIL_NOT_AVAILABLE"
	CodeEmailNotRegistered      errCode = "ERR_EMAIL_NOT_REGISTERED"
	CodeEventPublishingFailed   errCode = "ERR_EVENT_PUBLISHING_FAILED"
	CodeIdentifierNotRegistered errCode = "ERR_IDENTIFIER_NOT_REGISTERED"
	CodeInternal                errCode = "ERR_INTERNAL"
	CodeInvalidPayload          errCode = "ERR_INVALID_PAYLOAD"
	CodeJWTGenerationFailed     errCode = "ERR_JWT_GENERATION_FAILED"
	CodeMissingContextValue     errCode = "ERR_MISSING_CONTEXT_VALUE"
	CodeMissingMetadata         errCode = "ERR_MISSING_METADATA"
	CodeSessionNotFound         errCode = "ERR_SESSION_NOT_FOUND"
	CodeSessionNotOwned         errCode = "ERR_SESSION_NOT_OWNED"
	CodeTokenNotFound           errCode = "ERR_TOKEN_NOT_FOUND"
	CodeTokenNotOwned           errCode = "ERR_TOKEN_NOT_OWNED"
	CodeTypeConversionFailed    errCode = "ERR_TYPE_CONVERSION_FAILED"
	CodeUnauthenticated         errCode = "ERR_UNAUTHENTICATED"
	CodeUnauthorizedRole        errCode = "ERR_UNAUTHORIZED_ROLE"
	CodeUnauthorizedSubdomain   errCode = "ERR_UNAUTHORIZED_SUBDOMAIN"
	CodeUnknown                 errCode = "ERR_UNKNOWN"
	CodeUUIDGenerationFailed    errCode = "ERR_UUID_GENERATION_FAILED"
	CodeWrongPassword           errCode = "ERR_WRONG_PASSWORD"
	CodeWrongSubdomain          errCode = "ERR_WRONG_SUBDOMAIN"
)

// External Error Messages
const (
	MsgAuthAlreadyVerified string = "Account is already verified"
	MsgAuthNotFound        string = "Auth not found"
	MsgAuthNotVerified     string = "Require account verification"
	MsgEmailNotAvailable   string = "Email is already registered"
	MsgEmailNotRegistered  string = "Email is not registered"
	MsgInternalServer      string = "Internal server error"
	MsgInvalidCredentials  string = "Invalid credentials"
	MsgInvalidSession      string = "Invalid session"
	MsgInvalidToken        string = "Invalid token"
	MsgSessionNotFound     string = "Session not found"
	MsgUnauthenticated     string = "Unauthenticated"
	MsgUnauthorized        string = "Unauthorized"
)
