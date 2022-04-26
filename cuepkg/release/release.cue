package release

#Release: {
	// release name
	#name: string

	// releaee namespace
	#namespace: string

	// context name
	#context: string

	apiVersion: "octohelm.tech/v1alpha1"
	kind:       "Release"

	metadata: name:      "\(#name)"
	metadata: namespace: "\(#namespace)"
	metadata: labels: [K=string]: string

	metadata: labels: context: "\(#context)"

	for m in [_namespace, _rbac, _configuration, _storage, _network, _workload] {
		m & {#namespace: #namespace}
	}

	// wild kube resources
	spec: kube?: _
}
