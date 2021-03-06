// DO NOT EDIT THIS FILE DIRECTLY.
// generated by helm extractor.
package chart

_files: {
	"templates/service.yaml": 'apiVersion: v1\nkind: Service\nmetadata:\n  name: {{ include "test.fullname" . }}\n  labels:\n    {{- include "test.labels" . | nindent 4 }}\nspec:\n  type: {{ .Values.service.type }}\n  ports:\n    - port: {{ .Values.service.port }}\n      targetPort: http\n      protocol: TCP\n      name: http\n  selector:\n    {{- include "test.selectorLabels" . | nindent 4 }}\n'
}
