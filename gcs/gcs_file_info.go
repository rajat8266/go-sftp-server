package gcs

import (
	"io"
	"os"
	"time"

	"cloud.google.com/go/storage"
)

// GcsFileInfo represents a file or directory in Google Cloud Storage.
type GcsFileInfo struct {
	objAttr *storage.ObjectAttrs
	prefix  string
}

// Name returns the base name of the file or directory.
func (f *GcsFileInfo) Name() string {
	if f.objAttr.Prefix != "" {
		return f.objAttr.Prefix[len(f.prefix) : len(f.objAttr.Prefix)-1]
	}
	return f.objAttr.Name[len(f.prefix):]
}

// Size returns the length in bytes for regular files and is system-dependent for others.
func (f *GcsFileInfo) Size() int64 {
	return f.objAttr.Size
}

// Mode returns the file mode bits. It marks directories with a specific mode.
func (f *GcsFileInfo) Mode() os.FileMode {
	if f.IsDir() {
		return os.ModeDir | 0777
	}
	return 0777
}

// ModTime returns the modification time of the file or directory.
// For directories, it returns the current time.
func (f *GcsFileInfo) ModTime() time.Time {
	if f.IsDir() {
		return time.Now()
	}
	return f.objAttr.Updated
}

// IsDir checks if the GcsFileInfo represents a directory.
func (f *GcsFileInfo) IsDir() bool {
	if f.objAttr.Prefix != "" {
		return true
	}
	return len(f.objAttr.Name) > 0 && f.objAttr.Name[len(f.objAttr.Name)-1:] == "/" && f.objAttr.Size == 0
}

// Sys returns the underlying data source, which is nil in this case.
func (f *GcsFileInfo) Sys() interface{} {
	return nil
}

// ListerAt is a custom type for listing file information with support for offset-based access.
type ListerAt []os.FileInfo

// ListAt copies file information into the provided slice starting from the given offset.
// It returns the number of copied entries and an error if the end of the list is reached.
func (f ListerAt) ListAt(ls []os.FileInfo, offset int64) (int, error) {
	if offset >= int64(len(f)) {
		return 0, io.EOF
	}
	n := copy(ls, f[offset:])
	if n < len(ls) {
		return n, io.EOF
	}
	return n, nil
}
