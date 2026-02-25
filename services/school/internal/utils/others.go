package utils

import (
	"context"

	"github.com/ritchieridanko/klasshub/services/school/internal/constants"
	"go.opentelemetry.io/otel/trace"
)

// Get Request ID from Context
func CtxRequestID(ctx context.Context) string {
	if v, ok := ctx.Value(constants.CtxKeyRequestID).(string); ok {
		return v
	}
	return ""
}

// Get Request Meta (User Agent and IP Address) from Context
func CtxRequestMeta(ctx context.Context) (userAgent, ipAddress string) {
	userAgent, _ = ctx.Value(constants.CtxKeyUserAgent).(string)
	ipAddress, _ = ctx.Value(constants.CtxKeyIPAddress).(string)
	return
}

// Get Trace ID from Context
func CtxTraceID(ctx context.Context) string {
	if sp := trace.SpanFromContext(ctx); sp.SpanContext().HasTraceID() {
		return sp.SpanContext().TraceID().String()
	}
	return ""
}
