package tpgdclresource

import (
	"context"
	"fmt"

	glog "github.com/golang/glog"
)

// ReqCtxKey is the key type for storing values in the context.
// Context requires custom type key.
type ReqCtxKey string

// Keys used in context Value.
const (
	DoNotLogRequestsKey ReqCtxKey = "DoNotLogRequestsKey"
	APIRequestIDKey     ReqCtxKey = "APIRequestIDKey"
)

// APIRequestID returns the RequestID for the API call.
// APIRequestID is supposed to be used in log to help with debugging
// Since we do not want explicit error handling everywhere we are logging, so not throwing error.
// Its okay to print empty requestID in worse scenario.
func APIRequestID(ctx context.Context) string {
	val := ctx.Value(APIRequestIDKey)
	if val == nil {
		return ""
	}
	requestID, ok := val.(string)
	if !ok {
		glog.Warning("Could not convert APIRequestID val to string")
		return ""
	}
	return requestID
}

// ShouldLogRequest returns true if the request should be logged.
func ShouldLogRequest(ctx context.Context) (bool, error) {
	val := ctx.Value(DoNotLogRequestsKey)
	if val == nil {
		return true, nil
	}
	doNotLog, ok := val.(bool)
	if !ok {
		return false, fmt.Errorf("could not convert DoNotLogRequests value to bool")
	}
	return !doNotLog, nil
}

// ContextWithRequestID adds APIRequestID to ctx if APIRequestID is not present.
func ContextWithRequestID(ctx context.Context) context.Context {
	if APIRequestID(ctx) != "" {
		return ctx
	}
	return context.WithValue(ctx, APIRequestIDKey, CreateAPIRequestID())
}
