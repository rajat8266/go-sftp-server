package gcs

import (
	"context"
	"fmt"
	"log"
	"os"

	"cloud.google.com/go/storage"
	"google.golang.org/api/iterator"
)

// Gcs is a struct representing Google Cloud Storage client.
// Add relevant fields in the struct definition as needed.

// GetFileFromGCS retrieves a file from Google Cloud Storage.
// It returns a storage.Reader to read the file content.
func (fs *Gcs) GetFileFromGCS(ctx context.Context, filepath string) (*storage.Reader, error) {
	object := fs.Bucket.Object(filepath)

	reader, err := object.NewReader(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create reader for file %s: %w", filepath, err)
	}

	log.Printf("User: %s Read: File name: %s", User, filepath)
	return reader, nil
}

// WriteFileOnGCS writes a file to Google Cloud Storage.
// It returns a storage.Writer to write the file content.
func (fs *Gcs) WriteFileOnGCS(ctx context.Context, filepath string) (*storage.Writer, error) {
	if err := fs.CreateDirectoryTreeBySftpRequest(ctx, filepath); err != nil {
		return nil, fmt.Errorf("failed to create directory tree for file %s: %w", filepath, err)
	}

	object := fs.GetStorageHandlerByName(filepath)
	log.Printf("User: %s Write: File name: %s", User, filepath)

	return object.NewWriter(ctx), nil
}

// RenameFileOnGCS renames a file on Google Cloud Storage by copying
// the original file to a new location and deleting the old file.
func (fs *Gcs) RenameFileOnGCS(ctx context.Context, srcFilepath, dstFilepath string) error {
	src := fs.GetStorageHandlerByName(srcFilepath)
	dst := fs.GetStorageHandlerByName(dstFilepath)

	// Copy the object to the new location
	if _, err := dst.CopierFrom(src).Run(ctx); err != nil {
		return fmt.Errorf("failed to copy object from %s to %s: %w", srcFilepath, dstFilepath, err)
	}

	// Delete the original object
	if err := src.Delete(ctx); err != nil {
		return fmt.Errorf("failed to delete original object %s: %w", srcFilepath, err)
	}

	log.Printf("User: %s Renamed file from %s to %s", User, srcFilepath, dstFilepath)
	return nil
}

// RemoveFileOnGCS deletes a file from Google Cloud Storage.
func (fs *Gcs) RemoveFileOnGcs(ctx context.Context, filepath string) error {
	object := fs.GetStorageHandlerByName(filepath)
	if err := object.Delete(ctx); err != nil {
		return fmt.Errorf("failed to delete file %s: %w", filepath, err)
	}

	log.Printf("User: %s Removed file: %s", User, filepath)
	return nil
}

// MakeDirOnGCS creates a directory in Google Cloud Storage by creating
// a zero-byte object with a trailing slash in the name.
func (fs *Gcs) MakeDirOnGCS(ctx context.Context, dirpath string) error {
	object := fs.GetStorageHandlerByName(dirpath)
	writer := object.NewWriter(ctx)

	if err := writer.Close(); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", dirpath, err)
	}

	log.Printf("User: %s Created directory: %s", User, dirpath)
	return nil
}

// ListFileOnGCS lists files in a specified directory in Google Cloud Storage.
// It returns a list of file information in the directory.
func (fs *Gcs) ListFileOnGCS(ctx context.Context, prefix string) (ListerAt, error) {
	if prefix != "" {
		prefix += "/"
	}

	var list []os.FileInfo
	objects := fs.Bucket.Objects(ctx, &storage.Query{
		Delimiter: "/",
		Prefix:    prefix,
	})

	for {
		objAttrs, err := objects.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Printf("User: %s Error iterating directory %s: %v", User, prefix, err)
			return nil, fmt.Errorf("failed to list files in directory %s: %w", prefix, err)
		}

		// Skip the prefix itself
		if (prefix != "" && objAttrs.Prefix == prefix) || (objAttrs.Prefix == "" && objAttrs.Name == prefix) {
			continue
		}

		list = append(list, &GcsFileInfo{
			prefix:  prefix,
			objAttr: objAttrs,
		})
	}

	log.Printf("User: %s Listing files in directory: /%s", User, prefix)
	return ListerAt(list), nil
}

// StatsOnGCS retrieves metadata for a specific file in Google Cloud Storage.
func (fs *Gcs) StatsOnGCS(ctx context.Context, filepath string) (ListerAt, error) {
	if filepath == "/" {
		file := &GcsFileInfo{
			objAttr: &storage.ObjectAttrs{
				Prefix: "/",
			},
		}
		return ListerAt([]os.FileInfo{file}), nil
	}

	object := fs.Bucket.Object(filepath)

	log.Printf("User: %s Getting file info for %s", User, filepath)

	attrs, err := object.Attrs(ctx)
	if err != nil {
		log.Printf("User: %s failed getting attribute for object with err :%s", User, err)
	}

	if err == storage.ErrObjectNotExist {
		object := fs.Bucket.Object(filepath + "/")

		log.Printf("User: %s retrying file info for %s", User, filepath+"/")

		attrs, err = object.Attrs(ctx)
	}
	if err != nil {
		log.Printf("User: %s failed to get file info for %s: %s", User, filepath, err)
		if err == storage.ErrObjectNotExist {
			err = os.ErrNotExist
		}
		return nil, err
	}

	file := &GcsFileInfo{
		objAttr: attrs,
	}
	return ListerAt([]os.FileInfo{file}), nil
}
