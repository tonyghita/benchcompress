package benchcompress_test

import (
	"compress/zlib"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tonyghita/benchcompress"
)

func BenchmarkZlibCompress(b *testing.B) {
	cases := []struct {
		name  string
		level int
	}{
		{"DefaultCompression", zlib.DefaultCompression},
		{"BestCompression", zlib.BestCompression},
		{"BestSpeed", zlib.BestSpeed},
		{"HuffmanOnly", zlib.HuffmanOnly},
		{"NoCompression", zlib.NoCompression},
	}

	json := randomJSON(b)

	for _, c := range cases {
		b.Run(c.name, func(b *testing.B) {
			zl, err := benchcompress.NewZlib(c.level)
			require.NoError(b, err)
			benchmarkCompress(b, zl, json)
		})
	}
}

func BenchmarkZlibUncompress(b *testing.B) {
	cases := []struct {
		name  string
		level int
	}{
		{"DefaultCompression", zlib.DefaultCompression},
		{"BestCompression", zlib.BestCompression},
		{"BestSpeed", zlib.BestSpeed},
		{"HuffmanOnly", zlib.HuffmanOnly},
		{"NoCompression", zlib.NoCompression},
	}

	json := randomJSON(b)

	for _, c := range cases {
		b.Run(c.name, func(b *testing.B) {
			gz, err := benchcompress.NewZlib(c.level)
			require.NoError(b, err)
			benchmarkUncompress(b, gz, json)
		})
	}
}
