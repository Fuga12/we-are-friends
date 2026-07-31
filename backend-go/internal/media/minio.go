package media

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"strings"
	"time"

	"github.com/Fuga12/we-are-friends/backend-go/internal/config"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// Client provides URLs for objects. Database records deliberately contain only
// object keys, never transport-specific links.
type Client struct {
	client     *minio.Client
	bucket     string
	publicBase string
}

func New(cfg config.Config) (*Client, error) {
	endpoint, secure, err := endpoint(cfg.MinIOEndpoint)
	if err != nil {
		return nil, err
	}
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.MinIOAccessKey, cfg.MinIOSecretKey, ""),
		Secure: secure,
	})
	if err != nil {
		return nil, fmt.Errorf("create MinIO client: %w", err)
	}
	return &Client{client: client, bucket: cfg.MinIOBucket, publicBase: strings.TrimRight(cfg.MediaPublicBaseURL, "/")}, nil
}

func endpoint(raw string) (string, bool, error) {
	if !strings.Contains(raw, "://") {
		return raw, false, nil
	}
	u, err := url.Parse(raw)
	if err != nil || u.Host == "" || (u.Scheme != "http" && u.Scheme != "https") {
		return "", false, fmt.Errorf("MINIO_ENDPOINT must be host:port or an HTTP URL")
	}
	return u.Host, u.Scheme == "https", nil
}

// URL prefers stable public URLs for public media and otherwise makes a
// one-hour presigned URL for a private bucket.
func (c *Client) URL(ctx context.Context, objectKey string) (string, error) {
	if c.publicBase != "" {
		return c.publicBase + "/" + strings.TrimLeft(objectKey, "/"), nil
	}
	url, err := c.client.PresignedGetObject(ctx, c.bucket, objectKey, time.Hour, nil)
	if err != nil {
		return "", fmt.Errorf("presign object %q: %w", objectKey, err)
	}
	return url.String(), nil
}

func (c *Client) EnsureBucket(ctx context.Context) error {
	exists, err := c.client.BucketExists(ctx, c.bucket)
	if err != nil {
		return fmt.Errorf("check bucket: %w", err)
	}
	if !exists {
		if err := c.client.MakeBucket(ctx, c.bucket, minio.MakeBucketOptions{}); err != nil {
			return fmt.Errorf("create bucket: %w", err)
		}
	}
	return nil
}

func (c *Client) Upload(ctx context.Context, objectKey, contentType string, size int64, body io.Reader) error {
	_, err := c.client.PutObject(ctx, c.bucket, objectKey, body, size, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		return fmt.Errorf("upload object %q: %w", objectKey, err)
	}
	return nil
}

// Delete removes an object from the configured S3-compatible bucket.
func (c *Client) Delete(ctx context.Context, objectKey string) error {
	if err := c.client.RemoveObject(ctx, c.bucket, objectKey, minio.RemoveObjectOptions{}); err != nil {
		return fmt.Errorf("delete object %q: %w", objectKey, err)
	}
	return nil
}
