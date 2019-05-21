# Operator for storage monitoring mixin

This operator takes care of deployment of mixins of storage system in
kubernetes. The mixins is an abstraction framework that makes it easier
for developers to collaborate on common set of features like prometheus
rules and grafana dashboards and it also increases the re-usability.

# Deploy mixin with operator
To deploy the operator simply deploy to your cluster. The set of steps
to be followed for deployment of monstorak operator and related artifacts,
follow the below steps. The set of steps performed by operator involve
creation of namespace, role, role binding, its custom resources and
finally the storage alerting rules.

Create a namespace where the operator would be deployed and then run the
below commands -

```
oc create -f deploy/crds/alerts_v1alpha1_storagealert_crd.yaml
oc create -f deploy/serviceaccount.yaml
oc create -f deploy/role.yaml
oc create -f deploy/role_binding.yaml
oc create -f deploy/operator.yaml
oc create -f deploy/crds/alerts_v1alpha1_storagealert_cr.yaml
```
