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

A given deployment configuration for storage systems contains a Custome Resource
which contains configurations options like whether alerts are enabled, if
dashboards are enabled, name of the storage system and version details of storage
system.

# Custom resources

This section describes the fields in each of the custom resources.

## StorageMixinClass CR

The storage mixin class CR defines the mixins deployment specific configuration
for a storage system. A commented example is shown below:

```yaml
apiVersion: mixins.monstorak.org/v1alpha1
kind: StorageMixinClass
metadata:
  name: example-storage-mixin-class
spec:
  # Add fields here
  size: 3
  storageName: "name-of-the-storage-system e.g. ceph"
  storageVersion: "version-of-storage-system"
  storageAlert:
    enabled: true
    prometheusNamespace: "namespace-in-which-cr-to-be-deployed"
  storageDashboard:
    enabled: false
    grafanaNamespace: "namespace-in-which-grafana-should-be-deployed"
status:
  # TBD operator state (e.g. running, completed, pending etc)
  ...
```

All CRs live within the `mixins.monstorak.org` group and have version
`v1alpha1`. The storage cluster and all its artifacts are suposed to be
deployed in single namespace.  The `spec` field provides the main
configuration options.

The `storageAlert` section mentions if alerts are enabled for the said
storage system and operator needs to deploy the same.

The `storageDashboard` section mentions if grafan dashboard are enabled
for the storage system and operator needs to deploy the same.

The `storageName` field holds the name of the storage system e.g. ceph/gluster

The `storageVersion` field holds the version details of the storage system.
