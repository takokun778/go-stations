package middleware

import (
	"context"

	ua "github.com/mileusna/useragent"
)

type key string

const (
	ctxOS = key("os")
)

func SetOSCtx(parent context.Context, userAgent string) context.Context {
	os := ua.Parse(userAgent)

	return context.WithValue(parent, ctxOS, os.OS)
}

func GetOSCtx(ctx context.Context) string {
	v := ctx.Value(ctxOS)

	ua, ok := v.(string)

	if !ok {
		return ""
	}

	return ua
}
