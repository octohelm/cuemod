package nginx

import (
	"github.com/x/examples/pkg/release"
)

release.#Release & {
	#name:      "web"
	#namespace: "default"

	spec: configmap: "\(#name)-html": data: "index.html": _indexHTML

	spec: deployment: "\(#name)": {
		#volumes: html: {
			mountPath: "/usr/share/nginx/html"
			volume: configMap: name: "\(#name)-html"
		}

		#container: {
			name:            "\(#name)"
			image:           "\(#values.image.repository):\(#values.image.tag)"
			imagePullPolicy: "\(#values.image.pullPolicy)"

			#ports: {
				http:  80
				https: 443
			}
		}
	}
}
