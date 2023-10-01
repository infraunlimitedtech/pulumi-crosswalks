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
