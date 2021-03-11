local grafana = import 'grafana/grafana.libsonnet';
grafana.addDashboard('simple', (import 'grafana/example/dashboard-simple.libsonnet'), folder='Example')
