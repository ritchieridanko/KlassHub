package utils

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"net/url"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/ritchieridanko/klasshub/services/notification/internal/constants"
	"github.com/segmentio/kafka-go"
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

// Create a new tokenized URL
func GenerateTokenizedURL(baseURL, path, token string) (string, error) {
	u, err := url.Parse(baseURL)
	if err != nil {
		return "", err
	}

	u.Path = path
	q := u.Query()
	q.Set("token", token)

	u.RawQuery = q.Encode()
	return u.String(), nil
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

// Encode value to MIME Base64
func ToMIMEBase64(value string) string {
	return "=?UTF-8?B?" + base64.StdEncoding.EncodeToString([]byte(value)) + "?="
}

// Convert value to json.RawMessage
func ToRawMessage(value any) (json.RawMessage, error) {
	b, err := json.Marshal(value)
	if err != nil {
		return nil, err
	}
	return json.RawMessage(b), nil
}

// Convert timestamp to time
func ToTime(ts *timestamppb.Timestamp) *time.Time {
	if ts == nil {
		return nil
	}
	t := ts.AsTime()
	return &t
}

// Convert string to UUID
func ToUUID(value string) (uuid.UUID, error) {
	return uuid.Parse(value)
}
