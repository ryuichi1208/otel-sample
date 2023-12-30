apiVersion: "v1"
kind:       "Service"
metadata: {
	name:      "nats-js"
	namespace: "nats-ns"
	labels: {
		"app.kubernetes.io/name":     "nats"
		"app.kubernetes.io/instance": "nats-js"
		"app.kubernetes.io/version":  "2.6.2"
	}
}
spec: {
	selector: {
		"app.kubernetes.io/name":     "nats"
		"app.kubernetes.io/instance": "nats-js"
	}
	clusterIP: "None"
	ports: [{
		name:        "client"
		port:        4222
		appProtocol: "tcp"
	}, {
		name:        "cluster"
		port:        6222
		appProtocol: "tcp"
	}, {
		name:        "monitor"
		port:        8222
		appProtocol: "tcp"
	}, {
		name:        "metrics"
		port:        7777
		appProtocol: "tcp"
	}]
}

