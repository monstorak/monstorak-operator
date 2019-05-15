# Monitoring operator for storages

This operator takes care of deployment of mixins of storage system in Openshift and Kubernetes. The mixins are the mechanism for adding alerting rules, grafana dashboards and recording rules for prometheus in Openshift and Kubernetes.

# Deploy Operator with mixins
To deploy the operator simply deploy to your cluster. The set of steps to be followed
for deployment of monstorak operator and related artifacts, follow the below steps

## Create the namespace `storage-monitoring`
```
oc create -f deploy/kubernetes/01_namespace.yaml
```

## Create required cluster roles
```
oc create -f deploy/kubernetes/02_cluster-role.yaml
```

## Create cluster role binding
```
oc create -f deploy/kubernetes/03_cluster-rolebinding.yaml
```

## Deploy the operator
```
oc create -f deploy/kubernetes/04_operator.yaml
```

## Deploy the custom resources
```
oc create -f deploy/kubernets/05_storageAlerts.yaml
```

## Deploy the storage alert rules
```
oc create -f jsonnet/manifests/ceph-prometheus-rules.yaml
```
