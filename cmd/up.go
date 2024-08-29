package cmd

import (
	command "github.com/codescalersinternships/tfvpn/internal/cmd"
	"github.com/codescalersinternships/tfvpn/pkg/config"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var upCmd = &cobra.Command{
	Use:   "up",
	Short: "deploy application on the grid",
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		region, err := cmd.Flags().GetString("region")
		if err != nil {
			log.Fatal().Err(err).Send()
		}
		country, err := cmd.Flags().GetString("country")
		if err != nil {
			log.Fatal().Err(err).Send()
		}
		city, err := cmd.Flags().GetString("city")
		if err != nil {
			log.Fatal().Err(err).Send()
		}

		vpnConfig := config.VPNConfig{
			Region:  region,
			Country: country,
			City:    city,
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
	// still to discuss
	// upCmd.PersistentFlags().BoolP("dns", "-d", false, "use dns of the vpn server")
}
