package benchcompress

import (
	"io"
	"strings"
	"sync"

	"github.com/klauspost/compress/gzip"
)

type GzipOptimized struct {
	readers sync.Pool
	writers sync.Pool
}

func NewGzipOptimized(level int) (Compressor, error) {
	return &GzipOptimized{
		writers: sync.Pool{
			New: func() interface{} {
				w, _ := gzip.NewWriterLevel(nil, level)
				return w
			},
		},
	}, nil
}

func (g *GzipOptimized) Compress(s string) (string, error) {
	w := g.writers.Get().(*gzip.Writer)
	defer g.writers.Put(w)

	b := &strings.Builder{}
	w.Reset(b)

	if _, err := w.Write(string2ByteSlice(s)); err != nil {
		return "", err
	}

	if err := w.Close(); err != nil {
		return "", err
	}

	return b.String(), nil
}

func (g *GzipOptimized) Uncompress(s string) (string, error) {
	sr := strings.NewReader(s)

	r, ok := g.readers.Get().(*gzip.Reader)
	if !ok || r == nil {
		r, _ = gzip.NewReader(sr)
	} else {
		r.Reset(sr)
	}
	defer g.readers.Put(r)

	b := &strings.Builder{}
	if _, err := io.Copy(b, r); err != nil {
		return "", err
	}

	if err := r.Close(); err != nil {
		return "", err
	}

	return b.String(), nil
}
