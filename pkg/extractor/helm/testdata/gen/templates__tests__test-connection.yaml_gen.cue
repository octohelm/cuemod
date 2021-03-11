// DO NOT EDIT THIS FILE DIRECTLY.
// generated by helm extractor.
package chart

_files: {
	"templates/tests/test-connection.yaml": '''
		apiVersion: v1
		kind: Pod
		metadata:
		  name: "{{ include "test.fullname" . }}-test-connection"
		  labels:
		    {{- include "test.labels" . | nindent 4 }}
		  annotations:
		    "helm.sh/hook": test
		spec:
		  containers:
		    - name: wget
		      image: busybox
		      command: ['wget']
		      args: ['{{ include "test.fullname" . }}:{{ .Values.service.port }}']
		  restartPolicy: Never
		'''
}
