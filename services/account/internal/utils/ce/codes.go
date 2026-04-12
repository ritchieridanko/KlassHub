package ce

// Internal Error Codes
const (
	CodeAlreadyExists         errCode = "ERR_ALREADY_EXISTS"
	CodeAuthNotVerified       errCode = "ERR_AUTH_NOT_VERIFIED"
	CodeEventPublishingFailed errCode = "ERR_EVENT_PUBLISHING_FAILED"
	CodeFailedPrecondition    errCode = "ERR_FAILED_PRECONDITION"
	CodeInternal              errCode = "ERR_INTERNAL"
	CodeInvalidArgument       errCode = "ERR_INVALID_ARGUMENT"
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
)

// External Error Messages
const (
	MsgAuthNotVerified    string = "Require account verification"
	MsgInternalServer     string = "Internal server error"
	MsgInvalidCredentials string = "Invalid credentials"
	MsgUnauthenticated    string = "Unauthenticated"
	MsgUnauthorized       string = "Unauthorized"
)
