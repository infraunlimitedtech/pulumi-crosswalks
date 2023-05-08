package firewalld

import (
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"fmt"
	"github.com/pulumi/pulumi-command/sdk/go/command/remote"
)


type Rule interface{}

type PortRule struct {
	Name string
	Protocol    string
	Port        int
	Zone        string
}

type InterfaceRule struct {
	Name string
	Interface string
	Zone        string
}

type SourceRule struct {
	Name string
	Source string
	Zone        string
}

type ServiceRule struct {
	Name string
	Service string
	Zone        string
}

var FirewalldSSHRule = &PortRule{
	Name: "auto-port-to-internal:ssh",
	Zone: "internal",
	Port: 22,
	Protocol: "tcp",
}

func (f *Firewalld) AddRules(rules []Rule) error {
	for _, rule := range rules {
		r, err := f.addRule(rule)
		if err != nil {
			return err
		}
		f.DependsOn = append(f.DependsOn, r)
	}
	return nil
}

func GetFirewallDSSHRule() Rule {
	return FirewalldSSHRule
}

func (f *Firewalld) addRule(rule Rule) (pulumi.Resource, error) {
	cmdCreate, cmdDelete, name := "", "", ""
	// We can't reload firewalld after all rules if only deleted rules are presents in diff.
	// Reload after delete is a workaround for this problem.
	reloadAfterDelete := "sudo firewall-cmd --reload"
	switch r := rule.(type) {

	case *PortRule:
		cmdCreate = fmt.Sprintf("sudo firewall-cmd --permanent --add-port=%d/%s --zone=%s", r.Port, r.Protocol, r.Zone)
		cmdDelete = fmt.Sprintf("sudo firewall-cmd --permanent --remove-port=%d/%s --zone=%s && %s", r.Port, r.Protocol, r.Zone, reloadAfterDelete)
		name = r.Name

	case *InterfaceRule:
		cmdCreate = fmt.Sprintf("sudo firewall-cmd --permanent --add-interface=%s --zone=%s", r.Interface, r.Zone)
		cmdDelete = fmt.Sprintf("sudo firewall-cmd --permanent --remove-interface=%s --zone=%s && %s", r.Interface, r.Zone, reloadAfterDelete)
		name = r.Name

	case *SourceRule:
		cmdCreate = fmt.Sprintf("sudo firewall-cmd --permanent --add-source=%s --zone=%s", r.Source, r.Zone)
		cmdDelete = fmt.Sprintf("sudo firewall-cmd --permanent --remove-source=%s --zone=%s && %s", r.Source, r.Zone, reloadAfterDelete)
		name = r.Name

	default:
		return nil, fmt.Errorf("unknown rule type: %T with content %+v", r, r)
	}


	applied, err := remote.NewCommand(f.ctx, fmt.Sprintf("%s-ApplyFirewallRule-%s", f.NodeInfo.ID, name), &remote.CommandArgs{
		Connection: &remote.ConnectionArgs{
			Host:       f.NodeInfo.Host,
			User:       f.NodeInfo.User,
			PrivateKey: f.NodeInfo.PrivateKey,
		},

		Create: pulumi.String(cmdCreate),
		Delete: pulumi.String(cmdDelete),
	}, pulumi.DeleteBeforeReplace(true), pulumi.DependsOn(f.DependsOn))
	if err != nil {
		err = fmt.Errorf("error while adding rules to firewalld: %w", err)
		return nil, err
	}

	return applied, nil
}

func (f *Firewalld) RemoveService(service *ServiceRule) error {
	cmd := fmt.Sprintf("sudo firewall-cmd --permanent --remove-service=%s --zone=%s", service.Service, service.Zone)
	removed, err := remote.NewCommand(f.ctx, fmt.Sprintf("%s-RemoveService-%s", f.NodeInfo.ID, service.Name), &remote.CommandArgs{
		Connection: &remote.ConnectionArgs{
			Host:       f.NodeInfo.Host,
			User:       f.NodeInfo.User,
			PrivateKey: f.NodeInfo.PrivateKey,
		},

		Create: pulumi.String(cmd),
	}, pulumi.DependsOn(f.DependsOn))

	f.DependsOn = append(f.DependsOn, removed)

	if err != nil {
		return err
	}

	return nil
}