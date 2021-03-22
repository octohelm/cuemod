package release

import (
	"k8s.io/api/core/v1"

	apps_v1 "k8s.io/api/apps/v1"
	networking_v1 "k8s.io/api/networking/v1"
)

#ReleaseBase: {
	apiVersion: "octohelm.tech/v1alpha"
	kind:       "Release"

	#name:      string
	#namespace: string
	#context:   *"default" | string

	metadata: name:      "\(#name)"
	metadata: namespace: "\(#namespace)"
	metadata: labels: context: "\(#context)"

	spec: _
}

#Release: {
	#ReleaseBase

	#name:      string
	#namespace: string
	#context:   *"default" | string

	metadata: name:      "\(#name)"
	metadata: namespace: "\(#namespace)"
	metadata: labels: context: "\(#context)"

	spec: {
		namespace: v1.#Namespace & {
			metadata: name: "\(#namespace)"
		}

		pvc: [Name = _]: v1.#PersistentVolumeClaim & {
			metadata: name:      Name
			metadata: namespace: "\(#namespace)"
			metadata: labels: app: Name
		}

		configmap: [Name = _]: v1.#ConfigMap & {
			metadata: name:      Name
			metadata: namespace: "\(#namespace)"
		}

		deployment: [Name = _]: apps_v1.#Deployment & {
			metadata: name:      Name
			metadata: namespace: "\(#namespace)"

			spec: selector: matchLabels: app: Name
			spec: template: metadata: labels: app: Name
			spec: replicas: *1 | int

			#volumes: [string]: {
				v1.#VolumeMount

				volume: v1.#Volume
			}

			#container: v1.#Container & {

				#ports: [string]: int

				ports: [
					for n, cp in #ports {
						{
							name:          n
							containerPort: cp
						}
					},
				]

				#env: [string]: string

				env: [
					for n, v in #env {
						{
							name:  "\(n)"
							value: v
						}
					},
				]

				volumeMounts: [
					for n, vol in #volumes {
						name: n
						for k, v in vol if k != "volume" {
							"\(k)": v
						}
					},
				]
			}

			spec: template: spec: volumes: [
				for n, v in #volumes {
					{
						name: "\(n)"
						v.volume
					}
				},
			]

			spec: template: spec: containers: [
				#container,
			]
		}

		for x in [deployment] for n, v in x {
			service: "\(n)": v1.#Service & {
				metadata: name:      "\(n)"
				metadata: namespace: "\(#namespace)"

				spec: selector: app: v.spec.template.metadata.labels["app"]

				spec: ports: [
					for c in v.spec.template.spec.containers
					for p in c.ports {
						name:       *p.name | string
						port:       *p.containerPort | int
						targetPort: *p.containerPort | int
					},
				]
			}
		}

		ingress: [Name = _]: networking_v1.#Ingress & {
			metadata: name:      Name
			metadata: namespace: "\(#namespace)"

			metadata: labels: app: Name
		}
	}
}
