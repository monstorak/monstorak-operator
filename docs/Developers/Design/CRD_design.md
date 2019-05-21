This document describes a set of [Custom Resource Definitions
(CRDs)](https://kubernetes.io/docs/concepts/extend-kubernetes/api-extension/custom-resources/)
that will be used to deploy mixin and ensure user friendly configuration options
for monitoring storage systems.

# Overview

An instance of monstorak operator takes care of deployment of mixins for one or
more storage systems. It is expected that -
1. deployment of the storage cluster to be monitored is taken care by the user
2. prometheus can scrape metrics from this cluster

![Hierarchy of monstorak custom resources](crd_hierarchy.dot.png)

# Custom Resource Definitions

This section describes the fields in each of the custom resources.
All CRs live within the `alerts.monstorak.org` group and have version
`v1alpha1`.

## **StorageAlert**

### *Singular: StorageAlert, Plural:StorageAlerts*

This CR defines configuration that is used by the operator to deploy monitoring
mixin for the storage systems.

```YAML
kind: StorageAlert
spec:
  - storage: "ceph"
    - version: "mimic"
      alert: false
    - version: "nautilus"
      alert: true
      prometheus:
      - namespace: "monitoring"
      - label: "prometheus:k8s"
```
