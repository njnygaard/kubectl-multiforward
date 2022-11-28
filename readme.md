# Kubectl Multiforward

Forward to multiple services in your current Kubernetes context based on a config file.

## Example

```bash
kubectl multiforward staging
```

```plaintext
┌───────────────┬────────────────────────┐
│ Alertmanager  │ http://localhost:9093  │
│ Prometheus    │ http://localhost:9090  │
│ Grafana       │ http://localhost:3000  │
│ Kibana        │ http://localhost:5601  │
│ Elasticsearch │ https://localhost:9200 │
└───────────────┴────────────────────────┘
Monitoring Resources... ^C to exit
```
## Configuration
 - localAddress : (Optional) if you want to discover you local IP, please use 'auto'. You can set anything you want inside
 - localPort : the local machine port on which the service will be published 
 - servicePort : the service exposed port 
 
