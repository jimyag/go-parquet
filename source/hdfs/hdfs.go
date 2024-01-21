package hdfs

import (
	"github.com/colinmarc/hdfs/v2"

	"github.com/jimyag/go-parquet/source"
)

type HdfsFile struct {
	Hosts []string
	User  string

	Client     *hdfs.Client
	FilePath   string
	FileReader *hdfs.FileReader
	FileWriter *hdfs.FileWriter
}

func NewHdfsFileWriter(hosts []string, user string, name string) (source.ParquetFile, error) {
	res := &HdfsFile{
		Hosts:    hosts,
		User:     user,
		FilePath: name,
	}
	return res.Create(name)
}

func NewHdfsFileReader(hosts []string, user string, name string) (source.ParquetFile, error) {
	res := &HdfsFile{
		Hosts:    hosts,
		User:     user,
		FilePath: name,
	}
	return res.Open(name)
}

func (f *HdfsFile) Create(name string) (source.ParquetFile, error) {
	var err error
	hf := new(HdfsFile)
	hf.Hosts = f.Hosts
	hf.User = f.User
	hf.Client, err = hdfs.NewClient(hdfs.ClientOptions{
		Addresses: hf.Hosts,
		User:      hf.User,
	})
	hf.FilePath = name
	if err != nil {
		return hf, err
	}
	hf.FileWriter, err = hf.Client.Create(name)
	return hf, err

}
func (f *HdfsFile) Open(name string) (source.ParquetFile, error) {
	var (
		err error
	)
	if name == "" {
		name = f.FilePath
	}

	hf := new(HdfsFile)
	hf.Hosts = f.Hosts
	hf.User = f.User
	hf.Client, err = hdfs.NewClient(hdfs.ClientOptions{
		Addresses: hf.Hosts,
		User:      hf.User,
	})
	hf.FilePath = name
	if err != nil {
		return hf, err
	}
	hf.FileReader, err = hf.Client.Open(name)
	return hf, err
}
func (f *HdfsFile) Seek(offset int64, pos int) (int64, error) {
	return f.FileReader.Seek(offset, pos)
}

func (f *HdfsFile) Read(b []byte) (cnt int, err error) {
	var n int
	ln := len(b)
	for cnt < ln {
		n, err = f.FileReader.Read(b[cnt:])
		cnt += n
		if err != nil {
			break
		}
	}
	return cnt, err
}

func (f *HdfsFile) Write(b []byte) (n int, err error) {
	return f.FileWriter.Write(b)
}

func (f *HdfsFile) Close() error {
	if f.FileReader != nil {
		f.FileReader.Close()
	}
	if f.FileWriter != nil {
		f.FileWriter.Close()
	}
	if f.Client != nil {
		f.Client.Close()
	}
	return nil
}
