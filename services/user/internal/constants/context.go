package constants

type ctxKey string

const (
	CtxKeyAuth      ctxKey = "x-auth"
	CtxKeyRequestID ctxKey = "x-request-id"
	CtxKeySubdomain ctxKey = "x-subdomain"
)
