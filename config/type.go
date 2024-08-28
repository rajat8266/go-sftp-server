package config

// Config represents information go-sftp-server needs to operate
type Config struct {
	SftpAccounts           []SftpAccounts
	Address                string
	Port                   string
	KeyPathSSH             string
	SftpAuthorizedKeysFile string
	CredentialsFileGCS     string
}

// Account holds specific information for each account we support
type SftpAccounts struct {
	Username   string
	Password   string
	BucketName string
}
