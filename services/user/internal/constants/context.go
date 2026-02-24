package constants

type ctxKey string

const (
	CtxKeyIPAddress ctxKey = "x-ip-address"
	CtxKeyRequestID ctxKey = "x-request-id"
	CtxKeyUserAgent ctxKey = "x-user-agent"
)
