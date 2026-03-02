package utils

import (
	"context"
	"errors"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/google/uuid"
	"github.com/ritchieridanko/klasshub/services/gateway/internal/constants"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/metadata"
)

// Get Request ID from Context
func CtxRequestID(ctx context.Context) string {
	if v, ok := ctx.Value(constants.CtxKeyRequestID).(string); ok {
		return v
	}
	return ""
}

// Get Subdomain from Context
func CtxSubdomain(ctx *gin.Context, hostname, tld string) (string, error) {
	host, _, err := net.SplitHostPort(ctx.Request.Host)
	if err != nil {
		return "", err
	}
	if suffix := fmt.Sprintf(".%s.%s", hostname, tld); strings.HasSuffix(host, suffix) {
		return strings.TrimSuffix(host, suffix), nil
	}
	return "", errors.New("invalid hostname")
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

// Insert metadata into outgoing context.
// Transport details include User Agent and IP Address
func ToOutgoingCtx(ctx *gin.Context, includeTransportDetails bool) context.Context {
	c := ctx.Request.Context()

	pairs := make([]string, 0, 6)
	pairs = append(pairs, constants.MDKeyRequestID, CtxRequestID(c))

	if includeTransportDetails {
		pairs = append(pairs, constants.MDKeyUserAgent, ctx.Request.UserAgent())
		pairs = append(pairs, constants.MDKeyIPAddress, ctx.ClientIP())
	}
	return metadata.AppendToOutgoingContext(c, pairs...)
}

// Unwrap string value
func UnwrapString(sv *wrappers.StringValue) *string {
	if sv == nil {
		return nil
	}
	return &sv.Value
}

// Unwrap timestamp value
func UnwrapTimestamp(ts *timestamp.Timestamp) *time.Time {
	if ts == nil {
		return nil
	}
	t := ts.AsTime()
	return &t
}
