package handler

import (
	"context"
	"fmt"
	"io"

	"github.com/go-sftp-server/gcs"
	"github.com/pkg/sftp"
)

// SftpHandler handles SFTP operations using Google Cloud Storage as the backend.
type SftpHandler struct {
	storage *gcs.Gcs
}

// Handler initializes the SFTP handlers for different file operations.
func Handler(ctx context.Context) (*sftp.Handlers, error) {
	storage, err := gcs.GoogleCloudStorage(ctx)
	if err != nil {
		// we can return no error
		return nil, nil
	}

	handler := &SftpHandler{
		storage: storage,
	}

	return &sftp.Handlers{
		FileGet:  handler,
		FilePut:  handler,
		FileCmd:  handler,
		FileList: handler,
	}, nil
}

// Fileread reads a file from Google Cloud Storage and returns an io.ReaderAt interface.
func (fs *SftpHandler) Fileread(r *sftp.Request) (io.ReaderAt, error) {
	reader, err := fs.storage.GetFileFromGCS(r.Context(), r.Filepath[1:])
	if err != nil {
		return nil, err
	}

	return NewReadAtBuffer(reader)
}

// Filewrite creates an object on Google Cloud Storage and returns an io.WriterAt interface.
func (fs *SftpHandler) Filewrite(r *sftp.Request) (io.WriterAt, error) {
	writer, err := fs.storage.WriteFileOnGCS(r.Context(), r.Filepath[1:])
	if err != nil {
		return nil, err
	}

	return NewWriteAtBuffer(writer, []byte{}), nil
}

// Filecmd performs various file operations on Google Cloud Storage.
//
// Supported operations:
// - Setstat: Not implemented.
// - Rename: Renames an object on GCS.
// - Remove: Removes a file from GCS.
// - Mkdir: Creates a directory in the GCS bucket.
// - Rmdir: Not implemented.
// - Symlink: Not implemented.
func (fs *SftpHandler) Filecmd(r *sftp.Request) error {
	switch r.Method {
	case "Setstat":
		return nil

	case "Rename":
		err := fs.storage.RenameFileOnGCS(r.Context(), r.Filepath[1:], r.Target)
		if err != nil {
			return err
		}
		return nil

	case "Remove":
		err := fs.storage.RemoveFileOnGcs(r.Context(), r.Filepath[1:])
		if err != nil {
			return err
		}
		return nil

	case "Mkdir":
		err := fs.storage.MakeDirOnGCS(r.Context(), r.Filepath[1:]+"/")
		if err != nil {
			return err
		}
		return nil

	case "Rmdir":
		return fmt.Errorf("rmdir is not implemented")

	case "Symlink":
		return fmt.Errorf("symlink is not implemented")

	default:
		return fmt.Errorf("unsupported file command: %s", r.Method)
	}
}

// Filelist performs file listing and stat operations on Google Cloud Storage.
//
// Supported operations:
// - List: Returns a list of files in the format of os.FileInfo.
// - Stat: Returns file stats in the format of os.FileInfo.
// - Readlink: Not implemented.
func (fs *SftpHandler) Filelist(r *sftp.Request) (sftp.ListerAt, error) {
	switch r.Method {
	case "List":
		info, err := fs.storage.ListFileOnGCS(r.Context(), r.Filepath[1:])
		if err != nil {
			return nil, err
		}
		return info, nil

	case "Stat":
		stats, err := fs.storage.StatsOnGCS(r.Context(), r.Filepath[1:])
		if err != nil {
			return nil, err
		}
		return stats, nil

	case "Readlink":
		return nil, fmt.Errorf("readlink is not implemented")

	default:
		return nil, fmt.Errorf("unsupported file list command: %s", r.Method)
	}
}
