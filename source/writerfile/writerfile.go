package writerfile

import (
	"io"

	"github.com/jimyag/go-parquet/source"
)

type WriterFile struct {
	Writer io.Writer
}

func NewWriterFile(writer io.Writer) source.ParquetFile {
	return &WriterFile{Writer: writer}
}

func (w *WriterFile) Create(name string) (source.ParquetFile, error) {
	return w, nil
}

func (w *WriterFile) Open(name string) (source.ParquetFile, error) {
	return w, nil
}

func (w *WriterFile) Seek(offset int64, pos int) (int64, error) {
	return 0, nil
}

func (w *WriterFile) Read(b []byte) (int, error) {
	return 0, nil
}

func (w *WriterFile) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func (w *WriterFile) Close() error {
	return nil
}
