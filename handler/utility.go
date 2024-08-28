package handler

import (
	"bytes"
	"io"
	"sync"
)

// NewReadAtBuffer reads all data from the provided ReadCloser, closes it,
// and returns a ReaderAt backed by a buffer containing the read data.
func NewReadAtBuffer(r io.ReadCloser) (io.ReaderAt, error) {
	buf, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	if err = r.Close(); err != nil {
		return nil, err
	}

	return bytes.NewReader(buf), nil
}

// WriteAtBuffer is a thread-safe buffer that supports writing at specific offsets
// and provides an io.WriteCloser to write the buffer's content to an underlying writer.
type WriteAtBuffer struct {
	buf         []byte
	mu          sync.Mutex
	GrowthCoeff float64 // GrowthCoeff defines the growth rate of the internal buffer.
	Writer      io.WriteCloser
}

// NewWriteAtBuffer creates a new WriteAtBuffer with an initial buffer and an io.WriteCloser.
func NewWriteAtBuffer(w io.WriteCloser, buf []byte) *WriteAtBuffer {
	return &WriteAtBuffer{
		Writer:      w,
		buf:         buf,
		GrowthCoeff: 1.0, // Default growth coefficient is 1.
	}
}

// WriteAt writes the provided data at the specified position in the buffer.
// It expands the buffer if necessary.
func (b *WriteAtBuffer) WriteAt(p []byte, pos int64) (int, error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	pLen := len(p)
	expLen := pos + int64(pLen)

	// Expand buffer if needed
	if expLen > int64(len(b.buf)) {
		b.expandBuffer(expLen)
	}

	copy(b.buf[pos:], p)
	return pLen, nil
}

// expandBuffer expands the internal buffer to accommodate the specified length.
func (b *WriteAtBuffer) expandBuffer(expLen int64) {
	if b.GrowthCoeff < 1 {
		b.GrowthCoeff = 1
	}

	newBuf := make([]byte, expLen, int64(b.GrowthCoeff*float64(expLen)))
	copy(newBuf, b.buf)
	b.buf = newBuf
}

// Bytes returns a copy of the data in the buffer.
func (b *WriteAtBuffer) Bytes() []byte {
	b.mu.Lock()
	defer b.mu.Unlock()

	// Return a copy to avoid modifying the internal buffer directly
	return append([]byte(nil), b.buf...)
}

// Close writes the buffer's contents to the underlying writer and closes it.
func (b *WriteAtBuffer) Close() error {
	b.mu.Lock()
	defer b.mu.Unlock()

	if _, err := b.Writer.Write(b.buf); err != nil {
		return err
	}

	return b.Writer.Close()
}
