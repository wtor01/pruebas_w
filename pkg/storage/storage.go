package storage

import (
	"context"
	"io"
)

//go:generate mockery --case=snake --outpkg=mocks --output=./mocks --name=Storage
type Storage interface {
	Close() error
	Delete(ctx context.Context, path string) error
	ReadAll(ctx context.Context, path string) ([]byte, error)
	Copy(ctx context.Context, originPath, dstPath string) error
	NewReader(ctx context.Context, path string) (io.Reader, error)
}

type StorageCreator = func(ctx context.Context) (Storage, error)
