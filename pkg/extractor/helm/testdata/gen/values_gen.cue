// DO NOT EDIT THIS FILE DIRECTLY.
// generated by helm extractor.
package chart

values: {
	affinity?: [string]: _
	autoscaling?: {
		enabled?:                        *false | bool
		maxReplicas?:                    *100 | float64
		minReplicas?:                    *1 | float64
		targetCPUUtilizationPercentage?: *80 | float64
	}
	fullnameOverride?: *"" | string
	image?: {
		pullPolicy?: *"IfNotPresent" | string
		repository?: *"nginx" | string
		tag?:        *"" | string
	}
	imagePullSecrets?: [...]
	ingress?: {
		annotations?: [string]: _
		enabled?: *false | bool
		hosts?: [...{
			host?: *"chart-example.local" | string
			paths?: [...{
				backend?: {
					serviceName?: *"chart-example.local" | string
					servicePort?: *80 | float64
				}
				path?: *"/" | string
			}]
		}]
		tls?: [...]
	}
	nameOverride?: *"" | string
	nodeSelector?: [string]: _
	podAnnotations?: [string]: _
	podSecurityContext?: [string]: _
	replicaCount?: *1 | float64
	resources?: [string]: _
	securityContext?: [string]: _
	service?: {
		port?: *80 | float64
		type?: *"ClusterIP" | string
	}
	serviceAccount?: {
		annotations?: [string]: _
		create?: *true | bool
		name?:   *"" | string
	}
	tolerations?: [...]
}
