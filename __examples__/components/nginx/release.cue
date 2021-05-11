package nginx

import (
	"github.com/octohelm/cuem/release"
)

release.#Release & {
	#name:      "web"
	#namespace: "default"

	spec: {
		configMaps: "\(#name)-html": data: "index.html": _indexHTML

		deployments: "\(#name)": {
			#volumes: html: {
				mountPath: "/usr/share/nginx/html"
				volume: configMap: name: "\(#name)-html"
			}

			#containers: web: {
				image:           "\(#values.image.repository):\(#values.image.tag)"
				imagePullPolicy: "\(#values.image.pullPolicy)"

				#ports: {
					http:  80
					https: 443
				}
			}
		}
	}
}
