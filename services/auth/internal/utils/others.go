package utils

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/ritchieridanko/klasshub/services/auth/internal/constants"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Get IP Address from Context
func CtxIPAddress(ctx context.Context) string {
	if v, ok := ctx.Value(constants.CtxKeyIPAddress).(string); ok {
		return v
	}
	return ""
}

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

// Get User Agent from Context
func CtxUserAgent(ctx context.Context) string {
	if v, ok := ctx.Value(constants.CtxKeyUserAgent).(string); ok {
		return v
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

// Convert time value to timestamp
func ToTimestamp(t *time.Time) *timestamppb.Timestamp {
	if t == nil {
		return nil
	}
	return timestamppb.New(*t)
}
