# Operator for Storage Monitoring Mixins

Monstorak operator takes care of deployment of monitoring mixins for storage systems in **Openshift**.

Monitoring mixin is an abstraction framework that makes it easier
for developers to collaborate on common set of features like prometheus
rules and grafana dashboards. It brings modularity and increases re-usability of code.

# Deploying Monitoring Mixins with Operator

### 1. Deploy Operator

* Create a **Namespace** where the operator would be deployed,

  `oc create namespace storage-monitoring`

* Create **CRDs**

  `oc create -f deploy/crds/alerts_v1alpha1_storagealert_crd.yaml`

* Create a **Service Account**

  `oc create -f deploy/serviceaccount.yaml`

* Create **Role** and **Cluster Role**

  `oc create -f deploy/role.yaml`

* Create **Role Binding** and **Cluster Role Binding**

  `oc create -f deploy/role_binding.yaml`

* Create a **Deployment** of **Operator**

  `oc create -f deploy/operator.yaml`

### 2. Deploy Monitoring Mixins

* Create a **CR** of type **StorageAlert**

  * For Ceph storage alerts,

     `oc create -f deploy/example/ceph-storagealert.yaml`

  * For Noobaa storage alerts,

     `oc create -f deploy/example/noobaa-storagealert.yaml`

     **OR,**
  * For both,

     `oc create -f deploy/example/monstorak-storagealert.yaml`