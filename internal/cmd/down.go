package cmd

import (
	"context"
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/threefoldtech/tfgrid-sdk-go/grid-client/deployer"
	"github.com/threefoldtech/tfgrid-sdk-go/grid-client/graphql"
)

// Down stops the VPN server
func Down(ctx context.Context) error {
	client, ok := ctx.Value("client").(deployer.TFPluginClient)
	if !ok {
		return fmt.Errorf("failed to get grid client")
	}

	name := fmt.Sprintf("%d/%s", client.TwinID, "vpn")

	remoteIP, err := getRemoteIP(ctx, client, name)
	if err != nil {
		return err
	}

	err = killShuttle()
	if err != nil {
		return err
	}

	err = destroyVm(client, name)
	if err != nil {
		return err
	}

	err = removeHost(remoteIP)
	if err != nil {
		return err
	}

	return nil
}

// killShuttle kills the shuttle process
func killShuttle() error {
	// TODO: store the pid of the process somewhere and kill it
	cmd := exec.Command("pkill", "sshuttle")
	_, err := cmd.CombinedOutput()
	if err != nil {
		log.Error().Msgf("failed to kill the shuttle process: %s", err)
		return err
	}

	log.Info().Msg("killed the shuttle process")
	return nil
}

// destroyVm destroys the VM
func destroyVm(client deployer.TFPluginClient, name string) error {
	err := client.CancelByProjectName(name, true)
	if err != nil {
		log.Error().Msgf("failed to cancel deployment: %s", err)
		return err
	}

	log.Info().Msgf("deployment canceled successfully")
	return nil
}

// removeHost removes the host from the known hosts
func removeHost(remoteIP string) error {
	cmd := exec.Command("ssh-keygen", "-R", remoteIP)
	_, err := cmd.CombinedOutput()
	if err != nil {
		log.Error().Msgf("failed to remove the host: %s", err)
		return err
	}

	log.Info().Msgf("removed host %s from known hosts", remoteIP)
	return nil
}

// getRemoteIP gets the remote IP
func getRemoteIP(ctx context.Context, client deployer.TFPluginClient, name string) (string, error) {

	nodeID, err := getNodeID(client, name)
	if err != nil {
		return "", err
	}

	deployment, err := client.State.LoadVMFromGrid(ctx, nodeID, "vpn", "dl_vpn")
	if err != nil {
		log.Error().Msgf("failed to get deployment: %s", err)
		return "", err
	}

	ip := strings.Split(deployment.ComputedIP, "/")[0]

	return ip, nil
}

// getNodeID gets the node ID
func getNodeID(client deployer.TFPluginClient, name string) (uint32, error) {
	contracts, err := client.ContractsGetter.ListContractsOfProjectName(name)
	if err != nil {
		log.Error().Msgf("failed to list contracts: %s", err)
		return 0, err
	}

	if len(contracts.NodeContracts) == 0 {
		log.Error().Msgf("no vpn server found")
		return 0, err
	}

	nodeID := contracts.NodeContracts[0].NodeID

	err = appendToState(client, contracts, nodeID)
	if err != nil {
		return 0, err
	}

	return nodeID, nil
}

// appendToState appends the contract to the state
func appendToState(client deployer.TFPluginClient, contracts graphql.Contracts, nodeID uint32) error {
	var contractIDs []uint64

	for _, contract := range contracts.NodeContracts {
		contractID, err := strconv.ParseUint(contract.ContractID, 10, 64)
		if err != nil {
			log.Error().Msgf("failed to convert contract.ContractID to uint64: %s", err)
			return err
		}
		contractIDs = append(contractIDs, contractID)
	}

	client.State.StoreContractIDs(nodeID, contractIDs...)

	return nil
}
