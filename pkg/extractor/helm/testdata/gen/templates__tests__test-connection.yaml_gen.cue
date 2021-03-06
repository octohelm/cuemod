// DO NOT EDIT THIS FILE DIRECTLY.
// generated by helm extractor.
package chart

_files: {
	"templates/tests/test-connection.yaml": 'apiVersion: v1\nkind: Pod\nmetadata:\n  name: "{{ include "test.fullname" . }}-test-connection"\n  labels:\n    {{- include "test.labels" . | nindent 4 }}\n  annotations:\n    "helm.sh/hook": test\nspec:\n  containers:\n    - name: wget\n      image: busybox\n      command: [\'wget\']\n      args: [\'{{ include "test.fullname" . }}:{{ .Values.service.port }}\']\n  restartPolicy: Never\n'
}
