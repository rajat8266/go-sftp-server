package gcs

import (
	"context"
	"log"
	"path/filepath"

	"cloud.google.com/go/storage"
)

// CreateDirectoryTreeBySftpRequest ensures that the directory structure exists in Google Cloud Storage
// by creating directories as needed.
func (fs *Gcs) CreateDirectoryTreeBySftpRequest(ctx context.Context, filepath string) error {
	dirTree := fs.getDirectoryTree(filepath)

	for _, dir := range dirTree {
		dirWithSlash := dir + "/"
		object := fs.Bucket.Object(dirWithSlash)

		log.Printf("User: %s Checking if directory exists: /%s", User, dirWithSlash)

		_, err := object.Attrs(ctx)
		if err == storage.ErrObjectNotExist {
			log.Printf("User: %s Creating directory: /%s", User, dirWithSlash)

			writer := object.NewWriter(ctx)
			if err := writer.Close(); err != nil {
				return err
			}
		} else if err != nil {
			log.Printf("User: %s Failed to get object info for %s: %v", User, filepath, err)
			return err
		}
	}

	return nil
}

// getDirectoryTree generates a list of directories that need to be created for a given file path.
func (fs *Gcs) getDirectoryTree(dir string) []string {
	var dirsToCreate []string

	for {
		dir, _ = filepath.Split(dir)
		if dir == "" || dir == "/" {
			break
		}
		dir = filepath.Clean(dir)
		dirsToCreate = append([]string{dir}, dirsToCreate...)
	}

	return dirsToCreate
}

// GetStorageHandlerByName returns a Google Cloud Storage object handle for a given object name.
func (fs *Gcs) GetStorageHandlerByName(name string) *storage.ObjectHandle {
	return fs.Bucket.Object(name)
}
