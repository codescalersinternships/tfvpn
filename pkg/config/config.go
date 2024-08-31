// Package config for handling user configuration
package config

import (
	"fmt"

	"github.com/cosmos/go-bip39"
)

// Config struct that holds user configuration
type Config struct {
	Mnemonics string `json:"mnemonics"`
	Network   string `json:"network"`
}

// Validate validates the user configuration
func (c *Config) Validate() error {
	if c.Network == "" {
		c.Network = "dev"
	}
	if !(c.Network == "dev" || c.Network == "test" || c.Network == "qa" || c.Network == "main") {
		return fmt.Errorf("invalid network %s", c.Network)
	}

	if c.Mnemonics == "" {
		return fmt.Errorf("mnemonics is required")
	}
	if !bip39.IsMnemonicValid(c.Mnemonics) {
		return fmt.Errorf("entered mnemonics %s are invalid", c.Mnemonics)
	}

	return nil
}

// VPNConfig struct that holds VPN configuration
type VPNConfig struct {
	Country string
	City    string
	Region  string
}
