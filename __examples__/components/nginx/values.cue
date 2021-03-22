package nginx

#values: image: {
	repository: *"nginx" | string
	tag:        *"alpine" | string
	pullPolicy: *"IfNotPresent" | string
}
