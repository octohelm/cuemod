local d = import 'doc-util/main.libsonnet';
local grafana = import 'grafana/grafana.libsonnet';
local k = import 'k.libsonnet';
{
  config+:: {
    prometheus_url: 'http://prometheus',
  },
  '#withName':: d.fn(help='`name` is the name of the service. Required', args=[d.arg(name='name', type=d.T.string)]),
  withName(name): { name: name },

  namespace: k.core.v1.namespace.new(importstr './sub/name.txt'),

  prometheus_datasource:: grafana.datasource.new('prometheus', $.config.prometheus_url, type='prometheus', default=true),

  grafana: grafana
           + grafana.withAnonymous()
           + grafana.addFolder('Example')
           + (import 'dashboard.jsonnet')
           + grafana.addDatasource('prometheus', $.prometheus_datasource)
           ,
}