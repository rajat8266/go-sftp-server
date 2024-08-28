package gcs

var (
	CredentialFile string
	User           string
	BucketName     string
)

func SetConfigForGcs(user string, bucketName string, credentialFile string) {
	CredentialFile = credentialFile
	User = user
	BucketName = bucketName
}
