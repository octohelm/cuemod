// DO NOT EDIT THIS FILE DIRECTLY.
// generated by jsonnet extractor.
package src

import (
	grafana_example "github.com/grafana/jsonnet-libs/grafana/example:example"
	grafana "github.com/grafana/jsonnet-libs/grafana:grafana"
)

_files
_files: "dashboard.jsonnet": {
	imports: {
		"grafana/example/dashboard-simple.libsonnet": grafana_example["dashboard-simple.libsonnet"], "grafana/grafana.libsonnet": grafana["grafana.libsonnet"]
	}
	data: '''
		local grafana = import 'grafana/grafana.libsonnet';
		grafana.addDashboard('simple', (import 'grafana/example/dashboard-simple.libsonnet'), folder='Example')
		'''
}
