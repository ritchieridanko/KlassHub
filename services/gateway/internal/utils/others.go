package utils

import (
	"context"
	"strings"
	"time"

	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/google/uuid"
	"github.com/ritchieridanko/klasshub/services/gateway/internal/constants"
	"go.opentelemetry.io/otel/trace"
)

// Get Request ID from Context
func CtxRequestID(ctx context.Context) string {
	if v, ok := ctx.Value(constants.CtxKeyRequestID).(string); ok {
		return v
	}
	return ""
}

// Get Subdomain from Context
func CtxSubdomain(ctx context.Context) string {
	if v, ok := ctx.Value(constants.CtxKeySubdomain).(string); ok {
		return v
	}
	return ""
}

// Get Trace ID from Context
func CtxTraceID(ctx context.Context) string {
	if sp := trace.SpanFromContext(ctx); sp.SpanContext().HasTraceID() {
		return sp.SpanContext().TraceID().String()
	}
	return ""
}

// Create a new random UUID v7
func GenerateUUIDv7() (uuid.UUID, error) {
	return uuid.NewV7()
}

// Strip string of leading and trailing whitespaces
// and set it to all lowercase
func NormalizeString(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}

// Convert timestamp value to time
func ToTime(ts *timestamp.Timestamp) *time.Time {
	if ts == nil {
		return nil
	}
	t := ts.AsTime()
	return &t
}
