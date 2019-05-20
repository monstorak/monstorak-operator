# Monitoring operator for storages

This operator takes care of deployment of mixins of storage system in
kubernetes. The mixins are the mechanism for adding alerting rules,
grafana dashboards and recording rules for prometheus in kubernetes.

# Deploy Operator with mixins
To deploy the operator simply deploy to your cluster. The set of steps
to be followed for deployment of monstorak operator and related artifacts,
follow the below steps. These steps involve creation of namespace, cluster
role, role binding and deployment of operator, its custom resources and
finally the storage alerting rules.

```
oc create -f deploy/kubernetes/01_namespace.yaml
oc create -f deploy/kubernetes/02_cluster-role.yaml
oc create -f deploy/kubernetes/03_cluster-rolebinding.yaml
oc create -f deploy/kubernetes/04_operator.yaml
oc create -f deploy/kubernetes/05_storageAlerts.yaml
oc create -f jsonnet/manifests/ceph-prometheus-rules.yaml
```

The above steps consider sample deployment of ceph alerting rules.
