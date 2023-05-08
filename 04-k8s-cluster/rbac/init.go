package rbac

import (
	"k8s-cluster/rbac/sa"

	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/meta/v1"
	rbacv1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/rbac/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

type RBAC struct {
	ServiceAccounts *sa.ServiceAccounts
}

func Init(ctx *pulumi.Context) (*RBAC, error) {
	var pulumiRBACCfg *RBAC
	cfg := config.New(ctx, "")
	cfg.RequireSecretObject("rbac", &pulumiRBACCfg)

	_, err := rbacv1.NewClusterRoleBinding(ctx, "rbacInfraAdmins", &rbacv1.ClusterRoleBindingArgs{
		ApiVersion: pulumi.String("rbac.authorization.k8s.io/v1"),
		Kind:       pulumi.String("ClusterRoleBinding"),
		Metadata: &metav1.ObjectMetaArgs{
			Name: pulumi.String("infra:admins"),
		},
		Subjects: rbacv1.SubjectArray{
			&rbacv1.SubjectArgs{
				Kind:     pulumi.String("Group"),
				Name:     pulumi.String("infra:admins"),
				ApiGroup: pulumi.String("rbac.authorization.k8s.io"),
			},
		},
		RoleRef: &rbacv1.RoleRefArgs{
			Kind:     pulumi.String("ClusterRole"),
			Name:     pulumi.String("cluster-admin"),
			ApiGroup: pulumi.String("rbac.authorization.k8s.io"),
		},
	}, pulumi.DeleteBeforeReplace(true))
	if err != nil {
		return nil, err
	}

	return pulumiRBACCfg, err
}
