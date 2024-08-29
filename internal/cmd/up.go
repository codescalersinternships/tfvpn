package cmd

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/codescalersinternships/tfvpn/internal/filter"
	"github.com/codescalersinternships/tfvpn/pkg/config"
	"github.com/rs/zerolog/log"
	"github.com/threefoldtech/tfgrid-sdk-go/grid-client/deployer"
)

// still fails to connect to the vpn server, needs to be fixed
func Up(ctx context.Context, vpnCfg config.VPNConfig) error {
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
	if err := client.NetworkDeployer.Deploy(ctx, network); err != nil {
		return fmt.Errorf("failed to deploy network %w", err)
	}

	vm := buildVM(network.Name)
	deployedVM, err := deployVM(ctx, &client, vm, network)
	if err != nil {
		log.Info().Msg("an error occured while deploying the vm, cancelling...")
		if err := client.NetworkDeployer.Cancel(ctx, network); err != nil {
			return fmt.Errorf("failed to cancel network %w", err)
		}
		return fmt.Errorf("failed to deploy vm %w", err)
	}
	log.Info().Str("public_ip", deployedVM.ComputedIP).Msg("vpn server deployed successfully")

	fmt.Println("enter your local password to connect to the vpn server:")
	r := bufio.NewReader(os.Stdin)
	password, _ := r.ReadString('\n')
	password = strings.TrimSpace(password)
	log.Info().Msg("connecting to the vpn server...")
	time.Sleep(30 * time.Second)

	host := strings.Split(deployedVM.ComputedIP, "/")[0]
	if err := removeDuplicateHost(host); err != nil {
		return err
	}

	if err := addHostToKnownHosts(host); err != nil {
		return err
	}

	cmd := exec.Command("sshpass", "-p", password, "sshuttle", "-r", "root@"+host, "0.0.0.0/0")

	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Error().Msgf("failed to connect to the vpn server %s output %s", err, output)
		return fmt.Errorf("failed to connect to the vpn server %w", err)
	}
	log.Info().Msg("connection established successfully!")

	return nil
}
