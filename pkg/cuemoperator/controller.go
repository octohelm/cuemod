package cuemoperator

import (
	"context"

	"sigs.k8s.io/controller-runtime/pkg/client"
)

func ClientWithoutCache(c client.Client, r client.Reader) client.Client {
	return &clientWithoutCache{Client: c, r: r}
}

type clientWithoutCache struct {
	r client.Reader
	client.Client
}

func (c *clientWithoutCache) Get(ctx context.Context, key client.ObjectKey, obj client.Object) error {
	return c.r.Get(ctx, key, obj)
}

func (c *clientWithoutCache) List(ctx context.Context, key client.ObjectList, opts ...client.ListOption) error {
	return c.r.List(ctx, key, opts...)
}
