## SFTP Server with Google Cloud Bucket as File Storage
This project implements an SFTP server that uses Google Cloud Storage (GCS) as the backend storage system. It provides a secure, scalable, and configurable SFTP service that can be used for file transfers, leveraging GCS for data storage and management.

### Sftp Backend
**Overview**: The SFTP server is built using the Go programming language, with support for secure file transfer over SSH. It handles file operations such as reading, writing, listing, renaming, and deleting files in a GCS bucket.
**Key Components**:
- **Connection Handling**: Manages incoming SSH connections using `golang.org/x/crypto`, allowing for secure file transfer sessions.
- **SFTP Handlers**: Custom handlers for SFTP operations are implemented, interfacing with GCS for all file-related operations using `github.com/pkg/sftp`.
- **Session Management**: Supports multiple concurrent SFTP sessions, with each session linked to specific GCS buckets based on user authentication.

### Google Cloud Backend
**Overview**: The project uses Google Cloud Storage as the primary data store, where all files uploaded via SFTP are stored and managed.
**Key Components**:
- **GCS Client Setup**: A GCS client is initialized using service account credentials, allowing the application to interact with GCS buckets using `cloud.google.com/go/storage`.
- **Bucket and Object Operations**: The project includes functions for checking if a bucket exists, creating directories, and performing CRUD operations on objects within GCS.
- **User-Specific Buckets**: Each user is associated with a specific GCS bucket, ensuring data isolation and security.


### Configurations
**Overview**: The project is highly configurable, allowing administrators to define key settings via a JSON configuration file and environment variables.
**Key Components**:
- **Root Configuration**: Loads and parses a JSON configuration file that specifies settings such as the address and port for the SFTP server, GCS credentials, and user accounts.
- **SSH Server Configuration**: Configures the SSH server with the necessary settings, including private keys and authorized public keys for user authentication.
- **User and Bucket Mapping**: The configuration allows mapping of SFTP users to specific GCS buckets, ensuring that each user's files are stored in their designated bucket.

### Run
```./go-sftp-server --config-path path_to_config.json```
Generate a new config using `config/config.md`.

