package cuemoperator

import (
	"context"

	releasev1alpha1 "github.com/octohelm/cuemod/pkg/apis/release/v1alpha1"

	"github.com/octohelm/cuemod/pkg/kubernetes"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	"github.com/go-logr/logr"
	"github.com/octohelm/cuemod/pkg/plugins/kube"
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

	c := ctrl.NewControllerManagedBy(mgr).For(&releasev1alpha1.Release{})

	for i := range gvks {
		if gvks[i].Group == releasev1alpha1.SchemeGroupVersion.Group {
			continue
		}

		u := &unstructured.Unstructured{}
		u.SetGroupVersionKind(gvks[i])
		c = c.Owns(u, builder.OnlyMetadata)
	}

	return c.Complete(r)
}

func (r *ReleaseReconciler) Reconcile(ctx context.Context, request reconcile.Request) (reconcile.Result, error) {
	release := &releasev1alpha1.Release{}

	if err := r.APIReader.Get(ctx, request.NamespacedName, release); err != nil {
		if apierrors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, err
	}

	ret, err := kube.Process(release, nil)
	if err != nil {
		r.Log.Error(err, "eval template failed")
		return reconcile.Result{}, nil
	}

	release.Status.Resources = nil

	for i := range ret.Resources {
		o := ret.Resources[i]

		// skip namespace
		if o.GetObjectKind().GroupVersionKind().Kind == "Namespace" {
			continue
		}

		if err := controllerutil.SetControllerReference(release, o, r.Scheme); err != nil {
			r.Log.Error(err, "set controller failed")
			return reconcile.Result{}, nil
		}

		if err := kubernetes.ApplyOne(ctx, r.Client, o, false); err != nil {
			r.Log.Error(err, "apply template failed")
			return reconcile.Result{}, nil
		}

		release.Status.Resources = append(release.Status.Resources, releasev1alpha1.Resource{
			GroupVersionKind: o.GetObjectKind().GroupVersionKind(),
			Namespace:        o.GetNamespace(),
			Name:             o.GetName(),
		})

		if err := r.Client.Status().Update(ctx, release); err != nil {
			r.Log.Error(err, "update status err")
			return reconcile.Result{}, nil
		}
	}

	return reconcile.Result{}, nil
}
