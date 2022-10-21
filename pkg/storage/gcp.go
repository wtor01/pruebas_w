package storage

import (
	"cloud.google.com/go/storage"
	"context"
	"errors"
	"io"
	"io/ioutil"
	"strings"
)

type Gcp struct {
	client *storage.Client
}

func (s Gcp) NewReader(ctx context.Context, path string) (io.Reader, error) {
	bucket, object, err := s.getSplitPath(path)

	if err != nil {
		return nil, err
	}

	reader, err := s.client.Bucket(bucket).Object(object).NewReader(ctx)

	return reader, err
}

func NewStorageGCP(ctx context.Context) (Storage, error) {
	client, err := storage.NewClient(ctx)
	if err != nil {
		return Gcp{}, err
	}
	return Gcp{
		client: client,
	}, nil
}

func (s Gcp) Copy(ctx context.Context, originPath, dstPath string) error {
	bucketOrigin, objectOrigin, err := s.getSplitPath(originPath)

	if err != nil {
		return err
	}
	bucketDst, objectDst, err := s.getSplitPath(dstPath)

	if err != nil {
		return err
	}

	src := s.client.Bucket(bucketOrigin).Object(objectOrigin)

	dst := s.client.Bucket(bucketDst).Object(objectDst)

	if _, err := dst.CopierFrom(src).Run(ctx); err != nil {
		return err
	}

	return nil
}

func (s Gcp) Delete(ctx context.Context, path string) error {

	bucket, object, err := s.getSplitPath(path)

	if err != nil {
		return err
	}

	src := s.client.Bucket(bucket).Object(object).Delete(ctx)

	if src != nil {
		return err
	}

	return nil
}

func (s Gcp) getSplitPath(path string) (bucket string, object string, err error) {
	paths := strings.Split(path, "/")

	if len(paths) < 2 {
		return "", "", errors.New("invalid path")
	}
	bucket = paths[0]
	object = strings.Join(paths[1:], "/")

	return bucket, object, nil
}

func (s Gcp) Close() error {
	return s.client.Close()
}

func (s Gcp) ReadAll(ctx context.Context, path string) ([]byte, error) {
	bucket, object, err := s.getSplitPath(path)

	if err != nil {
		return []byte{}, err
	}

	reader, err := s.client.Bucket(bucket).Object(object).NewReader(ctx)

	if err != nil {
		return []byte{}, errors.New("invalid path")
	}

	defer reader.Close()

	return ioutil.ReadAll(reader)
}
