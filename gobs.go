package gobs

import (
	"bytes"
	"io"
	"io/fs"
	"path/filepath"
	"strings"
	"time"

	"github.com/daqiancode/gobs/obs"
)

func NewOBS(accessKey, secretKey, endPoint, bucket string) (*OBS, error) {
	obsClient, err := obs.New(accessKey, secretKey, endPoint)
	if err != nil {
		return nil, err
	}
	return &OBS{
		ObsClient: obsClient,
		Bucket:    bucket,
	}, nil
}

type OBS struct {
	*obs.ObsClient
	Bucket string
}

func (s *OBS) NormPath(path string) string {
	return strings.TrimLeft(path, "/\\")
}

type ObsFile struct {
	Name         string
	IsDir        bool
	Size         int64
	LastModified time.Time
}

//Write buf -> obs path
func (s *OBS) ListFile(prefix string, max int) ([]ObsFile, error) {
	prefix = s.NormPath(prefix)
	input := &obs.ListObjectsInput{}
	input.Bucket = s.Bucket
	input.Prefix = prefix
	input.MaxKeys = max
	output, err := s.ListObjects(input)
	r := make([]ObsFile, len(output.Contents))
	for i, v := range output.Contents {
		r[i].Name = v.Key
		r[i].Size = v.Size
		r[i].LastModified = v.LastModified
		r[i].IsDir = strings.HasSuffix(v.Key, "/")
	}
	return r, err
}

//GetReader path file reader, close reader after reading
func (s *OBS) GetReader(path string) (io.ReadCloser, error) {
	path = s.NormPath(path)
	input := &obs.GetObjectInput{}
	input.Bucket = s.Bucket
	input.Key = path
	output, err := s.GetObject(input)
	if err != nil {
		return nil, err
	}
	return output.Body, nil
}

func (s *OBS) GetMeta(path string) (*obs.GetObjectMetadataOutput, error) {
	input := &obs.GetObjectMetadataInput{
		Bucket: s.Bucket,
		Key:    path,
	}
	return s.GetObjectMetadata(input)
}

//Read path file to buffer
func (s *OBS) Read(path string) (*bytes.Buffer, error) {
	path = s.NormPath(path)
	src, err := s.GetReader(path)
	if err != nil {
		return nil, err
	}
	defer src.Close()
	buf := bytes.NewBuffer(nil)
	_, err = io.Copy(buf, src)
	if err != nil {
		return nil, err
	}
	return buf, nil
}

//Write buf -> obs path
func (s *OBS) Write(buf io.Reader, path string) error {
	path = s.NormPath(path)
	input := &obs.PutObjectInput{}
	input.Bucket = s.Bucket
	input.Key = path
	input.Body = buf
	_, err := s.PutObject(input)
	return err
}

//CopyDirectory copy src/* to dst/*
func (s *OBS) CopyDirectory(src, dst string, filter ...func(relpath string) bool) error {
	return filepath.Walk(src, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		input := &obs.PutFileInput{}
		input.Bucket = s.Bucket
		input.SourceFile = path
		path, _ = filepath.Rel(src, path)
		if len(filter) > 0 && !filter[0](path) {
			return nil
		}
		input.Key = filepath.Join(dst, path)
		input.Key = strings.TrimLeft(input.Key, "/\\")
		_, err = s.PutFile(input)
		return err
	})

}
