package config

import (
	"fmt"
	"log"
	"os"

	"github.com/go-sftp-server/gcs"
	"golang.org/x/crypto/ssh"
)

// LoadSSHConfig initializes the SSH server configuration.
func LoadSSHConfig() *ssh.ServerConfig {
	configSSH := &ssh.ServerConfig{
		NoClientAuth:  false,
		ServerVersion: "SSH-2.0-GCS-SFTP",
		AuthLogCallback: func(conn ssh.ConnMetadata, method string, err error) {
			logAuthAttempt(conn, method, err)
		},
	}

	privateKey := loadPrivateKey(RootConfig.KeyPathSSH)
	configSSH.AddHostKey(privateKey)

	processPublicKeyAuth(configSSH, &RootConfig)

	return configSSH
}

// loadPrivateKey reads and parses the private key from the specified path.
func loadPrivateKey(path string) ssh.Signer {
	privateBytes, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("failed to load private key from %s: %s", path, err)
	}

	privateKey, err := ssh.ParsePrivateKey(privateBytes)
	if err != nil {
		log.Fatalf("failed to parse private key: %s", err)
	}

	return privateKey
}

// logAuthAttempt logs the result of an authentication attempt.
func logAuthAttempt(conn ssh.ConnMetadata, _ string, err error) {
	if err != nil {
		log.Printf("User: %s Authentication Attempt from %s", conn.User(), conn.RemoteAddr())
	} else {
		log.Printf("User: %s Authentication Accepted from %s", conn.User(), conn.RemoteAddr())
	}
}

// processPublicKeyAuth configures public key authentication for the SSH server.
func processPublicKeyAuth(configSSH *ssh.ServerConfig, rootConfig *Config) {
	if rootConfig.SftpAuthorizedKeysFile == "" {
		return
	}

	authorizedKeys := loadAuthorizedKeys(rootConfig.SftpAuthorizedKeysFile)

	configSSH.PublicKeyCallback = func(conn ssh.ConnMetadata, auth ssh.PublicKey) (*ssh.Permissions, error) {
		for _, pubKey := range authorizedKeys {
			if comparePublicKeys(pubKey, auth) {
				return &ssh.Permissions{
					Extensions: map[string]string{
						"pubkey-fp": ssh.FingerprintSHA256(auth),
					},
				}, nil
			}
		}
		return nil, fmt.Errorf("unknown public key for %q", conn.User())
	}
}

// loadAuthorizedKeys reads and parses the authorized keys file.
func loadAuthorizedKeys(path string) []ssh.PublicKey {
	authorizedKeysBytes, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("failed to load authorized keys file from %s: %s", path, err)
	}

	var authorizedKeys []ssh.PublicKey
	for len(authorizedKeysBytes) > 0 {
		pubKey, _, _, rest, err := ssh.ParseAuthorizedKey(authorizedKeysBytes)
		if err != nil {
			if err.Error() == "ssh: no key found" {
				break
			}
			log.Fatalf("failed to parse authorized key: %s", err)
		}
		authorizedKeys = append(authorizedKeys, pubKey)
		authorizedKeysBytes = rest
	}
	return authorizedKeys
}

// comparePublicKeys compares two SSH public keys.
func comparePublicKeys(key1, key2 ssh.PublicKey) bool {
	return string(key1.Marshal()) == string(key2.Marshal())
}

// SetAccountForSSHConfig configures password authentication for the SSH server.
func SetAccountForSSHConfig(sshConfig *ssh.ServerConfig) *ssh.ServerConfig {
	sshConfig.PasswordCallback = func(c ssh.ConnMetadata, pass []byte) (*ssh.Permissions, error) {
		msg, exists := isValidAccount(c.User(), string(pass))
		if !exists {
			return nil, fmt.Errorf("%s and rejected for %q", msg, c.User())
		}
		return nil, nil
	}
	return sshConfig
}

// isValidAccount checks if the provided username and password match any configured account.
func isValidAccount(username, password string) (string, bool) {
	for _, acc := range RootConfig.SftpAccounts {
		if acc.Username == username && acc.Password == password {
			gcs.SetConfigForGcs(acc.Username, acc.BucketName, RootConfig.CredentialsFileGCS)
			return "", true
		}
	}
	return "either user or password not matched", false
}
