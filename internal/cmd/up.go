package cmd

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/codescalersinternships/tfvpn/internal/cmd/connect"
	"github.com/codescalersinternships/tfvpn/internal/filter"
	"github.com/codescalersinternships/tfvpn/pkg/config"
	"github.com/rs/zerolog/log"
	"github.com/threefoldtech/tfgrid-sdk-go/grid-client/deployer"
)

// Up deploys the vpn server on the grid and connects to it
func Up(ctx context.Context, vpnCfg config.VPNConfig) error {
	if err := installRequirements(); err != nil {
		return err
	}
	//nolint:staticcheck
	client, ok := ctx.Value("client").(deployer.TFPluginClient)
	if !ok {
		return fmt.Errorf("failed to get grid client")
	}
	nodeID, err := filter.FilterNode(ctx, client, vpnCfg)
	if err != nil {
		return err
	}

	network := buildNetwork(nodeID, fmt.Sprintf("%d/%s", client.TwinID, "vpn"))
	if err := client.NetworkDeployer.Deploy(ctx, &network); err != nil {
		return fmt.Errorf("failed to deploy network %w", err)
	}
	pubKey, err := getPublicKey()
	if err != nil {
		return err
	}
	vm := buildVM(network.Name, pubKey)
	dl, err := deployVM(ctx, &client, vm, &network)
	if err != nil {
		if err := rollback(ctx, &client, &dl, &network, err); err != nil {
			return err
		}
		return fmt.Errorf("failed to deploy vm %w", err)
	}
	deployedVM := dl.Vms[0]
	ipAddr := strings.Split(deployedVM.ComputedIP, "/")[0]
	log.Info().Str("public_ip", ipAddr).Msg("vpn server deployed successfully!")

	if err := connect.ConnectToVM(30*time.Second, "root", ipAddr); err != nil {
		if err := rollback(ctx, &client, &dl, &network, err); err != nil {
			return err
		}
		return err
	}
	if err := connect.ConnectToVPN(ipAddr); err != nil {
		if err := rollback(ctx, &client, &dl, &network, err); err != nil {
			return err
		}
		return err
	}

	return nil
}
