package benchcompress

import (
	"compress/zlib"
	"fmt"
	"io"
	"strings"
	"sync"
)

type Zlib struct {
	readers sync.Pool
	writers sync.Pool
}

func NewZlib(level int) (Compressor, error) {
	if level < zlib.HuffmanOnly || level > zlib.BestCompression {
		return nil, fmt.Errorf("gzip: invalid compression level: %d", level)
	}

	return &Zlib{
		writers: sync.Pool{
			New: func() interface{} {
				w, _ := zlib.NewWriterLevel(nil, level)
				return w
			},
		},
	}, nil
}

func (z *Zlib) Compress(s string) (string, error) {
	w := z.writers.Get().(*zlib.Writer)
	defer z.writers.Put(w)

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

type reader interface {
	io.ReadCloser
	zlib.Resetter
}

func (z *Zlib) Uncompress(s string) (string, error) {
	sr := strings.NewReader(s)
	r, ok := z.readers.Get().(reader)
	if !ok || r == nil {
		zr, _ := zlib.NewReader(sr)
		r, _ = zr.(reader)
	} else {
		r.Reset(sr, nil)
	}
	defer z.readers.Put(r)

	b := &strings.Builder{}
	if _, err := io.Copy(b, r); err != nil {
		return "", err
	}

	if err := r.Close(); err != nil {
		return "", err
	}

	return b.String(), nil
}
