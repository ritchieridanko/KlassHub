package utils

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/ritchieridanko/klasshub/services/user/internal/constants"
	"github.com/ritchieridanko/klasshub/services/user/internal/models"
	"go.opentelemetry.io/otel/trace"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var titlecaser = cases.Title(language.English)

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

// Convert timestamp to time
func ToTime(ts *timestamppb.Timestamp) *time.Time {
	if ts == nil {
		return nil
	}
	t := ts.AsTime()
	return &t
}

// Convert time to timestamp
func ToTimestamp(t *time.Time) *timestamppb.Timestamp {
	if t == nil {
		return nil
	}
	return timestamppb.New(*t)
}

// Set to all titlecase
func ToTitlecase(s string) string {
	return titlecaser.String(s)
}

// Set to all titlecase
// NOTE: Return nil if s is nil
func ToTitlecasePtr(s *string) *string {
	if s == nil {
		return nil
	}
	res := ToTitlecase(*s)
	return &res
}

// Strip string of leading and trailing whitespaces
// NOTE: Return nil if s is nil
func TrimSpacePtr(s *string) *string {
	if s == nil {
		return nil
	}
	res := strings.TrimSpace(*s)
	return &res
}
