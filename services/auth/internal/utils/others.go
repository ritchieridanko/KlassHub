package utils

import (
	"context"
	"strings"
	"time"

	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/google/uuid"
	"github.com/ritchieridanko/klasshub/services/auth/internal/constants"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
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

// Create a new random UUID v7
func GenerateUUIDv7() (uuid.UUID, error) {
	return uuid.NewV7()
}

// Strip string of leading and trailing whitespaces
// and set it to all lowercase
func NormalizeString(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}

// Wrap string value
func WrapString(s *string) *wrappers.StringValue {
	if s == nil {
		return nil
	}
	return wrapperspb.String(*s)
}

// Wrap time value
func WrapTime(t *time.Time) *timestamppb.Timestamp {
	if t == nil {
		return nil
	}
	return timestamppb.New(*t)
}
