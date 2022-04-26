// DO NOT EDIT THIS FILE DIRECTLY.
// generated by helm extractor.
package chart

_files: "templates/_helpers.tpl": '{{/*\nExpand the name of the chart.\n*/}}\n{{- define "test.name" -}}\n{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}\n{{- end }}\n\n{{/*\nCreate a default fully qualified app name.\nWe truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).\nIf release name contains chart name it will be used as a full name.\n*/}}\n{{- define "test.fullname" -}}\n{{- if .Values.fullnameOverride }}\n{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" }}\n{{- else }}\n{{- $name := default .Chart.Name .Values.nameOverride }}\n{{- if contains $name .Release.Name }}\n{{- .Release.Name | trunc 63 | trimSuffix "-" }}\n{{- else }}\n{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" }}\n{{- end }}\n{{- end }}\n{{- end }}\n\n{{/*\nCreate chart name and version as used by the chart label.\n*/}}\n{{- define "test.chart" -}}\n{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}\n{{- end }}\n\n{{/*\nCommon labels\n*/}}\n{{- define "test.labels" -}}\nhelm.sh/chart: {{ include "test.chart" . }}\n{{ include "test.selectorLabels" . }}\n{{- if .Chart.AppVersion }}\napp.kubernetes.io/version: {{ .Chart.AppVersion | quote }}\n{{- end }}\napp.kubernetes.io/managed-by: {{ .Release.Service }}\n{{- end }}\n\n{{/*\nSelector labels\n*/}}\n{{- define "test.selectorLabels" -}}\napp.kubernetes.io/name: {{ include "test.name" . }}\napp.kubernetes.io/instance: {{ .Release.Name }}\n{{- end }}\n\n{{/*\nCreate the name of the service account to use\n*/}}\n{{- define "test.serviceAccountName" -}}\n{{- if .Values.serviceAccount.create }}\n{{- default (include "test.fullname" .) .Values.serviceAccount.name }}\n{{- else }}\n{{- default "default" .Values.serviceAccount.name }}\n{{- end }}\n{{- end }}\n'
