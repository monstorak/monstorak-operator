package common

import (
	rbac "k8s.io/api/rbac/v1"
	errors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/client-go/kubernetes/typed/rbac/v1"
	clientcmd "k8s.io/client-go/tools/clientcmd"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
)

var rbacLog = logf.Log.WithName("common_rbac")

func newRbacV1Client() (*v1.RbacV1Client, error) {
	client, err := clientcmd.BuildConfigFromFlags("", "")
	if err != nil {
		return nil, err
	}
	return v1.NewForConfig(client)
}

func CreateOrUpdateRBAC() error {
	rbacLog.Info("Creating Role/RoleBinding")
	prometheusK8sRole := newRole(
		"storage-prometheus",
		"openshift-storage",
		newPolicyRules(
			newPolicyRule(
				[]string{""},
				[]string{"pods", "services", "endpoints"},
				[]string{},
				[]string{"list", "watch"},
				[]string{},
			),
		),
	)

	err := createOrUpdateRole(prometheusK8sRole)
	if err != nil {
		return err
	}

	subject := newSubject("ServiceAccount", "prometheus-k8s", "openshift-monitoring")
	subject.APIGroup = ""

	prometheusK8sRolebinding := newRoleBinding(
		"storage-prometheus",
		"storage-prometheus",
		"openshift-storage",
		newSubjects(
			subject,
		),
	)
	return createOrUpdateRoleBinding(prometheusK8sRolebinding)
}

func newRole(roleName string, namespace string, rules []rbac.PolicyRule) *rbac.Role {
	return &rbac.Role{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Role",
			APIVersion: rbac.SchemeGroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      roleName,
			Namespace: namespace,
		},
		Rules: rules,
	}
}

func newClusterRole(roleName string, rules []rbac.PolicyRule) *rbac.ClusterRole {
	return &rbac.ClusterRole{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ClusterRole",
			APIVersion: rbac.SchemeGroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: roleName,
		},
		Rules: rules,
	}
}

func newRoleBinding(bindingName, roleName string, namespace string, subjects []rbac.Subject) *rbac.RoleBinding {
	return &rbac.RoleBinding{
		TypeMeta: metav1.TypeMeta{
			Kind:       "RoleBinding",
			APIVersion: rbac.SchemeGroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      bindingName,
			Namespace: namespace,
		},
		RoleRef: rbac.RoleRef{
			Kind:     "Role",
			Name:     roleName,
			APIGroup: rbac.GroupName,
		},
		Subjects: subjects,
	}
}

func newClusterRoleBinding(bindingName, roleName string, subjects []rbac.Subject) *rbac.ClusterRoleBinding {
	return &rbac.ClusterRoleBinding{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ClusterRoleBinding",
			APIVersion: rbac.SchemeGroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: bindingName,
		},
		RoleRef: rbac.RoleRef{
			Kind:     "ClusterRole",
			Name:     roleName,
			APIGroup: rbac.GroupName,
		},
		Subjects: subjects,
	}
}

func newSubject(kind, name, namespace string) rbac.Subject {
	return rbac.Subject{
		Kind:      kind,
		Name:      name,
		Namespace: namespace,
		APIGroup:  rbac.GroupName,
	}
}

func newSubjects(subjects ...rbac.Subject) []rbac.Subject {
	return subjects
}

func newPolicyRule(apiGroups, resources, resourceNames, verbs, urls []string) rbac.PolicyRule {
	return rbac.PolicyRule{
		APIGroups:       apiGroups,
		Resources:       resources,
		ResourceNames:   resourceNames,
		Verbs:           verbs,
		NonResourceURLs: urls,
	}
}

func newPolicyRules(rules ...rbac.PolicyRule) []rbac.PolicyRule {
	return rules
}

func createOrUpdateRole(role *rbac.Role) error {
	rbacV1Client, err := newRbacV1Client()
	if err != nil {
		rbacLog.Error(err, "Failed to create rbacV1Client")
		return err
	}

	roles := rbacV1Client.Roles(role.Namespace)
	_, err = roles.Create(role)

	if !errors.IsAlreadyExists(err) {
		rbacLog.Error(err, "Failed to create Role", "role: ", role)
		return err
	}

	_, err = roles.Update(role)
	if err != nil {
		rbacLog.Error(err, "Failed to update Role", "role: ", role)
	}
	return err
}

func createOrUpdateRoleBinding(roleBinding *rbac.RoleBinding) error {
	rbacV1Client, err := newRbacV1Client()
	if err != nil {
		rbacLog.Error(err, "Failed to create rbacV1Client")
		return err
	}

	roleBindings := rbacV1Client.RoleBindings(roleBinding.Namespace)
	_, err = roleBindings.Create(roleBinding)

	if !errors.IsAlreadyExists(err) {
		rbacLog.Error(err, "Failed to create RoleBinding", "roleBinding: ", roleBinding)
		return err
	}

	_, err = roleBindings.Update(roleBinding)
	if err != nil {
		rbacLog.Error(err, "Failed to update RoleBinding", "roleBinding: ", roleBinding)
	}
	return err
}
