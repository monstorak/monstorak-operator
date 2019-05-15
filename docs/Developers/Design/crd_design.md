# Custom Resource Definitions

## **StorageCatalog**

### *Singular: StorageCatalog, Plural: StorageCatalogs*

```YAML
kind: StorageCatalog
spec:
  - storage: "ceph"
    version: ["mimic", "nautilus"]
    alert: true
    dashboard: true
```

## **StorageAlert**

### *Singular: StorageAlert, Plural: StorageAlerts*

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

## **StorageDashboard**

### *Singular: StorageDashboard, Plural: StorageDashboards*

```YAML
kind: StorageDashboard
spec:
  - storage: "ceph"
    - version: "mimic"
      dashboard: false
    - version: "nautilus"
      dashboard: true
      grafana:
      - namespace: "monitoring"
      - label: "grafana:k8s"
```

## **StorageMonitoringPlan**

### *Singular: StorageMonitoringPlan, Plural: StorageMonitoringPlans*

```YAML
kind: StorageMonitoringPlan
spec:
  TBD
```
