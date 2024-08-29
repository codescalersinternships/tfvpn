package cmd

import (
	"fmt"

	"github.com/codescalersinternships/tfvpn/pkg/config"
	"github.com/threefoldtech/tfgrid-sdk-go/grid-client/deployer"
)

// NewApp creates a new instance of the application
func setup(conf config.Config) (deployer.TFPluginClient, error) {
	client, err := deployer.NewTFPluginClient(conf.Mnemonics, deployer.WithNetwork(conf.Network))
	if err != nil {
		return deployer.TFPluginClient{}, fmt.Errorf("failed to load grid client %w", err)
	}

	return client, nil
}
