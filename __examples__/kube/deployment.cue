package kube

import (
	apps_v1 "k8s.io/api/apps/v1"
)

deployment: [Name=_]: apps_v1.#Deployment & {
	metadata: name: Name

	_labels: app: Name

	spec: selector: matchLabels: _labels
	spec: template: metadata: labels: _labels
}

deployment: nginx: spec: template: spec: {
	volumes: [
		{
			name: "secret-volume"
			secret: secretName: "proxy-secrets"
		},
		{
			name: "config-volume"
			configMap: name: "nginx"
		},
	]
	containers: [
		{
			name:  "nginx"
			image: "nginx:1.11.10-alpine"
			ports: [
				{
					containerPort: 80
				},
				{
					containerPort: 443
				},
			]
			volumeMounts: [
				{
					mountPath: "/etc/ssl"
					name:      "secret-volume"
				},
				{
					name:      "config-volume"
					mountPath: "/etc/nginx/nginx.conf"
					subPath:   "nginx.conf"
				},
			]
		}]
}
