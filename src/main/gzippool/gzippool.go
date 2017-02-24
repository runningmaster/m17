package gzippool

import (
	"bufio"
	"bytes"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"sync"

	"github.com/klauspost/compress/gzip"
)

var (
	readerPool = sync.Pool{
		New: func() interface{} {
			// hack: gzip empty string for init reader
			r := bytes.NewReader([]byte{
				0x1f, 0x8b, 0x8, 0x0, 0x0, 0x9, 0x6e, 0x88, 0x0, 0xff, 0x1,
				0x0, 0x0, 0xff, 0xff, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
			})
			v, _ := gzip.NewReader(r)
			return v
		},
	}

	writerPool = sync.Pool{
		New: func() interface{} {
			return gzip.NewWriter(ioutil.Discard)
		},
	}
)

// GetReader gets reader from pool.
func GetReader() *gzip.Reader {
	return readerPool.Get().(*gzip.Reader)
}

// PutReadCloser closes reader and puts it back to the pool.
func PutReader(c io.Closer) {
	if c != nil {
		_ = c.Close()
	}
	readerPool.Put(c)
}

// GetWriter gets writer from pool.
func GetWriter() *gzip.Writer {
	return writerPool.Get().(*gzip.Writer)
}

// PutWriter closes writer and puts it back to the pool.
func PutWriter(c io.Closer) {
	if c != nil {
		_ = c.Close()
	}
	writerPool.Put(c)
}

// NewResponseWriter returns gzip-wrapper for http.ResponseWriter
func NewResponseWriter(w *gzip.Writer, rw http.ResponseWriter) http.ResponseWriter {
	return &responseWriter{
		Writer:         w,
		ResponseWriter: rw,
	}
}

type responseWriter struct {
	io.Writer
	http.ResponseWriter
}

func (w *responseWriter) WriteHeader(code int) {
	if code == http.StatusNoContent {
		w.ResponseWriter.Header().Del("Content-Encoding")
	}
	w.ResponseWriter.WriteHeader(code)
}

func (w responseWriter) Write(b []byte) (int, error) {
	if w.ResponseWriter.Header().Get("Content-Type") == "" {
		w.ResponseWriter.Header().Set("Content-Type", http.DetectContentType(b))
	}

	n, err := w.Writer.Write(b)
	if n == 0 {
		w.ResponseWriter.Header().Del("Content-Encoding")
	}

	return n, err
}

func (w responseWriter) Flush() error {
	return w.Writer.(*gzip.Writer).Flush()
}

func (w responseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return w.ResponseWriter.(http.Hijacker).Hijack()
}

func (w *responseWriter) CloseNotify() <-chan bool {
	return w.ResponseWriter.(http.CloseNotifier).CloseNotify()
}

////////////
