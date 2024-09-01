package cmd

import (
	"github.com/codescalersinternships/tfvpn/pkg/config"
	"github.com/spf13/cobra"
)

func parseUpFlags(cmd *cobra.Command) (config.VPNConfig, error) {
	region, err := cmd.Flags().GetString("region")
	if err != nil {
		return config.VPNConfig{}, err
	}
	country, err := cmd.Flags().GetString("country")
	if err != nil {
		return config.VPNConfig{}, err
	}
	city, err := cmd.Flags().GetString("city")
	if err != nil {
		return config.VPNConfig{}, err
	}
	return config.VPNConfig{
		Region:  region,
		Country: country,
		City:    city,
	}, nil
}
