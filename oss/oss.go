package oss

import (
	"context"
	"io"
)

type Operation interface {
	Put(ctx context.Context, key string, data io.Reader, size int64, mimeType string) error
	Base64Put(ctx context.Context, key string, raw []byte, mimeType string) error
	Delete(key string) error
	Bucket() string
	Domain() string
	URL(key string) string
}
