//go:build !no_snappy
// +build !no_snappy

package compress

import (
	"github.com/jimyag/parquet-go/parquet"
	"github.com/klauspost/compress/snappy"
)

func init() {
	compressors[parquet.CompressionCodec_SNAPPY] = &Compressor{
		Compress: func(buf []byte) []byte {
			return snappy.Encode(nil, buf)
		},
		Uncompress: func(buf []byte) (bytes []byte, err error) {
			return snappy.Decode(nil, buf)
		},
	}
}
