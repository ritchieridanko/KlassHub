package ce

import (
	"net/http"

	"github.com/golang-jwt/jwt/v5"
)

// Internal Errors
var (
	ErrCookieNotFound  error = http.ErrNoCookie
	ErrInvalidJWTClaim error = jwt.ErrTokenInvalidClaims
	ErrJWTExpired      error = jwt.ErrTokenExpired
	ErrJWTMalformed    error = jwt.ErrTokenMalformed
)

// Internal Error Codes
const (
	CodeAlreadyExists          errCode = "ERR_ALREADY_EXISTS"
	CodeAuthNotVerified        errCode = "ERR_AUTH_NOT_VERIFIED"
	CodeInternal               errCode = "ERR_INTERNAL"
	CodeInvalidParams          errCode = "ERR_INVALID_PARAMS"
	CodeInvalidPayload         errCode = "ERR_INVALID_PAYLOAD"
	CodeInvalidRequestMetadata errCode = "ERR_INVALID_REQUEST_METADATA"
	CodeMissingContextValue    errCode = "ERR_MISSING_CONTEXT_VALUE"
	CodeNotFound               errCode = "ERR_NOT_FOUND"
	CodeRefreshTokenNotFound   errCode = "ERR_REFRESH_TOKEN_NOT_FOUND"
	CodeUnauthenticated        errCode = "ERR_UNAUTHENTICATED"
	CodeUnauthorized           errCode = "ERR_UNAUTHORIZED"
	CodeUnauthorizedRole       errCode = "ERR_UNAUTHORIZED_ROLE"
	CodeUnknown                errCode = "ERR_UNKNOWN"
	CodeUUIDGenerationFailed   errCode = "ERR_UUID_GENERATION_FAILED"
)

// External Error Messages
const (
	MsgAuthNotVerified  string = "Require account verification"
	MsgInternalServer   string = "Internal server error"
	MsgInvalidParams    string = "Invalid params"
	MsgInvalidPayload   string = "Invalid payload"
	MsgInvalidSession   string = "Invalid session"
	MsgResourceNotFound string = "Resource not found"
	MsgUnauthenticated  string = "Unauthenticated"
	MsgUnauthorized     string = "Unauthorized"
)
