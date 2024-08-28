package main

import (
	"github.com/go-sftp-server/config"
	"github.com/go-sftp-server/server"
)

func main() {
	// Load the root configuration from the JSON file.
	rootConfig := config.LoadRootConfig()

	// Load and configure SSH settings, including private key and authorized keys.
	sshConfig := config.LoadSSHConfig()

	// Validate User, Configure SSH to set the bucket name for logged-in users.
	sshConfig = config.SetAccountForSSHConfig(sshConfig)

	// Start the SFTP server with the configured settings.
	server.Serve(sshConfig, rootConfig.Address, rootConfig.Port)
}
