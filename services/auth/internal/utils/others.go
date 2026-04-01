package utils

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/ritchieridanko/klasshub/services/auth/internal/constants"
	"github.com/ritchieridanko/klasshub/services/auth/internal/models"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Get Auth Information from Context
func CtxAuth(ctx context.Context) *models.AuthContext {
	if v, ok := ctx.Value(constants.CtxKeyAuth).(*models.AuthContext); ok {
		return v
	}
	return nil
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

// Get Transport Information from Context
func CtxTransport(ctx context.Context) *models.TransportContext {
	if v, ok := ctx.Value(constants.CtxKeyTransport).(*models.TransportContext); ok {
		return v
	}
	return nil
}

// Create a new random UUID
func GenerateUUID() uuid.UUID {
	return uuid.New()
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

// Convert value to int64
func ToInt64(value any) (int64, error) {
	if v, ok := value.(string); ok {
		return strconv.ParseInt(v, 10, 64)
	}
	return 0, fmt.Errorf("unable to convert to int64: %v", value)
}

// Convert time to timestamp
func ToTimestamp(t *time.Time) *timestamppb.Timestamp {
	if t == nil {
		return nil
	}
	return timestamppb.New(*t)
}
