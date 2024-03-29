// *** WARNING: this file was generated by crd2pulumi. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

package v1alpha1

import (
	"context"
	"reflect"

	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Peer is a WireGuard peer that should have access to the VPN.
type PeerType struct {
	ApiVersion *string            `pulumi:"apiVersion"`
	Kind       *string            `pulumi:"kind"`
	Metadata   *metav1.ObjectMeta `pulumi:"metadata"`
	// Specification of the desired behavior of the Kilo Peer. More info: https://github.com/kubernetes/community/blob/master/contributors/devel/api-conventions.md#spec-and-status
	Spec PeerSpec `pulumi:"spec"`
}

// PeerTypeInput is an input type that accepts PeerTypeArgs and PeerTypeOutput values.
// You can construct a concrete instance of `PeerTypeInput` via:
//
//          PeerTypeArgs{...}
type PeerTypeInput interface {
	pulumi.Input

	ToPeerTypeOutput() PeerTypeOutput
	ToPeerTypeOutputWithContext(context.Context) PeerTypeOutput
}

// Peer is a WireGuard peer that should have access to the VPN.
type PeerTypeArgs struct {
	ApiVersion pulumi.StringPtrInput     `pulumi:"apiVersion"`
	Kind       pulumi.StringPtrInput     `pulumi:"kind"`
	Metadata   metav1.ObjectMetaPtrInput `pulumi:"metadata"`
	// Specification of the desired behavior of the Kilo Peer. More info: https://github.com/kubernetes/community/blob/master/contributors/devel/api-conventions.md#spec-and-status
	Spec PeerSpecInput `pulumi:"spec"`
}

func (PeerTypeArgs) ElementType() reflect.Type {
	return reflect.TypeOf((*PeerType)(nil)).Elem()
}

func (i PeerTypeArgs) ToPeerTypeOutput() PeerTypeOutput {
	return i.ToPeerTypeOutputWithContext(context.Background())
}

func (i PeerTypeArgs) ToPeerTypeOutputWithContext(ctx context.Context) PeerTypeOutput {
	return pulumi.ToOutputWithContext(ctx, i).(PeerTypeOutput)
}

// Peer is a WireGuard peer that should have access to the VPN.
type PeerTypeOutput struct{ *pulumi.OutputState }

func (PeerTypeOutput) ElementType() reflect.Type {
	return reflect.TypeOf((*PeerType)(nil)).Elem()
}

func (o PeerTypeOutput) ToPeerTypeOutput() PeerTypeOutput {
	return o
}

func (o PeerTypeOutput) ToPeerTypeOutputWithContext(ctx context.Context) PeerTypeOutput {
	return o
}

func (o PeerTypeOutput) ApiVersion() pulumi.StringPtrOutput {
	return o.ApplyT(func(v PeerType) *string { return v.ApiVersion }).(pulumi.StringPtrOutput)
}

func (o PeerTypeOutput) Kind() pulumi.StringPtrOutput {
	return o.ApplyT(func(v PeerType) *string { return v.Kind }).(pulumi.StringPtrOutput)
}

func (o PeerTypeOutput) Metadata() metav1.ObjectMetaPtrOutput {
	return o.ApplyT(func(v PeerType) *metav1.ObjectMeta { return v.Metadata }).(metav1.ObjectMetaPtrOutput)
}

// Specification of the desired behavior of the Kilo Peer. More info: https://github.com/kubernetes/community/blob/master/contributors/devel/api-conventions.md#spec-and-status
func (o PeerTypeOutput) Spec() PeerSpecOutput {
	return o.ApplyT(func(v PeerType) PeerSpec { return v.Spec }).(PeerSpecOutput)
}

type PeerMetadata struct {
}

// PeerMetadataInput is an input type that accepts PeerMetadataArgs and PeerMetadataOutput values.
// You can construct a concrete instance of `PeerMetadataInput` via:
//
//          PeerMetadataArgs{...}
type PeerMetadataInput interface {
	pulumi.Input

	ToPeerMetadataOutput() PeerMetadataOutput
	ToPeerMetadataOutputWithContext(context.Context) PeerMetadataOutput
}

type PeerMetadataArgs struct {
}

func (PeerMetadataArgs) ElementType() reflect.Type {
	return reflect.TypeOf((*PeerMetadata)(nil)).Elem()
}

func (i PeerMetadataArgs) ToPeerMetadataOutput() PeerMetadataOutput {
	return i.ToPeerMetadataOutputWithContext(context.Background())
}

func (i PeerMetadataArgs) ToPeerMetadataOutputWithContext(ctx context.Context) PeerMetadataOutput {
	return pulumi.ToOutputWithContext(ctx, i).(PeerMetadataOutput)
}

type PeerMetadataOutput struct{ *pulumi.OutputState }

func (PeerMetadataOutput) ElementType() reflect.Type {
	return reflect.TypeOf((*PeerMetadata)(nil)).Elem()
}

func (o PeerMetadataOutput) ToPeerMetadataOutput() PeerMetadataOutput {
	return o
}

func (o PeerMetadataOutput) ToPeerMetadataOutputWithContext(ctx context.Context) PeerMetadataOutput {
	return o
}

// Specification of the desired behavior of the Kilo Peer. More info: https://github.com/kubernetes/community/blob/master/contributors/devel/api-conventions.md#spec-and-status
type PeerSpec struct {
	// AllowedIPs is the list of IP addresses that are allowed for the given peer's tunnel.
	AllowedIPs []string `pulumi:"allowedIPs"`
	// Endpoint is the initial endpoint for connections to the peer.
	Endpoint *PeerSpecEndpoint `pulumi:"endpoint"`
	// PersistentKeepalive is the interval in seconds of the emission of keepalive packets by the peer. This defaults to 0, which disables the feature.
	PersistentKeepalive *int `pulumi:"persistentKeepalive"`
	// PresharedKey is the optional symmetric encryption key for the peer.
	PresharedKey *string `pulumi:"presharedKey"`
	// PublicKey is the WireGuard public key for the peer.
	PublicKey string `pulumi:"publicKey"`
}

// PeerSpecInput is an input type that accepts PeerSpecArgs and PeerSpecOutput values.
// You can construct a concrete instance of `PeerSpecInput` via:
//
//          PeerSpecArgs{...}
type PeerSpecInput interface {
	pulumi.Input

	ToPeerSpecOutput() PeerSpecOutput
	ToPeerSpecOutputWithContext(context.Context) PeerSpecOutput
}

// Specification of the desired behavior of the Kilo Peer. More info: https://github.com/kubernetes/community/blob/master/contributors/devel/api-conventions.md#spec-and-status
type PeerSpecArgs struct {
	// AllowedIPs is the list of IP addresses that are allowed for the given peer's tunnel.
	AllowedIPs pulumi.StringArrayInput `pulumi:"allowedIPs"`
	// Endpoint is the initial endpoint for connections to the peer.
	Endpoint PeerSpecEndpointPtrInput `pulumi:"endpoint"`
	// PersistentKeepalive is the interval in seconds of the emission of keepalive packets by the peer. This defaults to 0, which disables the feature.
	PersistentKeepalive pulumi.IntPtrInput `pulumi:"persistentKeepalive"`
	// PresharedKey is the optional symmetric encryption key for the peer.
	PresharedKey pulumi.StringPtrInput `pulumi:"presharedKey"`
	// PublicKey is the WireGuard public key for the peer.
	PublicKey pulumi.StringInput `pulumi:"publicKey"`
}

func (PeerSpecArgs) ElementType() reflect.Type {
	return reflect.TypeOf((*PeerSpec)(nil)).Elem()
}

func (i PeerSpecArgs) ToPeerSpecOutput() PeerSpecOutput {
	return i.ToPeerSpecOutputWithContext(context.Background())
}

func (i PeerSpecArgs) ToPeerSpecOutputWithContext(ctx context.Context) PeerSpecOutput {
	return pulumi.ToOutputWithContext(ctx, i).(PeerSpecOutput)
}

func (i PeerSpecArgs) ToPeerSpecPtrOutput() PeerSpecPtrOutput {
	return i.ToPeerSpecPtrOutputWithContext(context.Background())
}

func (i PeerSpecArgs) ToPeerSpecPtrOutputWithContext(ctx context.Context) PeerSpecPtrOutput {
	return pulumi.ToOutputWithContext(ctx, i).(PeerSpecOutput).ToPeerSpecPtrOutputWithContext(ctx)
}

// PeerSpecPtrInput is an input type that accepts PeerSpecArgs, PeerSpecPtr and PeerSpecPtrOutput values.
// You can construct a concrete instance of `PeerSpecPtrInput` via:
//
//          PeerSpecArgs{...}
//
//  or:
//
//          nil
type PeerSpecPtrInput interface {
	pulumi.Input

	ToPeerSpecPtrOutput() PeerSpecPtrOutput
	ToPeerSpecPtrOutputWithContext(context.Context) PeerSpecPtrOutput
}

type peerSpecPtrType PeerSpecArgs

func PeerSpecPtr(v *PeerSpecArgs) PeerSpecPtrInput {
	return (*peerSpecPtrType)(v)
}

func (*peerSpecPtrType) ElementType() reflect.Type {
	return reflect.TypeOf((**PeerSpec)(nil)).Elem()
}

func (i *peerSpecPtrType) ToPeerSpecPtrOutput() PeerSpecPtrOutput {
	return i.ToPeerSpecPtrOutputWithContext(context.Background())
}

func (i *peerSpecPtrType) ToPeerSpecPtrOutputWithContext(ctx context.Context) PeerSpecPtrOutput {
	return pulumi.ToOutputWithContext(ctx, i).(PeerSpecPtrOutput)
}

// Specification of the desired behavior of the Kilo Peer. More info: https://github.com/kubernetes/community/blob/master/contributors/devel/api-conventions.md#spec-and-status
type PeerSpecOutput struct{ *pulumi.OutputState }

func (PeerSpecOutput) ElementType() reflect.Type {
	return reflect.TypeOf((*PeerSpec)(nil)).Elem()
}

func (o PeerSpecOutput) ToPeerSpecOutput() PeerSpecOutput {
	return o
}

func (o PeerSpecOutput) ToPeerSpecOutputWithContext(ctx context.Context) PeerSpecOutput {
	return o
}

func (o PeerSpecOutput) ToPeerSpecPtrOutput() PeerSpecPtrOutput {
	return o.ToPeerSpecPtrOutputWithContext(context.Background())
}

func (o PeerSpecOutput) ToPeerSpecPtrOutputWithContext(ctx context.Context) PeerSpecPtrOutput {
	return o.ApplyT(func(v PeerSpec) *PeerSpec {
		return &v
	}).(PeerSpecPtrOutput)
}

// AllowedIPs is the list of IP addresses that are allowed for the given peer's tunnel.
func (o PeerSpecOutput) AllowedIPs() pulumi.StringArrayOutput {
	return o.ApplyT(func(v PeerSpec) []string { return v.AllowedIPs }).(pulumi.StringArrayOutput)
}

// Endpoint is the initial endpoint for connections to the peer.
func (o PeerSpecOutput) Endpoint() PeerSpecEndpointPtrOutput {
	return o.ApplyT(func(v PeerSpec) *PeerSpecEndpoint { return v.Endpoint }).(PeerSpecEndpointPtrOutput)
}

// PersistentKeepalive is the interval in seconds of the emission of keepalive packets by the peer. This defaults to 0, which disables the feature.
func (o PeerSpecOutput) PersistentKeepalive() pulumi.IntPtrOutput {
	return o.ApplyT(func(v PeerSpec) *int { return v.PersistentKeepalive }).(pulumi.IntPtrOutput)
}

// PresharedKey is the optional symmetric encryption key for the peer.
func (o PeerSpecOutput) PresharedKey() pulumi.StringPtrOutput {
	return o.ApplyT(func(v PeerSpec) *string { return v.PresharedKey }).(pulumi.StringPtrOutput)
}

// PublicKey is the WireGuard public key for the peer.
func (o PeerSpecOutput) PublicKey() pulumi.StringOutput {
	return o.ApplyT(func(v PeerSpec) string { return v.PublicKey }).(pulumi.StringOutput)
}

type PeerSpecPtrOutput struct{ *pulumi.OutputState }

func (PeerSpecPtrOutput) ElementType() reflect.Type {
	return reflect.TypeOf((**PeerSpec)(nil)).Elem()
}

func (o PeerSpecPtrOutput) ToPeerSpecPtrOutput() PeerSpecPtrOutput {
	return o
}

func (o PeerSpecPtrOutput) ToPeerSpecPtrOutputWithContext(ctx context.Context) PeerSpecPtrOutput {
	return o
}

func (o PeerSpecPtrOutput) Elem() PeerSpecOutput {
	return o.ApplyT(func(v *PeerSpec) PeerSpec { return *v }).(PeerSpecOutput)
}

// AllowedIPs is the list of IP addresses that are allowed for the given peer's tunnel.
func (o PeerSpecPtrOutput) AllowedIPs() pulumi.StringArrayOutput {
	return o.ApplyT(func(v *PeerSpec) []string {
		if v == nil {
			return nil
		}
		return v.AllowedIPs
	}).(pulumi.StringArrayOutput)
}

// Endpoint is the initial endpoint for connections to the peer.
func (o PeerSpecPtrOutput) Endpoint() PeerSpecEndpointPtrOutput {
	return o.ApplyT(func(v *PeerSpec) *PeerSpecEndpoint {
		if v == nil {
			return nil
		}
		return v.Endpoint
	}).(PeerSpecEndpointPtrOutput)
}

// PersistentKeepalive is the interval in seconds of the emission of keepalive packets by the peer. This defaults to 0, which disables the feature.
func (o PeerSpecPtrOutput) PersistentKeepalive() pulumi.IntPtrOutput {
	return o.ApplyT(func(v *PeerSpec) *int {
		if v == nil {
			return nil
		}
		return v.PersistentKeepalive
	}).(pulumi.IntPtrOutput)
}

// PresharedKey is the optional symmetric encryption key for the peer.
func (o PeerSpecPtrOutput) PresharedKey() pulumi.StringPtrOutput {
	return o.ApplyT(func(v *PeerSpec) *string {
		if v == nil {
			return nil
		}
		return v.PresharedKey
	}).(pulumi.StringPtrOutput)
}

// PublicKey is the WireGuard public key for the peer.
func (o PeerSpecPtrOutput) PublicKey() pulumi.StringPtrOutput {
	return o.ApplyT(func(v *PeerSpec) *string {
		if v == nil {
			return nil
		}
		return &v.PublicKey
	}).(pulumi.StringPtrOutput)
}

// Endpoint is the initial endpoint for connections to the peer.
type PeerSpecEndpoint struct {
	// DNSOrIP is a DNS name or an IP address.
	DnsOrIP PeerSpecEndpointDnsOrIP `pulumi:"dnsOrIP"`
	// Port must be a valid port number.
	Port int `pulumi:"port"`
}

// PeerSpecEndpointInput is an input type that accepts PeerSpecEndpointArgs and PeerSpecEndpointOutput values.
// You can construct a concrete instance of `PeerSpecEndpointInput` via:
//
//          PeerSpecEndpointArgs{...}
type PeerSpecEndpointInput interface {
	pulumi.Input

	ToPeerSpecEndpointOutput() PeerSpecEndpointOutput
	ToPeerSpecEndpointOutputWithContext(context.Context) PeerSpecEndpointOutput
}

// Endpoint is the initial endpoint for connections to the peer.
type PeerSpecEndpointArgs struct {
	// DNSOrIP is a DNS name or an IP address.
	DnsOrIP PeerSpecEndpointDnsOrIPInput `pulumi:"dnsOrIP"`
	// Port must be a valid port number.
	Port pulumi.IntInput `pulumi:"port"`
}

func (PeerSpecEndpointArgs) ElementType() reflect.Type {
	return reflect.TypeOf((*PeerSpecEndpoint)(nil)).Elem()
}

func (i PeerSpecEndpointArgs) ToPeerSpecEndpointOutput() PeerSpecEndpointOutput {
	return i.ToPeerSpecEndpointOutputWithContext(context.Background())
}

func (i PeerSpecEndpointArgs) ToPeerSpecEndpointOutputWithContext(ctx context.Context) PeerSpecEndpointOutput {
	return pulumi.ToOutputWithContext(ctx, i).(PeerSpecEndpointOutput)
}

func (i PeerSpecEndpointArgs) ToPeerSpecEndpointPtrOutput() PeerSpecEndpointPtrOutput {
	return i.ToPeerSpecEndpointPtrOutputWithContext(context.Background())
}

func (i PeerSpecEndpointArgs) ToPeerSpecEndpointPtrOutputWithContext(ctx context.Context) PeerSpecEndpointPtrOutput {
	return pulumi.ToOutputWithContext(ctx, i).(PeerSpecEndpointOutput).ToPeerSpecEndpointPtrOutputWithContext(ctx)
}

// PeerSpecEndpointPtrInput is an input type that accepts PeerSpecEndpointArgs, PeerSpecEndpointPtr and PeerSpecEndpointPtrOutput values.
// You can construct a concrete instance of `PeerSpecEndpointPtrInput` via:
//
//          PeerSpecEndpointArgs{...}
//
//  or:
//
//          nil
type PeerSpecEndpointPtrInput interface {
	pulumi.Input

	ToPeerSpecEndpointPtrOutput() PeerSpecEndpointPtrOutput
	ToPeerSpecEndpointPtrOutputWithContext(context.Context) PeerSpecEndpointPtrOutput
}

type peerSpecEndpointPtrType PeerSpecEndpointArgs

func PeerSpecEndpointPtr(v *PeerSpecEndpointArgs) PeerSpecEndpointPtrInput {
	return (*peerSpecEndpointPtrType)(v)
}

func (*peerSpecEndpointPtrType) ElementType() reflect.Type {
	return reflect.TypeOf((**PeerSpecEndpoint)(nil)).Elem()
}

func (i *peerSpecEndpointPtrType) ToPeerSpecEndpointPtrOutput() PeerSpecEndpointPtrOutput {
	return i.ToPeerSpecEndpointPtrOutputWithContext(context.Background())
}

func (i *peerSpecEndpointPtrType) ToPeerSpecEndpointPtrOutputWithContext(ctx context.Context) PeerSpecEndpointPtrOutput {
	return pulumi.ToOutputWithContext(ctx, i).(PeerSpecEndpointPtrOutput)
}

// Endpoint is the initial endpoint for connections to the peer.
type PeerSpecEndpointOutput struct{ *pulumi.OutputState }

func (PeerSpecEndpointOutput) ElementType() reflect.Type {
	return reflect.TypeOf((*PeerSpecEndpoint)(nil)).Elem()
}

func (o PeerSpecEndpointOutput) ToPeerSpecEndpointOutput() PeerSpecEndpointOutput {
	return o
}

func (o PeerSpecEndpointOutput) ToPeerSpecEndpointOutputWithContext(ctx context.Context) PeerSpecEndpointOutput {
	return o
}

func (o PeerSpecEndpointOutput) ToPeerSpecEndpointPtrOutput() PeerSpecEndpointPtrOutput {
	return o.ToPeerSpecEndpointPtrOutputWithContext(context.Background())
}

func (o PeerSpecEndpointOutput) ToPeerSpecEndpointPtrOutputWithContext(ctx context.Context) PeerSpecEndpointPtrOutput {
	return o.ApplyT(func(v PeerSpecEndpoint) *PeerSpecEndpoint {
		return &v
	}).(PeerSpecEndpointPtrOutput)
}

// DNSOrIP is a DNS name or an IP address.
func (o PeerSpecEndpointOutput) DnsOrIP() PeerSpecEndpointDnsOrIPOutput {
	return o.ApplyT(func(v PeerSpecEndpoint) PeerSpecEndpointDnsOrIP { return v.DnsOrIP }).(PeerSpecEndpointDnsOrIPOutput)
}

// Port must be a valid port number.
func (o PeerSpecEndpointOutput) Port() pulumi.IntOutput {
	return o.ApplyT(func(v PeerSpecEndpoint) int { return v.Port }).(pulumi.IntOutput)
}

type PeerSpecEndpointPtrOutput struct{ *pulumi.OutputState }

func (PeerSpecEndpointPtrOutput) ElementType() reflect.Type {
	return reflect.TypeOf((**PeerSpecEndpoint)(nil)).Elem()
}

func (o PeerSpecEndpointPtrOutput) ToPeerSpecEndpointPtrOutput() PeerSpecEndpointPtrOutput {
	return o
}

func (o PeerSpecEndpointPtrOutput) ToPeerSpecEndpointPtrOutputWithContext(ctx context.Context) PeerSpecEndpointPtrOutput {
	return o
}

func (o PeerSpecEndpointPtrOutput) Elem() PeerSpecEndpointOutput {
	return o.ApplyT(func(v *PeerSpecEndpoint) PeerSpecEndpoint { return *v }).(PeerSpecEndpointOutput)
}

// DNSOrIP is a DNS name or an IP address.
func (o PeerSpecEndpointPtrOutput) DnsOrIP() PeerSpecEndpointDnsOrIPPtrOutput {
	return o.ApplyT(func(v *PeerSpecEndpoint) *PeerSpecEndpointDnsOrIP {
		if v == nil {
			return nil
		}
		return &v.DnsOrIP
	}).(PeerSpecEndpointDnsOrIPPtrOutput)
}

// Port must be a valid port number.
func (o PeerSpecEndpointPtrOutput) Port() pulumi.IntPtrOutput {
	return o.ApplyT(func(v *PeerSpecEndpoint) *int {
		if v == nil {
			return nil
		}
		return &v.Port
	}).(pulumi.IntPtrOutput)
}

// DNSOrIP is a DNS name or an IP address.
type PeerSpecEndpointDnsOrIP struct {
	// DNS must be a valid RFC 1123 subdomain.
	Dns *string `pulumi:"dns"`
	// IP must be a valid IP address.
	Ip *string `pulumi:"ip"`
}

// PeerSpecEndpointDnsOrIPInput is an input type that accepts PeerSpecEndpointDnsOrIPArgs and PeerSpecEndpointDnsOrIPOutput values.
// You can construct a concrete instance of `PeerSpecEndpointDnsOrIPInput` via:
//
//          PeerSpecEndpointDnsOrIPArgs{...}
type PeerSpecEndpointDnsOrIPInput interface {
	pulumi.Input

	ToPeerSpecEndpointDnsOrIPOutput() PeerSpecEndpointDnsOrIPOutput
	ToPeerSpecEndpointDnsOrIPOutputWithContext(context.Context) PeerSpecEndpointDnsOrIPOutput
}

// DNSOrIP is a DNS name or an IP address.
type PeerSpecEndpointDnsOrIPArgs struct {
	// DNS must be a valid RFC 1123 subdomain.
	Dns pulumi.StringPtrInput `pulumi:"dns"`
	// IP must be a valid IP address.
	Ip pulumi.StringPtrInput `pulumi:"ip"`
}

func (PeerSpecEndpointDnsOrIPArgs) ElementType() reflect.Type {
	return reflect.TypeOf((*PeerSpecEndpointDnsOrIP)(nil)).Elem()
}

func (i PeerSpecEndpointDnsOrIPArgs) ToPeerSpecEndpointDnsOrIPOutput() PeerSpecEndpointDnsOrIPOutput {
	return i.ToPeerSpecEndpointDnsOrIPOutputWithContext(context.Background())
}

func (i PeerSpecEndpointDnsOrIPArgs) ToPeerSpecEndpointDnsOrIPOutputWithContext(ctx context.Context) PeerSpecEndpointDnsOrIPOutput {
	return pulumi.ToOutputWithContext(ctx, i).(PeerSpecEndpointDnsOrIPOutput)
}

func (i PeerSpecEndpointDnsOrIPArgs) ToPeerSpecEndpointDnsOrIPPtrOutput() PeerSpecEndpointDnsOrIPPtrOutput {
	return i.ToPeerSpecEndpointDnsOrIPPtrOutputWithContext(context.Background())
}

func (i PeerSpecEndpointDnsOrIPArgs) ToPeerSpecEndpointDnsOrIPPtrOutputWithContext(ctx context.Context) PeerSpecEndpointDnsOrIPPtrOutput {
	return pulumi.ToOutputWithContext(ctx, i).(PeerSpecEndpointDnsOrIPOutput).ToPeerSpecEndpointDnsOrIPPtrOutputWithContext(ctx)
}

// PeerSpecEndpointDnsOrIPPtrInput is an input type that accepts PeerSpecEndpointDnsOrIPArgs, PeerSpecEndpointDnsOrIPPtr and PeerSpecEndpointDnsOrIPPtrOutput values.
// You can construct a concrete instance of `PeerSpecEndpointDnsOrIPPtrInput` via:
//
//          PeerSpecEndpointDnsOrIPArgs{...}
//
//  or:
//
//          nil
type PeerSpecEndpointDnsOrIPPtrInput interface {
	pulumi.Input

	ToPeerSpecEndpointDnsOrIPPtrOutput() PeerSpecEndpointDnsOrIPPtrOutput
	ToPeerSpecEndpointDnsOrIPPtrOutputWithContext(context.Context) PeerSpecEndpointDnsOrIPPtrOutput
}

type peerSpecEndpointDnsOrIPPtrType PeerSpecEndpointDnsOrIPArgs

func PeerSpecEndpointDnsOrIPPtr(v *PeerSpecEndpointDnsOrIPArgs) PeerSpecEndpointDnsOrIPPtrInput {
	return (*peerSpecEndpointDnsOrIPPtrType)(v)
}

func (*peerSpecEndpointDnsOrIPPtrType) ElementType() reflect.Type {
	return reflect.TypeOf((**PeerSpecEndpointDnsOrIP)(nil)).Elem()
}

func (i *peerSpecEndpointDnsOrIPPtrType) ToPeerSpecEndpointDnsOrIPPtrOutput() PeerSpecEndpointDnsOrIPPtrOutput {
	return i.ToPeerSpecEndpointDnsOrIPPtrOutputWithContext(context.Background())
}

func (i *peerSpecEndpointDnsOrIPPtrType) ToPeerSpecEndpointDnsOrIPPtrOutputWithContext(ctx context.Context) PeerSpecEndpointDnsOrIPPtrOutput {
	return pulumi.ToOutputWithContext(ctx, i).(PeerSpecEndpointDnsOrIPPtrOutput)
}

// DNSOrIP is a DNS name or an IP address.
type PeerSpecEndpointDnsOrIPOutput struct{ *pulumi.OutputState }

func (PeerSpecEndpointDnsOrIPOutput) ElementType() reflect.Type {
	return reflect.TypeOf((*PeerSpecEndpointDnsOrIP)(nil)).Elem()
}

func (o PeerSpecEndpointDnsOrIPOutput) ToPeerSpecEndpointDnsOrIPOutput() PeerSpecEndpointDnsOrIPOutput {
	return o
}

func (o PeerSpecEndpointDnsOrIPOutput) ToPeerSpecEndpointDnsOrIPOutputWithContext(ctx context.Context) PeerSpecEndpointDnsOrIPOutput {
	return o
}

func (o PeerSpecEndpointDnsOrIPOutput) ToPeerSpecEndpointDnsOrIPPtrOutput() PeerSpecEndpointDnsOrIPPtrOutput {
	return o.ToPeerSpecEndpointDnsOrIPPtrOutputWithContext(context.Background())
}

func (o PeerSpecEndpointDnsOrIPOutput) ToPeerSpecEndpointDnsOrIPPtrOutputWithContext(ctx context.Context) PeerSpecEndpointDnsOrIPPtrOutput {
	return o.ApplyT(func(v PeerSpecEndpointDnsOrIP) *PeerSpecEndpointDnsOrIP {
		return &v
	}).(PeerSpecEndpointDnsOrIPPtrOutput)
}

// DNS must be a valid RFC 1123 subdomain.
func (o PeerSpecEndpointDnsOrIPOutput) Dns() pulumi.StringPtrOutput {
	return o.ApplyT(func(v PeerSpecEndpointDnsOrIP) *string { return v.Dns }).(pulumi.StringPtrOutput)
}

// IP must be a valid IP address.
func (o PeerSpecEndpointDnsOrIPOutput) Ip() pulumi.StringPtrOutput {
	return o.ApplyT(func(v PeerSpecEndpointDnsOrIP) *string { return v.Ip }).(pulumi.StringPtrOutput)
}

type PeerSpecEndpointDnsOrIPPtrOutput struct{ *pulumi.OutputState }

func (PeerSpecEndpointDnsOrIPPtrOutput) ElementType() reflect.Type {
	return reflect.TypeOf((**PeerSpecEndpointDnsOrIP)(nil)).Elem()
}

func (o PeerSpecEndpointDnsOrIPPtrOutput) ToPeerSpecEndpointDnsOrIPPtrOutput() PeerSpecEndpointDnsOrIPPtrOutput {
	return o
}

func (o PeerSpecEndpointDnsOrIPPtrOutput) ToPeerSpecEndpointDnsOrIPPtrOutputWithContext(ctx context.Context) PeerSpecEndpointDnsOrIPPtrOutput {
	return o
}

func (o PeerSpecEndpointDnsOrIPPtrOutput) Elem() PeerSpecEndpointDnsOrIPOutput {
	return o.ApplyT(func(v *PeerSpecEndpointDnsOrIP) PeerSpecEndpointDnsOrIP { return *v }).(PeerSpecEndpointDnsOrIPOutput)
}

// DNS must be a valid RFC 1123 subdomain.
func (o PeerSpecEndpointDnsOrIPPtrOutput) Dns() pulumi.StringPtrOutput {
	return o.ApplyT(func(v *PeerSpecEndpointDnsOrIP) *string {
		if v == nil {
			return nil
		}
		return v.Dns
	}).(pulumi.StringPtrOutput)
}

// IP must be a valid IP address.
func (o PeerSpecEndpointDnsOrIPPtrOutput) Ip() pulumi.StringPtrOutput {
	return o.ApplyT(func(v *PeerSpecEndpointDnsOrIP) *string {
		if v == nil {
			return nil
		}
		return v.Ip
	}).(pulumi.StringPtrOutput)
}

func init() {
	pulumi.RegisterOutputType(PeerTypeOutput{})
	pulumi.RegisterOutputType(PeerMetadataOutput{})
	pulumi.RegisterOutputType(PeerSpecOutput{})
	pulumi.RegisterOutputType(PeerSpecPtrOutput{})
	pulumi.RegisterOutputType(PeerSpecEndpointOutput{})
	pulumi.RegisterOutputType(PeerSpecEndpointPtrOutput{})
	pulumi.RegisterOutputType(PeerSpecEndpointDnsOrIPOutput{})
	pulumi.RegisterOutputType(PeerSpecEndpointDnsOrIPPtrOutput{})
}
