package hetzner

import (
	"gopkg.in/yaml.v3"
)

type CloudConfig struct {
	Users    []*UserCloudConfig
	FQDN     string
	GrowPart *GrowPartConfig
}

type UserCloudConfig struct {
	Name              string
	Sudo              string
	SSHAuthorizedKeys []string `yaml:"ssh_authorized_keys"`
}

type GrowPartConfig struct {
	Devices []string
}

func (c *CloudConfig) render() (string, error) {
	r := "#cloud-config\n"

	cfg, err := yaml.Marshal(&c)
	if err != nil {
		return "", err
	}

	return r + string(cfg), nil
}
