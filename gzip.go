package benchcompress

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"reflect"
	"runtime"
	"strings"
	"sync"
	"unsafe"

	"github.com/pkg/errors"
)

type Gzip struct {
	writers sync.Pool
	readers sync.Pool
}

func NewGzip(level int) (Compressor, error) {
	if level < gzip.HuffmanOnly || level > gzip.BestCompression {
		return nil, fmt.Errorf("gzip: invalid compression level: %d", level)
	}

	return &Gzip{
		writers: sync.Pool{
			New: func() interface{} {
				w, _ := gzip.NewWriterLevel(nil, level)
				return w
			},
		},
	}, nil
}

func (g *Gzip) Compress(s string) (string, error) {
	w := g.writers.Get().(*gzip.Writer)
	defer g.writers.Put(w)

	b := &strings.Builder{}
	w.Reset(b)

	if _, err := w.Write(string2ByteSlice(s)); err != nil {
		return "", errors.Wrap(err, "write")
	}

	if err := w.Close(); err != nil {
		return "", errors.Wrap(err, "close")
	}

	return b.String(), nil
}

func (g *Gzip) Uncompress(s string) (string, error) {
	sr := bytes.NewBufferString(s)

	r, ok := g.readers.Get().(*gzip.Reader)
	if !ok || r == nil {
		r, _ = gzip.NewReader(sr)
	} else {
		r.Reset(sr)
	}
	defer g.readers.Put(r)

	b := &strings.Builder{}
	if _, err := io.Copy(b, r); err != nil {
		return "", errors.Wrap(err, "copying to builder")
	}

	if err := r.Close(); err != nil {
		return "", errors.Wrap(err, "closing reader")
	}

	return b.String(), nil
}

// Use to avoid allocating extra memory for the type cast.
func string2ByteSlice(str string) (bs []byte) {
	strHdr := (*reflect.StringHeader)(unsafe.Pointer(&str))
	sliceHdr := (*reflect.SliceHeader)(unsafe.Pointer(&bs))
	sliceHdr.Data = strHdr.Data
	sliceHdr.Len = strHdr.Len
	sliceHdr.Cap = strHdr.Len
	// This KeepAlive line is essential to make the
	// String2ByteSlice function be always valid
	// when it is provided in other custom packages.
	runtime.KeepAlive(&str)
	return
}
