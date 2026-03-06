package constants

type ctxKey string

const (
	CtxKeyRequestID ctxKey = "x-request-id"
	CtxKeySubdomain ctxKey = "x-subdomain"
)
