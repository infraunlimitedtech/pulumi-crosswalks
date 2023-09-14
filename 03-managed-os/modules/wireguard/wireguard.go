package wireguard

import (
	"fmt"
	"managed-os/config"
	"managed-os/utils"
	"net"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/spigell/pulumi-file/sdk/go/file"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
	"inet.af/netaddr"
)

const (
	listenPort = 51822
)

type Cluster struct {
	Ctx                *pulumi.Context
	Iface              string
	Info               pulumi.AnyOutput
	Nodes              []*config.Node
	InfraLayerNodeInfo pulumi.AnyOutput
}

type CreatedCluster struct {
	Peers        pulumi.AnyOutput
	MasterConfig pulumi.StringOutput
	ListenPort   int
	Resources    map[string]pulumi.Resource
}

func GetRequiredPkgs() []string {
	return []string{"wireguard-tools"}
}

func (c *Cluster) Manage(deps []map[string]pulumi.Resource) (*CreatedCluster, error) {
	resources := make(map[string]pulumi.Resource)
	c.Nodes = append(c.Nodes, &config.Node{
		ID: "mgmt",
		Wireguard: config.Wireguard{
			MgmtNode: true,
			IP:       "10.10.1.1",
		},
	})
	wgPeers := buildWgPeers(c.Nodes, c.Info, c.InfraLayerNodeInfo)

	done := &CreatedCluster{}

	for _, node := range c.Nodes {
		if node.Wireguard.MgmtNode {
			done.MasterConfig = generateWgConfig(wgPeers, node)
		} else {
			deployed, err := file.NewRemote(c.Ctx, fmt.Sprintf("%s-WGCluster", node.ID), &file.RemoteArgs{
				Connection: &file.ConnectionArgs{
					Address:    pulumi.Sprintf("%s:22", utils.ExtractValueFromPulumiMapMap(c.InfraLayerNodeInfo, node.ID, "ip")),
					User:       utils.ExtractValueFromPulumiMapMap(c.InfraLayerNodeInfo, node.ID, "user"),
					PrivateKey: utils.ExtractValueFromPulumiMapMap(c.InfraLayerNodeInfo, node.ID, "key"),
				},
				Hooks: &file.HooksArgs{
					CommandAfterCreate: pulumi.String("sudo systemctl enable wg-quick@kubewg0 && sudo systemctl restart wg-quick@kubewg0"),
					CommandAfterUpdate: pulumi.String("sudo systemctl restart wg-quick@kubewg0"),
				},
				UseSudo: pulumi.Bool(true),
				Path:    pulumi.String("/etc/wireguard/kubewg0.conf"),
				Content: generateWgConfig(wgPeers, node),
			}, pulumi.DependsOn(utils.ConvertMapSliceToSliceByKey(deps, node.ID)),
				pulumi.RetainOnDelete(true),
				pulumi.Timeouts(&pulumi.CustomTimeouts{Create: "5m"}),
			)
			if err != nil {
				return nil, err
			}

			resources[node.ID] = deployed
		}
	}

	done.ListenPort = listenPort
	done.Peers = wgPeers
	done.Resources = resources

	return done, nil
}

