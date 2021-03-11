package jsonnet_demo

import (
	"encoding/json"
	"extension/jsonnet"

	"github.com/grafana/jsonnet-libs/grafana"
)

"grafana": jsonnet.Eval({
	data: '''
		local grafana = import 'github.com/grafana/jsonnet-libs/grafana/grafana.libsonnet';
		
		{
		    config+:: (import 'config.jsonnet'),
		
		    prometheus_datasource:: grafana.datasource.new('prometheus', $.config.prometheus_url, type='prometheus', default=true),
		
		    grafana: grafana
		         + grafana.withAnonymous()
		         + grafana.addFolder('Example')
		         + grafana.addDatasource('prometheus', $.prometheus_datasource)
		         ,
		}
		'''

	imports: "github.com/grafana/jsonnet-libs/grafana/grafana.libsonnet": grafana["grafana.libsonnet"]
	imports: "config.jsonnet": code: json.Marshal({prometheus_url: 'http://prometheus'})
})
