package utils

import (
	"context"
	"strings"
	"time"

	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/ritchieridanko/klasshub/services/school/internal/constants"
	"github.com/ritchieridanko/klasshub/services/school/internal/models"
	"github.com/segmentio/kafka-go"
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

// Get Request ID from Kafka Message
func MsgRequestID(msg kafka.Message) string {
	for _, header := range msg.Headers {
		if header.Key == constants.EventHeaderKeyRequestID {
			return string(header.Value)
		}
	}
	return ""
}

// Strip string of leading and trailing whitespaces
// and set it to all lowercase
func NormalizeString(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}

// Strip string of leading and trailing whitespaces
// and set it to all lowercase.
// NOTE: Return nil if s is nil
func NormalizeStringPtr(s *string) *string {
	if s == nil {
		return nil
	}
	res := NormalizeString(*s)
	return &res
}

// Convert timestamp to time
func ToTime(ts *timestamp.Timestamp) *time.Time {
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
