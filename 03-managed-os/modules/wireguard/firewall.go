package wireguard

import (
	"managed-os/modules/firewall/firewalld"
)

func GetRequiredFirewalldRule() *firewalld.PortRule {
	return &firewalld.PortRule{
		Name:     "auto-source-for-public:wireguard",
		Protocol: "udp",
		Port:     listenPort,
		Zone:     "public",
	}
}
