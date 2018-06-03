package store

import (
	"fmt"
	"io"
	"os"

	minio "github.com/minio/minio-go"
)

type minioObjectStore struct {
	Bucket string
	Client *minio.Client
	Opts   MinioOpts
}

// MinioOpts are provided to configure the Minio storage client
type MinioOpts struct {
	URL    string
	Bucket string
	Region string
	UseSSL bool
}

// NewMinioObjectStore creates and returns an ObjectStore interface that uses
// Minio client under the hood.
func NewMinioObjectStore(opts MinioOpts) (ObjectStore, error) {

	accessKey := os.Getenv("AWS_ACCESS_KEY_ID")
	secretKey := os.Getenv("AWS_SECRET_ACCESS_KEY")

	client, err := minio.New(opts.URL, accessKey, secretKey, opts.UseSSL)
	if err != nil {
		return nil, fmt.Errorf("Minio error: %s", err.Error())
	}

	if err = client.MakeBucket(opts.Bucket, opts.Region); err != nil {
		return nil, fmt.Errorf("Minio bucket error: %s", err.Error())
	}

	return &minioObjectStore{
		Bucket: opts.Bucket,
		Client: client,
		Opts:   opts,
	}, nil
}

func (m *minioObjectStore) Get(objectName string, opts minio.GetObjectOptions) (io.Reader, error) {
	return m.Client.GetObject(m.Bucket, objectName, opts)
}

func (m *minioObjectStore) Put(objectName string, reader io.Reader, size int64, opts minio.PutObjectOptions) (n int64, err error) {
	return m.Client.PutObject(m.Bucket, objectName, reader, size, opts)
}

func (m *minioObjectStore) Remove(objectName string) error {
	return m.Client.RemoveObject(m.Bucket, objectName)
}
