package kilo

import (
	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/core/v1"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Namespace function creates ns with reasonable defaults for kilo.
func CreateNS(ctx *pulumi.Context, name string) (*corev1.Namespace, error) {
	ns, err := corev1.NewNamespace(ctx, name, &corev1.NamespaceArgs{
		Metadata: &metav1.ObjectMetaArgs{
			Name: pulumi.String(name),
		},
	})
	if err != nil {
		return nil, err
	}

	return ns, nil
}
