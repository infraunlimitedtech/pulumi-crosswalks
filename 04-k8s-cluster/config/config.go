package config

type CRDS struct {
	Install bool
	Path    string
}

type Firewalls struct {
	Hetzner   *Firewall
	Firewalld *Firewall
}

type Firewall struct {
	Managed bool
}

type HelmParams struct {
	Version string
}

type Status bool

func (s *Status) WithDefault(t bool) Status {
	if s == nil {
		return Status(t)
	}
	return *s
}
