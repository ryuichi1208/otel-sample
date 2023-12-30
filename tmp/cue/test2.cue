apiVersion: "apps/v1"
kind:       "Deployment"
metadata: {
    name: "nginx-deployment"
    labels: app: "nginx"
}
spec: {
    replicas: 3
    selector: matchLabels: app: "nginx"
    template: {
        metadata: labels: app: "nginx"
        spec: containers: [{
            name:  "nginx"
            image: "nginx:1.14.2"
            ports: [{
                containerPort: 80
            }]
        }]
    }
}
