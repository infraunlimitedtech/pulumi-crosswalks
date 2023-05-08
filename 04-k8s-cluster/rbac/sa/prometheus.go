package sa

import (
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/meta/v1"
	rbacv1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/rbac/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"

	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/core/v1"
)

func PrometheusAccount(ns *corev1.Namespace, ctx *pulumi.Context) error {
	serviceName := "infra-prometheus"

	serviceAccount, err := corev1.NewServiceAccount(ctx, "promServiceAccount", &corev1.ServiceAccountArgs{
		ApiVersion: pulumi.String("v1"),
		Kind:       pulumi.String("ServiceAccount"),
		Metadata: &metav1.ObjectMetaArgs{
			Name:      pulumi.String(serviceName),
			Namespace: ns.Metadata.Name(),
		},
	})
	if err != nil {
		return err
	}

	clusterRole, err := rbacv1.NewClusterRole(ctx, "promClusterRole", &rbacv1.ClusterRoleArgs{
		ApiVersion: pulumi.String("rbac.authorization.k8s.io/v1"),
		Kind:       pulumi.String("ClusterRole"),
		Metadata: &metav1.ObjectMetaArgs{
			Name: serviceAccount.Metadata.Name(),
		},
		Rules: rbacv1.PolicyRuleArray{
			&rbacv1.PolicyRuleArgs{
				ApiGroups: pulumi.StringArray{
					pulumi.String(""),
				},
				Resources: pulumi.StringArray{
					pulumi.String("nodes"),
					pulumi.String("nodes/metrics"),
				},
				Verbs: pulumi.StringArray{
					pulumi.String("list"),
					pulumi.String("get"),
					pulumi.String("watch"),
				},
			},
			&rbacv1.PolicyRuleArgs{
				ApiGroups: pulumi.StringArray{
					pulumi.String(""),
				},
				Resources: pulumi.StringArray{
					pulumi.String("endpoints"),
					pulumi.String("pods"),
					pulumi.String("services"),
				},
				Verbs: pulumi.StringArray{
					pulumi.String("list"),
					pulumi.String("get"),
					pulumi.String("watch"),
				},
			},
		},
	})
	if err != nil {
		return err
	}

	_, err = rbacv1.NewClusterRoleBinding(ctx, "promClusterRoleBinding", &rbacv1.ClusterRoleBindingArgs{
		ApiVersion: pulumi.String("rbac.authorization.k8s.io/v1"),
		Kind:       pulumi.String("ClusterRoleBinding"),
		Metadata: &metav1.ObjectMetaArgs{
			Name:      serviceAccount.Metadata.Name().Elem(),
			Namespace: serviceAccount.Metadata.Namespace().Elem(),
		},
		RoleRef: &rbacv1.RoleRefArgs{
			ApiGroup: pulumi.String("rbac.authorization.k8s.io"),
			Kind:     pulumi.String("ClusterRole"),
			Name:     clusterRole.Metadata.Name().Elem(),
		},
		Subjects: rbacv1.SubjectArray{
			&rbacv1.SubjectArgs{
				Kind:      pulumi.String("ServiceAccount"),
				Name:      serviceAccount.Metadata.Name().Elem(),
				Namespace: serviceAccount.Metadata.Namespace().Elem(),
			},
		},
	})
	if err != nil {
		return err
	}

	_, err = corev1.NewSecret(ctx, "promTokenSecret", &corev1.SecretArgs{
		ApiVersion: pulumi.String("v1"),
		Kind:       pulumi.String("Secret"),
		Metadata: &metav1.ObjectMetaArgs{
			Name:      serviceAccount.Metadata.Name().Elem(),
			Namespace: serviceAccount.Metadata.Namespace().Elem(),
			Annotations: pulumi.StringMap{
				"kubernetes.io/service-account.name": serviceAccount.Metadata.Name().Elem(),
			},
		},
		Type: pulumi.String("kubernetes.io/service-account-token"),
	})
	if err != nil {
		return err
	}

	return nil
}
