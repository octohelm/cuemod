package release

import (
	"k8s.io/api/core/v1"

	rbac_v1 "k8s.io/api/rbac/v1"
)

_rbac: {
	#namespace: string

	spec: {
		serviceAccounts: [Name=_]: v1.#ServiceAccount & {
			metadata: name:      Name
			metadata: namespace: "\(#namespace)"

			// auto create ClusterRole or Role and RoleBindings with #role and #rules
			#role: "ClusterRole" | "Role"
			#rules: [...rbac_v1.#PolicyRule]
		}

		for n, sa in serviceAccounts {
			if sa.#role == "ClusterRole" {
				{
					clusterRoles: "\(n)": rbac_v1.#ClusterRole & {
						metadata: name:      "\(n)"
						metadata: namespace: "\(#namespace)"
						rules: sa.#rules
					}

					clusterRoleBindings: "\(n)": rbac_v1.#ClusterRoleBinding & {
						metadata: name:      "\(n)"
						metadata: namespace: "\(#namespace)"

						subjects: [{
							kind:      "ServiceAccount"
							name:      "\(n)"
							namespace: "\(#namespace)"
						}]

						roleRef: {
							kind:     "ClusterRole"
							name:     "\(n)"
							apiGroup: "rbac.authorization.k8s.io"
						}
					}
				}
			}

			if sa.#role == "Role" {
				{
					roles: "\(n)": rbac_v1.#Role & {
						metadata: name:      "\(n)"
						metadata: namespace: "\(#namespace)"
						rules: sa.#rules
					}

					roleBindings: "\(n)": rbac_v1.#RoleBinding & {
						metadata: name:      "\(n)"
						metadata: namespace: "\(#namespace)"

						subjects: [
							{
								kind:      "ServiceAccount"
								name:      "\(n)"
								namespace: "\(#namespace)"
							},
						]

						roleRef: {
							kind:     "Role"
							name:     "\(n)"
							apiGroup: "rbac.authorization.k8s.io"
						}
					}
				}
			}
		}
	}
}
