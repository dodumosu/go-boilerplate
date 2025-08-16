package api

import (
	"context"
	"net/http"
)

type ContextKey string

const (
	CorrelationIDKey ContextKey = "correlationID"
	RequestKey       ContextKey = "httpRequest"
	UserIDKey        ContextKey = "userID"
)

func GetUserID(ctx context.Context) (string, bool) {
	val := ctx.Value(UserIDKey)
	if val == nil {
		return "", false
	}

	userID, ok := val.(string)
	return userID, ok && userID != ""
}

// GetCorrelationID retrieves the correlation ID from the context.
func GetCorrelationID(ctx context.Context) string {
	val := ctx.Value(CorrelationIDKey)
	if val == nil {
		return "" // Or a default one if you prefer
	}
	id, _ := val.(string)
	return id
}

// GetRequest retrieves the original *http.Request from the context.
// Returns the request and true if found.
func GetRequest(ctx context.Context) (*http.Request, bool) {
	val := ctx.Value(RequestKey)
	if val == nil {
		return nil, false
	}
	req, ok := val.(*http.Request)
	return req, ok
}

// WithUserID sets the userID in the context and returns the new context.
func WithUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, UserIDKey, userID)
}

// WithCorrelationID sets the correlationID in the context.
func WithCorrelationID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, CorrelationIDKey, id)
}

// WithRequest sets the *http.Request in the context.
func WithRequest(ctx context.Context, r *http.Request) context.Context {
	return context.WithValue(ctx, RequestKey, r)
}
