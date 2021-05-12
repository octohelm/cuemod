package cuemoperator

import (
	"context"

	"github.com/octohelm/cuemod/pkg/kubernetes"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	"github.com/go-logr/logr"
	"github.com/octohelm/cuemod/pkg/cuex"
	"github.com/octohelm/cuemod/pkg/plugins/kube"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type ReleaseReconciler struct {
	Client    client.Client
	APIReader client.Reader
	Log       logr.Logger
	Scheme    *runtime.Scheme
}

func (r *ReleaseReconciler) SetupWithManager(mgr ctrl.Manager) error {
	gvks, err := allWatchableGroupVersionKinds(mgr.GetConfig())
	if err != nil {
		return err
	}

	c := ctrl.NewControllerManagedBy(mgr).For(&corev1.Secret{})

	for i := range gvks {
		u := &unstructured.Unstructured{}
		u.SetGroupVersionKind(gvks[i])

		c = c.Owns(u, builder.OnlyMetadata)
	}

	return c.Complete(r)
}

func (r *ReleaseReconciler) Reconcile(ctx context.Context, request reconcile.Request) (reconcile.Result, error) {
	releaseSecret := &corev1.Secret{}

	if err := r.APIReader.Get(ctx, request.NamespacedName, releaseSecret); err != nil {
		if apierrors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, err
	}

	if releaseSecret.Type != SecretTypeRelease {
		return reconcile.Result{}, nil
	}

	if template, ok := releaseSecret.Data["template"]; ok {
		ret, err := evalTemplate(template)
		if err != nil {
			r.Log.Error(err, "eval template failed")
			return reconcile.Result{}, nil
		}

		for i := range ret.Resources {
			o := ret.Resources[i]

			// skip namespace
			if o.GetObjectKind().GroupVersionKind().Kind == "Namespace" {
				continue
			}

			if err := controllerutil.SetControllerReference(releaseSecret, o, r.Scheme); err != nil {
				r.Log.Error(err, "set controller failed")
				return reconcile.Result{}, nil
			}

			if err := kubernetes.ApplyOne(ctx, r.Client, o, false); err != nil {
				r.Log.Error(err, "apply template failed")
				return reconcile.Result{}, nil
			}
		}
	}

	return reconcile.Result{}, nil
}

func evalTemplate(t []byte) (*kube.LoadResult, error) {
	i, err := cuex.InstanceFromRaw(t)
	if err != nil {
		return nil, err
	}
	jsonraw, err := cuex.Eval(i, cuex.JSON)
	if err != nil {
		return nil, err
	}
	r, err := kube.ReleaseFromJSONRaw(jsonraw)
	if err != nil {
		return nil, err
	}
	return kube.Process(r, nil)
}
