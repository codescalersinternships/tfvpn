// Package filter provides functions for filtering available nodes based on the VPN configuration.
package filter

import (
	"context"
	"fmt"

	"github.com/codescalersinternships/tfvpn/pkg/config"
	"github.com/threefoldtech/tfgrid-sdk-go/grid-client/deployer"
	proxy_types "github.com/threefoldtech/tfgrid-sdk-go/grid-proxy/pkg/types"
)

// FilterNode filters the available nodes based on the vpn config
func FilterNode(ctx context.Context, client deployer.TFPluginClient, config config.VPNConfig) (uint32, error) {
	healthy := true
	freeIPs := uint64(1)
	filter := proxy_types.NodeFilter{
		Healthy: &healthy,
		Region:  &config.Region,
		City:    &config.City,
		Country: &config.Country,
		FreeIPs: &freeIPs,
	}

	nodes, _, err := client.GridProxyClient.Nodes(ctx, filter, proxy_types.Limit{})
	if err != nil {
		return 0, err
	}

	if len(nodes) == 0 || (len(nodes) == 1 && nodes[0].NodeID == 1) {
		return 0, fmt.Errorf("no available nodes")
	}

	nodeID := uint32(0)
	for _, node := range nodes {
		if node.NodeID != 1 {
			nodeID = uint32(node.NodeID)
			break
		}
	}

	return nodeID, nil
}
