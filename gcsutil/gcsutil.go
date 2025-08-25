package gcsutil

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"cloud.google.com/go/storage"
	"google.golang.org/api/iterator"
)

// Opsi A: Fungsi stateless — membuat client setiap pemanggilan.
// Mengembalikan map[string]interface{} hasil unmarshal JSON pertama yang match dengan prefix+secCode.
func GetFile(ctx context.Context, secCode, bucketName, objectPrefix string) (map[string]interface{}, error) {
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("gcs new client: %w", err)
	}
	defer client.Close()

	it := client.Bucket(bucketName).Objects(ctx, &storage.Query{
		Prefix: objectPrefix + secCode,
	})

	attrs, err := it.Next()
	if err == iterator.Done {
		return nil, errors.New("object not found")
	}
	if err != nil {
		return nil, fmt.Errorf("list objects: %w", err)
	}

	rc, err := client.Bucket(bucketName).Object(attrs.Name).NewReader(ctx)
	if err != nil {
		return nil, fmt.Errorf("new reader %q: %w", attrs.Name, err)
	}
	defer rc.Close()

	dataBytes, err := io.ReadAll(rc)
	if err != nil {
		return nil, fmt.Errorf("read object %q: %w", attrs.Name, err)
	}

	var dataMap map[string]interface{}
	if err := json.Unmarshal(dataBytes, &dataMap); err != nil {
		return nil, fmt.Errorf("unmarshal json %q: %w", attrs.Name, err)
	}

	return dataMap, nil
}

// Opsi B: Client reusable — lebih efisien bila dipanggil berulang.
// Gunakan New() sekali, lalu Close() saat shutdown.
type Client struct {
	st *storage.Client
}

func New(ctx context.Context) (*Client, error) {
	st, err := storage.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("gcs new client: %w", err)
	}
	return &Client{st: st}, nil
}

func (c *Client) Close() error {
	if c == nil || c.st == nil {
		return nil
	}
	return c.st.Close()
}

// GetFile melakukan hal yang sama seperti fungsi stateless, namun memakai client yang sama.
func (c *Client) GetFile(ctx context.Context, secCode, bucketName, objectPrefix string) (map[string]interface{}, error) {
	if c == nil || c.st == nil {
		return nil, errors.New("gcs client is nil")
	}

	it := c.st.Bucket(bucketName).Objects(ctx, &storage.Query{
		Prefix: objectPrefix + secCode,
	})

	attrs, err := it.Next()
	if err == iterator.Done {
		return nil, errors.New("object not found")
	}
	if err != nil {
		return nil, fmt.Errorf("list objects: %w", err)
	}

	rc, err := c.st.Bucket(bucketName).Object(attrs.Name).NewReader(ctx)
	if err != nil {
		return nil, fmt.Errorf("new reader %q: %w", attrs.Name, err)
	}
	defer rc.Close()

	dataBytes, err := io.ReadAll(rc)
	if err != nil {
		return nil, fmt.Errorf("read object %q: %w", attrs.Name, err)
	}

	var dataMap map[string]interface{}
	if err := json.Unmarshal(dataBytes, &dataMap); err != nil {
		return nil, fmt.Errorf("unmarshal json %q: %w", attrs.Name, err)
	}

	return dataMap, nil
}
