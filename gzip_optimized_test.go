package benchcompress_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tonyghita/benchcompress"
	"github.com/klauspost/compress/gzip"
)

func BenchmarkGzipOptimizedCompress(b *testing.B) {
	cases := []struct {
		name  string
		level int
	}{
		{"DefaultCompression", gzip.DefaultCompression},
		{"BestCompression", gzip.BestCompression},
		{"BestSpeed", gzip.BestSpeed},
		{"HuffmanOnly", gzip.HuffmanOnly},
		{"NoCompression", gzip.NoCompression},
	}

	json := randomJSON(b)

	for _, c := range cases {
		b.Run(c.name, func(b *testing.B) {
			gz, err := benchcompress.NewGzipOptimized(c.level)
			require.NoError(b, err)
			benchmarkCompress(b, gz, json)
		})
	}
}

func BenchmarkGzipOptimizedUncompress(b *testing.B) {
	cases := []struct {
		name  string
		level int
	}{
		{"DefaultCompression", gzip.DefaultCompression},
		{"BestCompression", gzip.BestCompression},
		{"BestSpeed", gzip.BestSpeed},
		{"HuffmanOnly", gzip.HuffmanOnly},
		{"NoCompression", gzip.NoCompression},
	}

	json := randomJSON(b)

	for _, c := range cases {
		b.Run(c.name, func(b *testing.B) {
			gz, err := benchcompress.NewGzipOptimized(c.level)
			require.NoError(b, err)
			benchmarkUncompress(b, gz, json)
		})
	}
}
