# deploying
egg is deployed into the egg namespace by default

you need to install clickhouse yourself, i recommend the clickhouse operator: https://github.com/Altinity/clickhouse-operator

## steps
1. deploy clickhouse
2. change the CLICKHOUSE_URI for ingress and egrees in deployment.yaml
3. run `kubectl apply -k .` to apply all the manifests (using kustomize)