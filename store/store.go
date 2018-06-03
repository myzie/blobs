package store

//go:generate mockgen -source=store.go -package store -destination mock.go

import (
	"io"

	minio "github.com/minio/minio-go"
)

// ObjectStore is an interface used to put, get, and remove objects
type ObjectStore interface {

	// Get an object from storage
	Get(objectName string, opts minio.GetObjectOptions) (io.Reader, error)

	// Put an object into storage
	Put(objectName string, reader io.Reader, size int64,
		opts minio.PutObjectOptions) (n int64, err error)

	// Remove an object from storage
	Remove(objectName string) error
}
