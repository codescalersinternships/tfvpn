package cmd

import (
	"context"
	"net"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/threefoldtech/tfgrid-sdk-go/grid-client/deployer"
	"github.com/threefoldtech/tfgrid-sdk-go/grid-client/workloads"
	"github.com/threefoldtech/zos/pkg/gridtypes"
)

// deployVM deploys a vm on the grid
func deployVM(ctx context.Context, client *deployer.TFPluginClient, vm workloads.VM, network *workloads.ZNet) (workloads.VM, error) {
	dl := workloads.NewDeployment("dl_vpn", network.Nodes[0], network.SolutionType, nil, network.Name, nil, nil, []workloads.VM{vm}, nil)
	if err := dl.Validate(); err != nil {
		return workloads.VM{}, err
	}

	log.Info().Str("name", vm.Name).Uint32("node_id", dl.NodeID).Msg("deploying vpn server...")
	if err := client.DeploymentDeployer.Deploy(ctx, &dl); err != nil {
		return workloads.VM{}, err
	}

	deployedVM, err := client.State.LoadVMFromGrid(ctx, dl.NodeID, vm.Name, dl.Name)
	if err != nil {
		return workloads.VM{}, errors.Wrapf(err, "failed to load vm from node %d", dl.NodeID)
	}

	return deployedVM, nil
}

func buildNetwork(nodeID uint32, projectName string) *workloads.ZNet {
	return &workloads.ZNet{
		Name: "vpn",
		IPRange: gridtypes.NewIPNet(net.IPNet{
			IP:   net.IPv4(10, 20, 0, 0),
			Mask: net.CIDRMask(16, 32),
		}),
		AddWGAccess:  false,
		Nodes:        []uint32{nodeID},
		SolutionType: projectName,
	}
}

func buildVM(networkName string) workloads.VM {
	return workloads.VM{
		Name:        "vpn",
		Flist:       "https://hub.grid.tf/tf-official-apps/threefoldtech-ubuntu-22.04.flist",
		Entrypoint:  "/sbin/zinit init",
		CPU:         1,
		Memory:      2048,
		RootfsSize:  15360,
		NetworkName: networkName,
		PublicIP:    true,
	}
}
