local ceph = (import 'ceph-mixins/mixin.libsonnet');
local noobaa = (import 'noobaa-mixins/mixin.libsonnet');
local gluster = (import 'gluster-mixins/mixin.libsonnet');
local prometheusRule(storageNamespace, storageName, alertType, prometheusLabel) = {
  apiVersion: 'monitoring.coreos.com/v1',
  kind: 'PrometheusRule',
  metadata: {
    labels: {
      prometheus: prometheusLabel,
      role: 'alert-rules',
    },
    name: 'prometheus-' + storageName + '-rules',
    namespace: storageNamespace,
  },
  spec: alertType,
};

{
  'ceph-prometheus-rules': prometheusRule('default', 'ceph', ceph.prometheusAlerts, 'k8s'),
  'noobaa-prometheus-rules': prometheusRule('default', 'noobaa', noobaa.prometheusAlerts, 'k8s'),
  'gluster-prometheus-rules': prometheusRule('default', 'gluster', gluster.prometheusAlerts, 'k8s'),
}
