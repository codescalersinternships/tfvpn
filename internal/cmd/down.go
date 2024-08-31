package cmd

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/threefoldtech/tfgrid-sdk-go/grid-client/deployer"
)

// Down stops the VPN server
func Down(ctx context.Context) error {
	client, ok := ctx.Value("client").(deployer.TFPluginClient)
	if !ok {
		return fmt.Errorf("failed to get grid client")
	}

	log.Info().Msg("disconnecting from vpn server")
	projectName := fmt.Sprintf("vm/%d/%s", client.TwinID, "vpn")

	err := killShuttle()
	if err != nil {
		return err
	}
	remoteIP, err := getRemoteIP(ctx, client, projectName)
	if err != nil {
		return err
	}
	err = destroyVM(client, projectName)
	if err != nil {
		return err
	}
	err = removeHost(remoteIP)
	if err != nil {
		return err
	}

	log.Info().Msg("disconnected from vpn server")
	return nil
}
