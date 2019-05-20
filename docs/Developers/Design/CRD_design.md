This document describes the set of [Custom Resource Definitions
(CRDs)](https://kubernetes.io/docs/concepts/extend-kubernetes/api-extension/custom-resources/)
that will be used to deploy mixins for a storage system. The actual implementation
of the resources described here will be phased in during development. The
purpose of this document is to provide the overall structure, ensuring the end
result provides necessary configurability in a user-friendly manner.

# Overview

A single Monstorak operator can take care of deployment of mixins for one or
more storage systems. Each storage cluster should already be deployed within
the kubernetes cluster, or the storage cluster could be running on hosts outside
of kubernetes cluster. The capabilities of the operator remains same in both
the modes of the deployment of storage cluster and the same set of CRDs should
be used for both the cases.

A given deployment configuration for storage systems contains a Custom Resource
which contains configurations options like whether alerts are enabled, if
dashboards are enabled, name of the storage system and version details of storage
system.

![Hierarchy of monstorak custom resources](crd_hierarchy.dot.png)

# Custom Resource Definitions

This section describes the fields in each of the custom resources.
All CRs live within the `mixins.monstorak.org` group and have version
`v1alpha1`.

## **StorageAlert**

### *Singular: StorageAlert, Plural:StorageAlerts*

The StorageAlert custom resource defines a storage system for which mixins could
be deployed using the operator. This CR captures the names of the storage system,
its supported version details, namespace in which alerts should be deployed.

```YAML
apiVersion: mixins.monstorak.org/v1alpha1
kind: StorageAlert
spec:
  storageType: "ceph"
  storageVersion: ["3.12", "3.13"]
  storageAlert:
  - labelSelector: 
    key1: value1
    key2: value2
  - prometheusNamespace: "monitoring"
```

## **StorageMonitoringPlan**

### *Singular: StorageMonitoringPlan, Plural: StorageMonitoringPlans*

The StorageMonitoringPlan custom resource defines a plan which should be
executed by operator controller to deploy the required mixins artifacts.

```YAML
apiVersion: mixins.monstorak.org/v1alpha1
kind: StorageMonitoringPlan
spec:
  - storage: "ceph"
  - version: "mimic"
    - alert: false
    - dashbaord: false
  - version: "nautilus"
    - alert: true
    - dashboard: true
```
