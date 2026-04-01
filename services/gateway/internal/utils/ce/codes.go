package ce

import "net/http"

// Internal Errors
var (
	ErrCookieNotFound error = http.ErrNoCookie
)

// Internal Error Codes
const (
	CodeAlreadyExists          errCode = "ERR_ALREADY_EXISTS"
	CodeInternal               errCode = "ERR_INTERNAL"
	CodeInvalidParams          errCode = "ERR_INVALID_PARAMS"
	CodeInvalidPayload         errCode = "ERR_INVALID_PAYLOAD"
	CodeInvalidRequestMetadata errCode = "ERR_INVALID_REQUEST_METADATA"
	CodeInvalidSubdomain       errCode = "ERR_INVALID_SUBDOMAIN"
	CodeMissingContextValue    errCode = "ERR_MISSING_CONTEXT_VALUE"
	CodeNotFound               errCode = "ERR_NOT_FOUND"
	CodeRefreshTokenNotFound   errCode = "ERR_REFRESH_TOKEN_NOT_FOUND"
	CodeUnauthenticated        errCode = "ERR_UNAUTHENTICATED"
	CodeUnauthorized           errCode = "ERR_UNAUTHORIZED"
	CodeUnknown                errCode = "ERR_UNKNOWN"
	CodeUUIDGenerationFailed   errCode = "ERR_UUID_GENERATION_FAILED"
)

// External Error Messages
const (
	MsgInternalServer   string = "Internal server error"
	MsgInvalidParams    string = "Invalid params"
	MsgInvalidPayload   string = "Invalid payload"
	MsgInvalidSession   string = "Invalid session"
	MsgInvalidSubdomain string = "Invalid host domain"
)
