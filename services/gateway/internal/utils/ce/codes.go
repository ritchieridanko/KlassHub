package ce

// Internal Error Codes
const (
	CodeInternal               errCode = "ERR_INTERNAL"
	CodeInvalidPayload         errCode = "ERR_INVALID_PAYLOAD"
	CodeInvalidRequestMetadata errCode = "ERR_INVALID_REQUEST_METADATA"
	CodeNotFound               errCode = "ERR_NOT_FOUND"
	CodeUnauthenticated        errCode = "ERR_UNAUTHENTICATED"
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
