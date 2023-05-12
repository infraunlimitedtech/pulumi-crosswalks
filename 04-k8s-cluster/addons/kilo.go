package addons

import (
	"k8s-cluster/packages/kilo"
)

const (
	vpnPort     = 31200
	serviceName = "kilo"
)

func (a *Addons) RunKilo() (*kilo.StartedKilo, error) {
	a.Kilo.Name = serviceName
	a.Kilo.Port = vpnPort

	deployed, err := kilo.RunKilo(a.ctx, a.Namespace, a.Kilo)
	if err != nil {
		return nil, err
	}

	return deployed, nil
}
