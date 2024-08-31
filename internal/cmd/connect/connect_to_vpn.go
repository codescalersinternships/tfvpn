package connect

import (
	"fmt"
	"os/exec"

	"github.com/rs/zerolog/log"
)

func ConnectToVPN(ipAddr string) error {
	log.Info().Msg("connecting to the vpn server...")
	cmd := exec.Command("sshuttle", "-r", "root@"+ipAddr, "0.0.0.0/0")

	err := cmd.Start()
	if err != nil {
		return fmt.Errorf("failed to connect to the vpn server %w", err)
	}
	log.Info().Msg("connection established successfully!")
	return nil
}
