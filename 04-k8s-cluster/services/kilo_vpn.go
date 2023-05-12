package services

import (
	"k8s-cluster/packages/kilo"
)

const (
	vpnPort     = 32200
	serviceName = "kilo-vpn"
)

func (s *Services) RunKiloVPN() (*kilo.StartedKilo, error) {
	s.KiloVPN.Name = serviceName
	s.KiloVPN.Port = vpnPort
	ns, err := kilo.CreateNS(s.ctx, "kilo-vpn")
	if err != nil {
		return nil, err
	}

	deployed, err := kilo.RunKilo(s.ctx, ns, s.KiloVPN)
	if err != nil {
		return nil, err
	}

	return deployed, nil
}
