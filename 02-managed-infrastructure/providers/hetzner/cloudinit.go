package hetzner

import (
	"gopkg.in/yaml.v3"
)

type CloudConfig struct {
	SshPwauth bool `yaml:"ssh_pwauth"`
	Users    []*UserCloudConfig
	Hostname  string
	GrowPart *GrowPartConfig
}

type UserCloudConfig struct {
	Name              string
	Sudo              string
	SSHAuthorizedKeys []string `yaml:"ssh_authorized_keys"`
	LockPasswd        bool     `yaml:"lock_passwd"`
	Passwd            string
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
