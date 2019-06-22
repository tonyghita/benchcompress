package benchcompress

type Compressor interface {
	Compress(string) (string, error)
	Uncompress(string) (string, error)
}