func buildWgPeers(nodes []*config.Node, wgInfo pulumi.AnyOutput, infraNodesInfo pulumi.AnyOutput) pulumi.AnyOutput {
	return pulumi.ToSecret(pulumi.All(wgInfo, infraNodesInfo).ApplyT(func(args []interface{}) []Peer {
		peers := make([]Peer, 0)

		info, ok := args[0].(map[string]interface{})
		if !ok {
			info = make(map[string]interface{})
		}

		key, pub, ip := "", "", ""
		m := make(map[string]netaddr.IP)

		for _, node := range nodes {
			if info[node.ID] != nil && node.Wireguard.IP == "" {
				k := info[node.ID].(map[string]interface{})
				ip = k["ip"].(string)
				m[node.ID] = netaddr.MustParseIP(ip)
			}
			if node.Wireguard.IP != "" {
				ip = node.Wireguard.IP
				m[node.ID] = netaddr.MustParseIP(ip)
			}
		}

		for _, node := range nodes {
			if info[node.ID] == nil {
				generated, _ := wgtypes.GeneratePrivateKey()
				key = generated.String()
				pub = generated.PublicKey().String()
			} else {
				k := info[node.ID].(map[string]interface{})
				key = k["privatekey"].(string)
				pub = k["publickey"].(string)
			}

			if m[node.ID].IsZero() {
				ipo, _, err := net.ParseCIDR(node.Wireguard.CIDR)
				if err != nil {
					panic(fmt.Sprintf("Can not parse CIDR for Wireguard! Is it a valid network? (%s)", err.Error()))
				}

				start := netaddr.MustParseIP(ipo.String())
				i := start.Next()
				for !free(m, i) {
					i = i.Next()
				}
				m[node.ID] = i
			}

			peer := Peer{
				ID:          node.ID,
				PrivateKey:  key,
				PublicKey:   pub,
				PrivateAddr: m[node.ID].String(),
			}

			if !node.Wireguard.MgmtNode {
				infraNodeInfo := args[1].(map[string]interface{})[node.ID].(map[string]interface{})
				peer.PublicAddr = infraNodeInfo["ip"].(string)
			}

			peers = append(peers, peer)
		}
		return peers
	})).(pulumi.AnyOutput)
}

func generateWgConfig(wgPeersOutput pulumi.AnyOutput, self *config.Node) pulumi.StringOutput {
	return pulumi.ToSecret(pulumi.Unsecret(pulumi.All(wgPeersOutput).ApplyT(func(args []interface{}) string {
		wgPeers := args[0].([]Peer)
		if len(self.Wireguard.AdditionalPeers) > 0 {
			for _, p := range self.Wireguard.AdditionalPeers {
				additionalPeer := Peer{
					PublicKey:  p.PublicKey,
					Endpoint:   p.Endpoint,
					AllowedIps: p.AllowedIps,
				}
				wgPeers = append(wgPeers, additionalPeer)
			}
		}

		peersWithoutSelf := ToPeers(wgPeers).without(self.ID)

		for k, v := range peersWithoutSelf {
			peersWithoutSelf[k].PersistentKeepalive = 25
			if len(peersWithoutSelf[k].AllowedIps) == 0 {
				peersWithoutSelf[k].AllowedIps = []string{fmt.Sprintf("%s/32", v.PrivateAddr)}
			}
			if v.PublicAddr != "" {
				peersWithoutSelf[k].Endpoint = fmt.Sprintf("%s:%d", v.PublicAddr, listenPort)
			}
		}

		selfPeer := ToPeers(wgPeers).Get(self.ID)

		config := &WgConfig{
			Peer: peersWithoutSelf.getWgPeers(),
			Interface: WgInterface{
				Address:    selfPeer.PrivateAddr,
				PrivateKey: selfPeer.PrivateKey,
				ListenPort: listenPort,
			},
		}

		wgConfig, err := renderConfig(config)
		if err != nil {
			panic(fmt.Sprintf("Error while render Wireguard config %e", err))
		}
		return wgConfig
	}))).(pulumi.StringOutput)
}

func (w *CreatedCluster) ConvertPeersToMapMap() pulumi.StringMapMapOutput {
	return w.Peers.ApplyT(func(v interface{}) map[string]map[string]string {
		m := make(map[string]map[string]string)
		for _, peer := range v.([]Peer) {
			m[peer.ID] = make(map[string]string)
			m[peer.ID]["ip"] = peer.PrivateAddr
			m[peer.ID]["privatekey"] = peer.PrivateKey
			m[peer.ID]["publickey"] = peer.PublicKey
		}
		return m
	}).(pulumi.StringMapMapOutput)
}

func free(m map[string]netaddr.IP, match netaddr.IP) bool {
	for _, n := range m {
		if n == match {
			return false
		}
	}
	return true
}
