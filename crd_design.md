``` yaml
# YAML

apiVersion: "KUBERNETES_API_VERSION_USED"
kind: CRD
metadata:
  # Name for this node
  name: storagealerts.alerts.monstorak.org
  # CRD is namespaced
  namespace: monstorak
  annotations:
    # Applied by operator when it creates/manages this object from a template.
    # When this is present, contents will be dynamically adjusted accd to the
    # template in the cluster CR.
    # When this annotation is present, the admin may only modify
    # .spec.desiredState or delete the CR. Any other change will be
    # overwritten.
    # operator.monstorak.org/template: template-name
  version: VERSION_STRING
spec:
  # Nodes belong to a cluster
  # not needed for now, will work on multi-cluster later
  cluster: CLUSTER_NAME
  # Only 1 of external | storage
  storage:
    # Initially only ceph will be supported
    # and later extend it to others
    - storageType: [ceph | gluster | cassandra | ...]
      storageVersion: STORAGE_VERSION_USED
      storageAlert:
      - enabled: [yes | no]
        labelSelctor: MAP_OF_KEY_VALUE_PAIRS
        prometheusNamespace: OPTIONAL_STR_DEFINES_PROMETHEUS_NAMESPACE
      storageDashboard: 
      - enabled: [yes | no]
        labelSelctor: MAP_OF_KEY_VALUE_PAIRS
        grafanaNamespace: OPTIONAL_STR_DEFINES_GRAFANA_NAMESPACE
    ...
  nodeAffinity:
    # https://kubernetes.io/docs/concepts/configuration/assign-pod-node/#affinity-and-anti-affinity
    ...
status:
  # TBD operator state
  # Possible states: (enabled | deleting | disabled)
  # see other 'currentState' keys in other operators
  currentState: enabled
```
