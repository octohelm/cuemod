package manifest

// This is a list of "cluster-wide" resources harvested from `kubectl api-resources --namespaced=false`
// This helps us to know which objects we should NOT apply namespaces to automatically.
// We can add to this list periodically if new types are added.
// This only applies to built-in kubernetes types. CRDs will need to be handled with annotations.
var clusterWideKinds = map[string]bool{
	"APIService":                     true,
	"CertificateSigningRequest":      true,
	"ClusterRole":                    true,
	"ClusterRoleBinding":             true,
	"ComponentStatus":                true,
	"CSIDriver":                      true,
	"CSINode":                        true,
	"CustomResourceDefinition":       true,
	"MutatingWebhookConfiguration":   true,
	"Namespace":                      true,
	"Node":                           true,
	"NodeMetrics":                    true,
	"PersistentVolume":               true,
	"PodSecurityPolicy":              true,
	"PriorityClass":                  true,
	"RuntimeClass":                   true,
	"SelfSubjectAccessReview":        true,
	"SelfSubjectRulesReview":         true,
	"StorageClass":                   true,
	"SubjectAccessReview":            true,
	"TokenReview":                    true,
	"ValidatingWebhookConfiguration": true,
	"VolumeAttachment":               true,
}

func IgnoreNamespace() func(m Object) Object {
	return func(m Object) Object {
		if m.GetObjectKind().GroupVersionKind().Kind == "Namespace" {
			return nil
		}
		return m
	}
}

func WithNamespace(def string) ProcessFunc {
	if def == "" {
		return nil
	}

	return func(m Object) Object {
		namespaced := true
		if clusterWideKinds[m.GetObjectKind().GroupVersionKind().Kind] {
			namespaced = false
		}
		if namespaced {
			m.SetNamespace(def)
		}
		return m
	}
}
