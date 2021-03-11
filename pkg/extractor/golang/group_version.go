package golang

import (
	"context"

	"k8s.io/apimachinery/pkg/runtime/schema"
)

type contextGroupVersionKind struct {
}

func groupVersionKindFromContext(ctx context.Context) *schema.GroupVersionKind {
	if i, ok := ctx.Value(contextGroupVersionKind{}).(*schema.GroupVersionKind); ok {
		return i
	}
	return nil
}

func withGroupVersionKind(ctx context.Context, i *schema.GroupVersionKind) context.Context {
	return context.WithValue(ctx, contextGroupVersionKind{}, i)
}
