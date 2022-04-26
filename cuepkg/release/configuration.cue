package release

import (
	policy_v1beta1 "k8s.io/api/policy/v1beta1"
	autoscaling_v2beta1 "k8s.io/api/autoscaling/v2beta1"

)

_configuration: {
	#namespace: string

	spec: {
		horizontalPodAutoscalers: [Name = _]: autoscaling_v2beta1.#HorizontalPodAutoscaler & {
			metadata: name:      Name
			metadata: namespace: #namespace
		}

		podDisruptionBudgets: [Name = _]: policy_v1beta1.#PodDisruptionBudget & {
			metadata: name:      Name
			metadata: namespace: #namespace
		}
	}
}
