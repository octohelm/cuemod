package kubernetes

import (
	"context"
	"fmt"
	"strings"

	releasev1alpha1 "github.com/octohelm/cuemod/pkg/apis/release/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"

	"github.com/octohelm/cuemod/pkg/kubernetes/manifest"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"

	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var (
	scheme = runtime.NewScheme()
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(releasev1alpha1.AddToScheme(scheme))
}

func NewClientForContext(contextName string) (*KubeClient, error) {
	kubeConfig, err := ResolveKubeConfig(contextName)
	if err != nil {
		return nil, err
	}

	c, e := client.New(kubeConfig, client.Options{
		Scheme: scheme,
	})
	if e != nil {
		return nil, e
	}

	return &KubeClient{
		Client: c,
		Config: kubeConfig,
	}, nil
}

func NewClient() (*KubeClient, error) {
	return NewClientForContext("")
}

type KubeClient struct {
	client.Client
	Config *rest.Config
}

func (c *KubeClient) Info() string {
	return fmt.Sprintf("cluster '%s'", c.Config.Host)
}

func (c *KubeClient) AllListableGroupVersionKinds() (gvks []schema.GroupVersionKind, err error) {
	return allListableGroupVersionKinds(c.Config)
}

func (c *KubeClient) ListAll(ctx context.Context, groupVersionKinds []schema.GroupVersionKind, listOptions ...client.ListOption) (l manifest.List, err error) {
	for _, gvk := range groupVersionKinds {
		// skip none list
		rolist, err := manifest.NewListForGroupVersionKind(gvk)
		if err != nil {
			continue
		}

		if err := c.List(ctx, rolist, listOptions...); err != nil {
			return nil, err
		}

		if meta.LenList(rolist) == 0 {
			continue
		}

		list, err := meta.ExtractList(rolist)
		if err != nil {
			return nil, err
		}

		for i := range list {
			o, _ := manifest.ObjectFromRuntimeObject(list[i])

			// skip resources controllered by other
			if len(o.GetOwnerReferences()) == 0 {
				l = append(l, o)
			}
		}
	}

	return
}

func allListableGroupVersionKinds(conf *rest.Config) (gvks []schema.GroupVersionKind, err error) {
	dc, e := discovery.NewDiscoveryClientForConfig(conf)
	if e != nil {
		return nil, e
	}

	preferredResources, err := dc.ServerPreferredResources()
	if err != nil {
		return nil, err
	}

	for i := range preferredResources {
		pr := preferredResources[i]

		if len(pr.APIResources) == 0 {
			continue
		}

		gv, err := schema.ParseGroupVersion(pr.GroupVersion)
		if err != nil {
			continue
		}

		for _, resource := range pr.APIResources {
			if len(resource.Verbs) == 0 {
				continue
			}

			if strings.Contains(strings.Join(resource.Verbs, ","), "list") {
				gvks = append(gvks, gv.WithKind(resource.Kind))
			}
		}
	}
	return
}
