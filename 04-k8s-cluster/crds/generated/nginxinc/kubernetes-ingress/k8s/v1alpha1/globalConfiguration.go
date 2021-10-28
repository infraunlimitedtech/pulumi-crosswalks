// *** WARNING: this file was generated by crd2pulumi. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

package v1alpha1

import (
	"context"
	"reflect"

	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// GlobalConfiguration defines the GlobalConfiguration resource.
type GlobalConfiguration struct {
	pulumi.CustomResourceState

	ApiVersion pulumi.StringPtrOutput     `pulumi:"apiVersion"`
	Kind       pulumi.StringPtrOutput     `pulumi:"kind"`
	Metadata   metav1.ObjectMetaPtrOutput `pulumi:"metadata"`
	// GlobalConfigurationSpec is the spec of the GlobalConfiguration resource.
	Spec GlobalConfigurationSpecPtrOutput `pulumi:"spec"`
}

// NewGlobalConfiguration registers a new resource with the given unique name, arguments, and options.
func NewGlobalConfiguration(ctx *pulumi.Context,
	name string, args *GlobalConfigurationArgs, opts ...pulumi.ResourceOption) (*GlobalConfiguration, error) {
	if args == nil {
		args = &GlobalConfigurationArgs{}
	}

	args.ApiVersion = pulumi.StringPtr("k8s.nginx.org/v1alpha1")
	args.Kind = pulumi.StringPtr("GlobalConfiguration")
	var resource GlobalConfiguration
	err := ctx.RegisterResource("kubernetes:k8s.nginx.org/v1alpha1:GlobalConfiguration", name, args, &resource, opts...)
	if err != nil {
		return nil, err
	}
	return &resource, nil
}

// GetGlobalConfiguration gets an existing GlobalConfiguration resource's state with the given name, ID, and optional
// state properties that are used to uniquely qualify the lookup (nil if not required).
func GetGlobalConfiguration(ctx *pulumi.Context,
	name string, id pulumi.IDInput, state *GlobalConfigurationState, opts ...pulumi.ResourceOption) (*GlobalConfiguration, error) {
	var resource GlobalConfiguration
	err := ctx.ReadResource("kubernetes:k8s.nginx.org/v1alpha1:GlobalConfiguration", name, id, state, &resource, opts...)
	if err != nil {
		return nil, err
	}
	return &resource, nil
}

// Input properties used for looking up and filtering GlobalConfiguration resources.
type globalConfigurationState struct {
	ApiVersion *string            `pulumi:"apiVersion"`
	Kind       *string            `pulumi:"kind"`
	Metadata   *metav1.ObjectMeta `pulumi:"metadata"`
	// GlobalConfigurationSpec is the spec of the GlobalConfiguration resource.
	Spec *GlobalConfigurationSpec `pulumi:"spec"`
}

type GlobalConfigurationState struct {
	ApiVersion pulumi.StringPtrInput
	Kind       pulumi.StringPtrInput
	Metadata   metav1.ObjectMetaPtrInput
	// GlobalConfigurationSpec is the spec of the GlobalConfiguration resource.
	Spec GlobalConfigurationSpecPtrInput
}

func (GlobalConfigurationState) ElementType() reflect.Type {
	return reflect.TypeOf((*globalConfigurationState)(nil)).Elem()
}

type globalConfigurationArgs struct {
	ApiVersion *string            `pulumi:"apiVersion"`
	Kind       *string            `pulumi:"kind"`
	Metadata   *metav1.ObjectMeta `pulumi:"metadata"`
	// GlobalConfigurationSpec is the spec of the GlobalConfiguration resource.
	Spec *GlobalConfigurationSpec `pulumi:"spec"`
}

// The set of arguments for constructing a GlobalConfiguration resource.
type GlobalConfigurationArgs struct {
	ApiVersion pulumi.StringPtrInput
	Kind       pulumi.StringPtrInput
	Metadata   metav1.ObjectMetaPtrInput
	// GlobalConfigurationSpec is the spec of the GlobalConfiguration resource.
	Spec GlobalConfigurationSpecPtrInput
}

func (GlobalConfigurationArgs) ElementType() reflect.Type {
	return reflect.TypeOf((*globalConfigurationArgs)(nil)).Elem()
}

type GlobalConfigurationInput interface {
	pulumi.Input

	ToGlobalConfigurationOutput() GlobalConfigurationOutput
	ToGlobalConfigurationOutputWithContext(ctx context.Context) GlobalConfigurationOutput
}

func (*GlobalConfiguration) ElementType() reflect.Type {
	return reflect.TypeOf((*GlobalConfiguration)(nil))
}

func (i *GlobalConfiguration) ToGlobalConfigurationOutput() GlobalConfigurationOutput {
	return i.ToGlobalConfigurationOutputWithContext(context.Background())
}

func (i *GlobalConfiguration) ToGlobalConfigurationOutputWithContext(ctx context.Context) GlobalConfigurationOutput {
	return pulumi.ToOutputWithContext(ctx, i).(GlobalConfigurationOutput)
}

type GlobalConfigurationOutput struct {
	*pulumi.OutputState
}

func (GlobalConfigurationOutput) ElementType() reflect.Type {
	return reflect.TypeOf((*GlobalConfiguration)(nil))
}

func (o GlobalConfigurationOutput) ToGlobalConfigurationOutput() GlobalConfigurationOutput {
	return o
}

func (o GlobalConfigurationOutput) ToGlobalConfigurationOutputWithContext(ctx context.Context) GlobalConfigurationOutput {
	return o
}

func init() {
	pulumi.RegisterOutputType(GlobalConfigurationOutput{})
}
