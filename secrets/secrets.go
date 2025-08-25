package secrets

import (
	"context"
	"fmt"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	secretmanagerpb "cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
)

// ====== Opsi 1: Fungsi Cepat (stateless) ======
// Get membuat client setiap kali dipanggil (cukup untuk script/CLI).
func Get(projectID, secretID string) (string, error) {
	ctx := context.Background()
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to create secretmanager client: %w", err)
	}
	defer client.Close()

	name := fmt.Sprintf("projects/%s/secrets/%s/versions/latest", projectID, secretID)
	res, err := client.AccessSecretVersion(ctx, &secretmanagerpb.AccessSecretVersionRequest{
		Name: name,
	})
	if err != nil {
		return "", fmt.Errorf("failed to access secret version: %w", err)
	}
	return string(res.Payload.GetData()), nil
}

// ====== Opsi 2: Client Reusable (disarankan untuk server/long-running) ======

// Client membungkus secretmanager.Client agar bisa di-reuse.
type Client struct {
	sm *secretmanager.Client
}

// New membuat Client reusable. Panggil Close() saat shutdown.
func New(ctx context.Context) (*Client, error) {
	sm, err := secretmanager.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create secretmanager client: %w", err)
	}
	return &Client{sm: sm}, nil
}

// Close menutup koneksi underlying Secret Manager.
func (c *Client) Close() error {
	if c == nil || c.sm == nil {
		return nil
	}
	return c.sm.Close()
}

// Get mengambil secret versi "latest" menggunakan client yang sama.
func (c *Client) Get(ctx context.Context, projectID, secretID string) (string, error) {
	if c == nil || c.sm == nil {
		return "", fmt.Errorf("secrets.Client is nil")
	}
	name := fmt.Sprintf("projects/%s/secrets/%s/versions/latest", projectID, secretID)
	res, err := c.sm.AccessSecretVersion(ctx, &secretmanagerpb.AccessSecretVersionRequest{
		Name: name,
	})
	if err != nil {
		return "", fmt.Errorf("failed to access secret version: %w", err)
	}
	return string(res.Payload.GetData()), nil
}
