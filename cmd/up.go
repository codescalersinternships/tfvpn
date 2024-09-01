package cmd

import (
	command "github.com/codescalersinternships/tfvpn/internal/cmd"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var upCmd = &cobra.Command{
	Use:   "up",
	Short: "deploy application on the grid",
	Run: func(cmd *cobra.Command, args []string) {
		vpnConfig, err := parseUpFlags(cmd)
		if err != nil {
			log.Fatal().Err(err).Send()
		}
		if err := command.Up(cmd.Context(), vpnConfig); err != nil {
			log.Fatal().Err(err).Send()
		}
	},
}

func init() {
	upCmd.PersistentFlags().StringP("region", "", "", "region to deploy the vpn server on it")
	upCmd.PersistentFlags().StringP("country", "", "", "country to deploy the vpn server on it")
	upCmd.PersistentFlags().StringP("city", "", "", "city to deploy the vpn server on it")
}
