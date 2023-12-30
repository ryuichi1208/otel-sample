deployments: [_name=string]: {
    _image: string
    _port:  int

    apiVersion: "apps/v1"
    kind:       "Deployment"
    metadata: {
        name: _name + "-deployment"
        labels: app: _name
    }
    spec: {
        replicas: 3
        selector: matchLabels: app: _name
        template: {
            metadata: labels: app: _name
            spec:
	        securityContext: {
                    runAsUser:  1000
                    runAsGroup: 3000
                    fsGroup:    2000
                }
	        containers: [{
                    name:  _name
                    image: _image
                    ports: [{
                        containerPort: _port
                    }]
            }]
        }
    }
}

