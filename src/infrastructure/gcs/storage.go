package gcs

import (
	"cloud.google.com/go/storage"
	"context"
	"io"
	"log"
)

type Service interface {
	Read(ctx context.Context, name string) (io.ReadCloser, error)
	Write(ctx context.Context, object string, file io.Reader) error
	Bucket() string
}

func NewService(bucket string) Service {
	return &service{bucket: bucket}
}

type service struct {
	bucket     string
}

func (s *service) Bucket() string {
	return s.bucket
}

func (s *service) new(ctx context.Context) (*storage.Client, error) {
	return storage.NewClient(ctx)
}

// https://github.com/300481/pricenotifier/blob/b7dc3d672dc0c99d056f7253da6517e2edb29a06/pkg/persistence/persistence.go
func (s *service) Read(ctx context.Context, name string) (io.ReadCloser, error) {
	client, err := s.new(ctx)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	log.Println(s.bucket, name)
	return client.Bucket(s.bucket).Object(name).NewReader(ctx)
}

func (s *service) Write(ctx context.Context, object string, file io.Reader) error {
	//defer file.Close()
	client, err := s.new(ctx)
	if err != nil {
		return err
	}
	defer client.Close()

	w :=  client.Bucket(s.bucket).Object(object).NewWriter(ctx)
	// Make public
	//w.ACL = []storage.ACLRule{{Entity: storage.AllUsers, Role: storage.RoleReader}}
	defer w.Close()
	if _, err := io.Copy(w, file); err != nil {
		return err
	}
	return nil
}

func (s *service) Delete(ctx context.Context, object string) error {
	client, err := s.new(ctx)
	if err != nil {
		return err
	}
	defer client.Close()

	return client.Bucket(s.bucket).Object(object).Delete(ctx)
}
