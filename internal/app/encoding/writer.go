package encoding

import (
	"compress/gzip"
	"net/http"
)

type compressWriter struct {
	writer   http.ResponseWriter
	gzWriter *gzip.Writer
}

func newCompressWriter(w http.ResponseWriter) *compressWriter {
	return &compressWriter{
		writer:   w,
		gzWriter: gzip.NewWriter(w),
	}
}

// Header Implements method of ResponseWriter interface
func (c *compressWriter) Header() http.Header {
	return c.writer.Header()
}

// WriteHeader Implements method of ResponseWriter interface
func (c *compressWriter) WriteHeader(statusCode int) {
	c.writer.WriteHeader(statusCode)
}

// Write Implements method of ResponseWriter interface
func (c *compressWriter) Write(p []byte) (int, error) {
	return c.gzWriter.Write(p)
}

// Close Implements method of ResponseWriter interface
func (c *compressWriter) Close() error {
	return c.gzWriter.Close()
}
