This document describes the actions that have to be executed successfully for
monstorak's CRDs to be considered reconciled.

![ActionGraph](workflow.dot.png)

# Pre-requisites

## Cluster monitoring deployed
This operator expects that cluster monitoring stack is already deployed in
k8s and all required services like prometheus and its related artifacts
are already available.

## Storage cluster deployed
The ceph mixins (alerting rules) are meant for storage, so the operator expects
the storage stack to be deployed and available already. The operator
![Rook](https://github.com/rook/rook/) takes care of deployment of ceph cluster.

## Monitoring Enabled For Storage
The namespace in which storage is deployed, monitoring needs to be enabled for
the same already. Below commands can be used to enable the required configurations

```
- oc label namespace <storage-namespoace> openshift.io/cluster-monitoring=true
- oc policy add-role-to-user view system:serviceaccount:openshift-monitoring:prometheus-k8s -n <storage-namespace>
```

# Deploy storage mixin
Mixins are mechanism to write alerting rules and grafana dashboard templates in
kubernetes. The project ![ceph-mixins](https://github.com/ceph/ceph-mixins/)
provides the alerting rules for ceph storage system. The mixins use a templating
language called jsonnet for writting the alerting rules and grafana dashboard
templates, which ultimately needs to be compiled into prometheus rules YAML and
grafana dashboard JSON templates.

This operator take care of deploying the prometheus alerting rules generated out
of ![ceph-mixins](https://github.com/ceph/ceph-mixins/) in a kubernetes cluster.

The monstorak operator consists of two controllers. The first one takes care of
the instanciating/reconciling custom resources like StorageCatalog,
StorageAlert and StorageDasboard. Once custom resources are available, the
second controller forms storage monitoring plans and executes them one by one
based on how many storage system and what artificats to be deployed.
