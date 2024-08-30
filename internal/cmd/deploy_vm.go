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

func deployVM(ctx context.Context, client *deployer.TFPluginClient, vm workloads.VM, network *workloads.ZNet) (workloads.Deployment, error) {
	dl := workloads.NewDeployment("dl_vpn", network.Nodes[0], network.SolutionType, nil, network.Name, nil, nil, []workloads.VM{vm}, nil)
	if err := dl.Validate(); err != nil {
		return workloads.Deployment{}, err
	}

	log.Info().Str("name", vm.Name).Uint32("node_id", dl.NodeID).Msg("deploying vpn server...")
	if err := client.DeploymentDeployer.Deploy(ctx, &dl); err != nil {
		return workloads.Deployment{}, err
	}

	deployedDl, err := client.State.LoadDeploymentFromGrid(ctx, dl.NodeID, dl.Name)
	if err != nil {
		return workloads.Deployment{}, errors.Wrapf(err, "failed to load vm from node %d", dl.NodeID)
	}

	return deployedDl, nil
}

func buildNetwork(nodeID uint32, projectName string) workloads.ZNet {
	return workloads.ZNet{
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

func rollback(ctx context.Context, client *deployer.TFPluginClient, dl *workloads.Deployment, net *workloads.ZNet, err error) error {
	log.Info().Msg("an error occurred while deploying, canceling all deployments")
	log.Info().Msg("canceling network...")
	if err := client.NetworkDeployer.Cancel(ctx, net); err != nil {
		return err
	}

	log.Info().Msg("canceling deployments..")
	if err := client.DeploymentDeployer.Cancel(ctx, dl); err != nil {
		return err
	}

	log.Info().Msg("deployment canceled successfully")
	return err
}
