package demo

import (
	"github.com/x/examples/kube"
)

apiVersion: "tanka.dev/v1alpha1"
kind:       "Environment"
metadata: name:     "demo"
spec: apiServer:    "https://172.16.0.7:8443"
spec: injectLabels: true
spec: namespace:    "default"
data: kube
