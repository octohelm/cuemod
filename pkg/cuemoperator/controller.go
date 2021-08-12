package cuemoperator

import (
	"context"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

func SetupWithManager(mgr ctrl.Manager) error {
	return SetupReconcilerWithManager(
		mgr,
		&ReleaseReconciler{
			APIReader: mgr.GetAPIReader(),
			Client:    mgr.GetClient(),
			Log:       mgr.GetLogger().WithName("controllers").WithName("Release"),
			Scheme:    mgr.GetScheme(),
		},
	)
}

type Reconciler interface {
	SetupWithManager(mgr ctrl.Manager) error
}

func SetupReconcilerWithManager(mgr manager.Manager, reconcilers ...Reconciler) error {
	for i := range reconcilers {
		if err := reconcilers[i].SetupWithManager(mgr); err != nil {
			return err
		}
	}
	return nil
}

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
