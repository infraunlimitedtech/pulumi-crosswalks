package libvirt

import (
	"fmt"

	"github.com/pulumi/pulumi-libvirt/sdk/go/libvirt"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func (i *ComputeConfig) manage(ctx *pulumi.Context, cfg *HypervisorConfig) (map[string]map[string]interface{}, error) {
	provider, err := libvirt.NewProvider(ctx, fmt.Sprintf("%s-provider", cfg.Name), &libvirt.ProviderArgs{
		Uri: pulumi.String(cfg.URI),
	})
	if err != nil {
		return nil, err
	}

	pool, err := libvirt.NewPool(ctx, fmt.Sprintf("%s-defaultPool", cfg.Name), &libvirt.PoolArgs{
		Name: pulumi.String(i.Storage.Name),
		Type: pulumi.String("dir"),
		Path: pulumi.Sprintf("/var/lib/libvirt/pools/%s", i.Storage.Name),
	}, pulumi.Provider(provider))
	if err != nil {
		return nil, err
	}

	net, err := libvirt.NewNetwork(ctx, fmt.Sprintf("%s-defaultNetwork", cfg.Name), &libvirt.NetworkArgs{
		Name:      pulumi.String(i.Network.Name),
		Autostart: pulumi.Bool(true),
		Mode:      pulumi.String("nat"),
		Addresses: pulumi.StringArray{
			pulumi.String(cfg.Network.CIDR),
		},
		Dns: &libvirt.NetworkDnsArgs{
			Enabled: pulumi.Bool(true),
		},
	}, pulumi.Provider(provider))
	if err != nil {
		return nil, err
	}

	nodes := make(map[string]map[string]interface{})

	for _, vm := range cfg.Machines {
		combustion, err := libvirt.NewVolume(ctx, fmt.Sprintf("%s-%s-combustionVolume", cfg.Name, vm.ID), &libvirt.VolumeArgs{
			Pool:   pool.Name,
			Name:   pulumi.Sprintf("%s-microos-combustion", vm.ID),
			Source: pulumi.String(i.Images.Combustion),
		}, pulumi.Provider(provider), pulumi.Protect(false))
		if err != nil {
			return nil, err
		}
		base, err := libvirt.NewVolume(ctx, fmt.Sprintf("%s-%s-volume", cfg.Name, vm.ID), &libvirt.VolumeArgs{
			Pool:   pool.Name,
			Name:   pulumi.String(vm.ID),
			Source: pulumi.String(i.Images.Base),
		}, pulumi.Provider(provider), pulumi.Protect(false))
		if err != nil {
			return nil, err
		}

		domain, err := libvirt.NewDomain(ctx, fmt.Sprintf("%s-%s", cfg.Name, vm.ID), &libvirt.DomainArgs{
			Name:      pulumi.String(vm.ID),
			Vcpu:      pulumi.Int(2),
			Memory:    pulumi.Int(2048),
			Autostart: pulumi.Bool(true),
			QemuAgent: pulumi.Bool(true),
			BootDevices: &libvirt.DomainBootDeviceArray{
				&libvirt.DomainBootDeviceArgs{
					Devs: pulumi.StringArray{
						pulumi.String("hd"),
					},
				},
			},
			NetworkInterfaces: &libvirt.DomainNetworkInterfaceArray{
				&libvirt.DomainNetworkInterfaceArgs{
					NetworkName:  net.Name,
					Hostname:     pulumi.String(vm.ID),
					WaitForLease: pulumi.Bool(true),
				},
			},
			Disks: &libvirt.DomainDiskArray{
				&libvirt.DomainDiskArgs{
					VolumeId: base.ID(),
				},
				&libvirt.DomainDiskArgs{
					VolumeId: combustion.ID(),
				},
			},
		}, pulumi.Provider(provider))
		if err != nil {
			return nil, err
		}
		nodes[vm.ID] = make(map[string]interface{})
		nodes[vm.ID]["ip"] = domain.NetworkInterfaces.Index(pulumi.Int(0)).Addresses().Index(pulumi.Int(0))
	}
	return nodes, nil
}
