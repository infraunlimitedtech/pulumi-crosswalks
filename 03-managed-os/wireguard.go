package main

import (
	"fmt"
	"managed-os/utils/wireguard"
)

type Wireguard struct {
	PrivateAddr     string
	PublicKey       string
	PrivateKey      string
	AdditionalPeers []AdditionalPeer `json:"additional_peers" yaml:"additional_peers"`
}

type AdditionalPeer struct {
	AllowedIps []string `json:"allowed_ips" yaml:"allowed_ips"`
	PublicKey  string
}

func buildWgPeers(nodes []Node, self Node) []wireguard.Peer {
	peers := make([]wireguard.Peer, 0)
	for _, node := range nodes {
		peer := wireguard.Peer{
			ID:          node.ID,
			PrivateKey:  node.Wireguard.PrivateKey,
			PrivateAddr: node.Wireguard.PrivateAddr,
			PublicKey:   node.Wireguard.PublicKey,
			PublicAddr:  node.PublicIP,
		}
		peers = append(peers, peer)
	}
	if len(self.Wireguard.AdditionalPeers) > 0 {
		for _, p := range self.Wireguard.AdditionalPeers {
			additionalPeer := wireguard.Peer{
				PublicKey:  p.PublicKey,
				AllowedIps: p.AllowedIps,
			}
			peers = append(peers, additionalPeer)
		}
	}
	return peers
}

func renderWgConfig(peers []wireguard.Peer, self Node) (string, error) {
	peersWithoutSelf := wireguard.ToPeers(peers).Without(self.ID)

	for k, v := range peersWithoutSelf {
		peersWithoutSelf[k].PersistentKeepalive = 25
		if len(peersWithoutSelf[k].AllowedIps) == 0 {
			peersWithoutSelf[k].AllowedIps = []string{fmt.Sprintf("%s/32", v.PrivateAddr)}
		}
		if v.PublicAddr != "" {
			peersWithoutSelf[k].Endpoint = fmt.Sprintf("%s:%d", v.PublicAddr, 51820)
		}
	}

	config := &wireguard.WgConfig{
		Peer: peersWithoutSelf.GetWgPeers(),
		Interface: wireguard.WgInterface{
			Address:    self.Wireguard.PrivateAddr,
			PrivateKey: self.Wireguard.PrivateKey,
			ListenPort: wgListenPort,
		},
	}

	wgConfig, err := wireguard.RenderConfig(config)
	if err != nil {
		return "", err
	}

	return wgConfig, nil
}
