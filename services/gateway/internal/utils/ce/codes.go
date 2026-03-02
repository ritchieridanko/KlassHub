package ce

// Internal Error Codes
const (
	CodeAlreadyExists        errCode = "ERR_ALREADY_EXISTS"
	CodeInternal             errCode = "ERR_INTERNAL"
	CodeInvalidPayload       errCode = "ERR_INVALID_PAYLOAD"
	CodeInvalidSubdomain     errCode = "ERR_INVALID_SUBDOMAIN"
	CodeNotFound             errCode = "ERR_NOT_FOUND"
	CodeUnauthenticated      errCode = "ERR_UNAUTHENTICATED"
	CodeUUIDGenerationFailed errCode = "ERR_UUID_GENERATION_FAILED"
	CodeUnknown              errCode = "ERR_UNKNOWN"
)

// External Error Messages
const (
	MsgInternalServer   string = "Internal server error"
	MsgInvalidParams    string = "Invalid params"
	MsgInvalidPayload   string = "Invalid payload"
	MsgInvalidSubdomain string = "Invalid subdomain"
)
