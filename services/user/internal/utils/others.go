package utils

import (
	"context"
	"time"

	"github.com/ritchieridanko/klasshub/services/user/internal/constants"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Get Request ID from Context
func CtxRequestID(ctx context.Context) string {
	if v, ok := ctx.Value(constants.CtxKeyRequestID).(string); ok {
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

// Convert time value to timestamp
func ToTimestamp(t *time.Time) *timestamppb.Timestamp {
	if t == nil {
		return nil
	}
	return timestamppb.New(*t)
}
