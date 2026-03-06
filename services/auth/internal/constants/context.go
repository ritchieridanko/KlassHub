package constants

type ctxKey string

const (
	CtxKeyIPAddress ctxKey = "x-ip-address"
	CtxKeyRequestID ctxKey = "x-request-id"
	CtxKeySubdomain ctxKey = "x-subdomain"
	CtxKeyUserAgent ctxKey = "x-user-agent"
)
