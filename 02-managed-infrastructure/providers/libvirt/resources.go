package libvirt

import (
	"fmt"
	"github.com/pulumi/pulumi-libvirt/sdk/go/libvirt"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"os"
)

func manageLibvirtHost(ctx *pulumi.Context, cfg HypervisorConfig) (map[string]map[string]interface{}, error) {
	provider, err := libvirt.NewProvider(ctx, fmt.Sprintf("%s-provider", cfg.Name), &libvirt.ProviderArgs{
		Uri: pulumi.String(cfg.URI),
	})

	if err != nil {
		return nil, err
	}

	pool, err := libvirt.NewPool(ctx, fmt.Sprintf("%s-defaultPool", cfg.Name), &libvirt.PoolArgs{
		Name: pulumi.String("infraunlimited"),
		Type: pulumi.String("dir"),
		Path: pulumi.String("/var/lib/libvirt/pools/infraunlimited"),
	}, pulumi.Provider(provider))

	if err != nil {
		return nil, err
	}

	networkXlst, err := os.ReadFile("providers/libvirt/network.xlst")
	if err != nil {
		return nil, err
	}

	net, err := libvirt.NewNetwork(ctx, fmt.Sprintf("%s-defaultNetwork", cfg.Name), &libvirt.NetworkArgs{
		Name:      pulumi.String("infraunlimited-routed"),
		Autostart: pulumi.Bool(true),
		Mode:      pulumi.String("route"),
		Bridge:    pulumi.String("qemu-br0"),
		Addresses: pulumi.StringArray{
			pulumi.String(cfg.NetworkCIDR),
		},
		Dns: &libvirt.NetworkDnsArgs{
			Enabled: pulumi.Bool(true),
		},
		Xml: &libvirt.NetworkXmlArgs{
			Xslt: pulumi.String(networkXlst),
		},
	}, pulumi.Provider(provider))
	if err != nil {
		return nil, err
	}

	nodes := make(map[string]map[string]interface{})

	for _, vm := range cfg.Machines {
		volume, err := libvirt.NewVolume(ctx, fmt.Sprintf("%s-%s-volume", cfg.Name, vm.ID), &libvirt.VolumeArgs{
			Pool:   pool.Name,
			Name:   pulumi.String(vm.ID),
			Format: pulumi.String("raw"),
			Source: pulumi.String("http://192.168.0.72:8080/kvm-images/k3os-v21.raw"),
		}, pulumi.Provider(provider), pulumi.Protect(true))
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
					WaitForLease: pulumi.Bool(true),
				},
			},
			Disks: &libvirt.DomainDiskArray{
				&libvirt.DomainDiskArgs{
					VolumeId: volume.ID(),
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
