## Configuration guide for go-sftp-server
This document provides an overview of the configuration settings required for the go-sftp-server application. These settings are defined in a JSON configuration file and are loaded at runtime.

### Configuration Structure
The main configuration file is parsed into a Config struct, which holds various settings for the application, including paths to SSH keys, SFTP accounts, and Google Cloud Storage (GCS) credentials.

Example Configuration File:

```json
{
  "SftpAccounts": [
    {
      "Username": "user1",
      "Password": "password1",
      "BucketName": "bucket-name-1"
    },
    {
      "Username": "user2",
      "Password": "password2",
      "BucketName": "bucket-name-2"
    }
  ],
  "Address": "127.0.0.1",
  "Port" : "2024",
  "KeyPathSSH": "path_to_private_key",
  "SftpAuthorizedKeysFile": "",
  "CredentialsFileGCS": "path_to_gcs_cred_file"
}
```

### Config Fields
- **KeyPathSSH**:
Path to the SSH private key used by the SFTP server for secure connections.
Generate using `ssh-keygen -b 2048 -t rsa -f filepath.txt`

- **CredentialsFileGCS**:
Path to the Google Cloud Storage credentials JSON file. This file is used to authenticate the application to access GCS buckets.
`generate key of service account with access to buckets read write or create a service account with access to required buckets`

- **SftpAuthorizedKeysFile**:
Path to the file containing authorized SSH public keys for users who are allowed to connect. If this field is not set, only password authentication will be used.

- **SftpAccounts**:
List of SFTP accounts, where each account contains the following fields:

  - **Username**:
  The username for the SFTP account.

  - **Password**:
  The password for the SFTP account. This is used for password-based authentication.

  - **BucketName**:
  The name of the Google Cloud Storage bucket that the user will be granted access to.



