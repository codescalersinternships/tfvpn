package cmd

import (
	"context"
	"os/exec"
	"strconv"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/threefoldtech/tfgrid-sdk-go/grid-client/deployer"
	"github.com/threefoldtech/tfgrid-sdk-go/grid-client/graphql"
)

// killShuttle kills the shuttle process
func killShuttle() error {
	// TODO: store the pid of the process somewhere and kill it
	log.Info().Msg("killing sshuttle...")
	cmd := exec.Command("pkill", "sshuttle")
	_, err := cmd.CombinedOutput()
	if err != nil {
		log.Error().Msgf("failed to kill the shuttle process %v", err)
		return err
	}

	log.Info().Msg("killed the shuttle process")
	return nil
}

// destroyVM destroys the VM
func destroyVM(client deployer.TFPluginClient, projectName string) error {
	log.Info().Msg("canceling deployment...")
	err := client.CancelByProjectName(projectName)
	if err != nil {
		log.Error().Msgf("failed to cancel deployment %s", err)
		return err
	}

	log.Info().Msgf("deployment canceled successfully")
	return nil
}

// removeHost removes the host from the known hosts
func removeHost(remoteIP string) error {
	log.Info().Msgf("removing host %s from known hosts...", remoteIP)
	cmd := exec.Command("ssh-keygen", "-R", remoteIP)
	_, err := cmd.CombinedOutput()
	if err != nil {
		log.Error().Msgf("failed to remove the host %s", err)
		return err
	}

	log.Info().Msgf("host %s is removed", remoteIP)
	return nil
}

// getRemoteIP gets the remote IP
func getRemoteIP(ctx context.Context, client deployer.TFPluginClient, projectName string) (string, error) {
	nodeID, err := getNodeID(client, projectName)
	if err != nil {
		return "", err
	}

	deployment, err := client.State.LoadVMFromGrid(ctx, nodeID, "vpn", "dl_vpn")
	if err != nil {
		log.Error().Msgf("failed to get deployment %s", err)
		return "", err
	}

	return strings.Split(deployment.ComputedIP, "/")[0], nil
}

// getNodeID gets the node ID
func getNodeID(client deployer.TFPluginClient, projectName string) (uint32, error) {
	contracts, err := client.ContractsGetter.ListContractsOfProjectName(projectName)
	if err != nil {
		log.Error().Msgf("failed to list contracts %s", err)
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
			log.Error().Msgf("failed to convert contract.ContractID to uint64 %s", err)
			return err
		}
		contractIDs = append(contractIDs, contractID)
	}

	client.State.StoreContractIDs(nodeID, contractIDs...)

	return nil
}
