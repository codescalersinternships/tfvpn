package cmd

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// NOT SUFFICIENTLY TESTED YET, will possibly be replaced by ssh package
func removeDuplicateHost(host string) error {
	knownHostsPath := filepath.Join(os.Getenv("HOME"), ".ssh", "known_hosts")
	tempFile, err := os.CreateTemp("", "known_hosts")
	if err != nil {
		return fmt.Errorf("failed to create temp file %w", err)
	}
	defer tempFile.Close()

	file, err := os.Open(knownHostsPath)
	if err != nil {
		return fmt.Errorf("failed to open known_hosts file %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if !strings.Contains(line, host) {
			_, err = tempFile.WriteString(line + "\n")
			if err != nil {
				return fmt.Errorf("failed to write to temp file %w", err)
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading known_hosts file %w", err)
	}

	if err := os.Rename(tempFile.Name(), knownHostsPath); err != nil {
		return fmt.Errorf("failed to replace known_hosts file %w", err)
	}

	return nil
}

func addHostToKnownHosts(host string) error {
	cmd := exec.Command("ssh-keyscan", "-H", host)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to fetch host key %w", err)
	}

	file, err := os.OpenFile(filepath.Join(os.Getenv("HOME"), ".ssh", "known_hosts"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open known_hosts file %w", err)
	}
	defer file.Close()

	_, err = file.Write(output)
	if err != nil {
		return fmt.Errorf("failed to write to known_hosts file %w", err)
	}
	return nil
}
