package benchcompress_test

import (
	"compress/flate"
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tonyghita/benchcompress"
)

func TestCompressors(t *testing.T) {
	cases := []struct {
		name string
		comp func(level int) (benchcompress.Compressor, error)
	}{
		{"Gzip", benchcompress.NewGzip},
		{"GzipOptimized", benchcompress.NewGzipOptimized},
		{"Zlib", benchcompress.NewZlib},
	}

	json := randomJSON(t)

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			comp, err := c.comp(flate.DefaultCompression)
			require.NoError(t, err)
			require.NotNil(t, comp)

			c, err := comp.Compress(json)
			require.NoError(t, err)
			require.NotEmpty(t, c)

			u, err := comp.Uncompress(c)
			require.NoError(t, err)
			require.Equal(t, json, u)
		})
	}
}

func benchmarkCompress(b *testing.B, c benchcompress.Compressor, in string) {
	var out string
	var err error
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		out, err = c.Compress(in)
	}

	require.NoError(b, err)
	require.NotEmpty(b, out)

	// b.Logf("%+.2f%%", (float64(len(out)-len(in))/float64(len(in)))*100.0)
}

func benchmarkUncompress(b *testing.B, c benchcompress.Compressor, in string) {
	comp, err := c.Compress(in)
	require.NoError(b, err)
	require.NotEmpty(b, comp)
	var out string
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		out, err = c.Uncompress(comp)
	}

	require.NoError(b, err)
	require.Equal(b, in, out)
}

func randomJSON(tb testing.TB) string {
	json, err := ioutil.ReadFile(filepath.Join("testdata", "random.json"))
	require.NoError(tb, err)

	return string(json)
}
