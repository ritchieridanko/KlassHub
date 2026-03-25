package ce

// Internal Error Codes
const (
	CodeAlreadyExists          errCode = "ERR_ALREADY_EXISTS"
	CodeInternal               errCode = "ERR_INTERNAL"
	CodeInvalidContextValue    errCode = "ERR_INVALID_CONTEXT_VALUE"
	CodeInvalidPayload         errCode = "ERR_INVALID_PAYLOAD"
	CodeInvalidRequestMetadata errCode = "ERR_INVALID_REQUEST_METADATA"
	CodeInvalidSubdomain       errCode = "ERR_INVALID_SUBDOMAIN"
	CodeNotFound               errCode = "ERR_NOT_FOUND"
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
	MsgInvalidSubdomain string = "Invalid host domain"
)
