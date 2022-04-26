package release

import (
	"encoding/hex"
	"encoding/json"
	"crypto/sha256"

	core_v1 "k8s.io/api/core/v1"
	apps_v1 "k8s.io/api/apps/v1"
)

_workload: {
	#namespace:             string
	#serviceSelectorLabels: *["app"] | [...string]

	spec: {
		configMaps: [Name = _]: core_v1.#ConfigMap & {
			metadata: name:      Name
			metadata: namespace: #namespace
		}

		secrets: [Name = _]: core_v1.#Secret & {
			metadata: name:      Name
			metadata: namespace: #namespace
			type: *"Opaque" | string
		}

		_template: {
			#volumes: [string]: {
				core_v1.#VolumeMount

				name:   string
				volume: core_v1.#Volume
			}

			_container: {
				#name: string

				core_v1.#Container & {
					name: #name

					#ports: [string]:   int
					#envVars: [string]: string | core_v1.#EnvVarSource

					ports: [
						for n, cp in #ports {
							{
								name:          n
								containerPort: cp
							}
						},
					]

					env: [
						for n, v in #envVars {
							let isString = (v & string) != _|_

							if isString {
								{
									name:  n
									value: v
								}
							}

							if !isString {
								{
									name:      n
									valueFrom: v
								}
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
			}

			#initContainers: [ContainerName = _]: _container & {
				#name: ContainerName
			}

			#containers: [ContainerName = _]: _container & {
				#name: ContainerName
			}

			spec: template: spec: volumes: [
				for n, v in #volumes {
					{
						name: n
						v.volume
					}
				},
			]

			spec: template: spec: initContainers: [
				for c in #initContainers {c},
			]

			spec: template: spec: containers: [
				for c in #containers {c},
			]

			// auto added checksum fo secrets & configMaps
			for n, v in #volumes {
				for serectName, serect in secrets {
					for volumeSourceName, volumeSource in v.volume {
						if volumeSourceName == "serect" {
							if volumeSource.name == serectName {
								spec: template: metadata: annotations: "checksum/\(n)": "\(hex.Encode(sha256.Sum256(json.Marshal(serect.data))))"
							}
						}
					}
				}
				for configMapName, cm in configMaps {
					for volumeSourceName, volumeSource in v.volume {
						if volumeSourceName == "configMap" {
							if volumeSource.name == configMapName {
								spec: template: metadata: annotations: "checksum/\(n)": "\(hex.Encode(sha256.Sum256(json.Marshal(cm.data))))"
							}
						}
					}
				}
			}
		}

		deployments: [Name = _]: apps_v1.#Deployment & {
			metadata: name:      Name
			metadata: namespace: #namespace

			metadata: labels: app: Name

			spec: template: metadata: labels: app: Name
			spec: selector: matchLabels: app: Name

			spec: replicas: *1 | int

			_template
		}

		daemonSets: [Name = _]: apps_v1.#DaemonSet & {
			metadata: name:      Name
			metadata: namespace: #namespace

			metadata: labels: app: Name

			spec: template: metadata: labels: app: Name
			spec: selector: matchLabels: app: Name

			_template
		}

		statefulSets: [Name = _]: apps_v1.#StatefulSet & {
			metadata: name:      Name
			metadata: namespace: #namespace

			metadata: labels: app: Name

			spec: template: metadata: labels: app: Name
			spec: selector: matchLabels: app: Name

			_template
		}

		for x in [deployments, daemonSets, statefulSets]
		for n, w in x {
			if len(w.spec.template.spec.containers) > 0 && len(w.spec.template.spec.containers[0].ports) > 0 {
				services: "\(n)": core_v1.#Service & {
					spec: selector: {
						for label in #serviceSelectorLabels {
							"\(label)": w.spec.template.metadata.labels[label]
						}
					}

					spec: ports: [
						for c in w.spec.template.spec.containers
						for p in c.ports {
							name:       *p.name | string
							port:       *p.containerPort | int
							targetPort: *p.containerPort | int
						},
					]
				}
			}
		}
	}
}
