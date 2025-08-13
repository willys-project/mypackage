package functions

import (
	"context"
	"fmt"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
	"github.com/golang-jwt/jwt/v5" // Pastikan kamu sudah menginstal package ini
)

// getJWTClaim mengekstrak satu klaim dari token JWT
func GetJWTClaim(token *jwt.Token, claimKey string) (interface{}, error) {
	if token == nil {
		return nil, fmt.Errorf("nil token provided")
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if claimValue, found := claims[claimKey]; found {
			return claimValue, nil
		}
		return nil, fmt.Errorf("claim not found: %s", claimKey)
	}
	return nil, fmt.Errorf("invalid token or claims")
}

func GetSecret(projectID, secretID string) (string, error) {
	// Create a context and Secret Manager client
	ctx := context.Background()
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to create secretmanager client: %v", err)
	}
	defer client.Close()

	// Build the secret version name
	name := fmt.Sprintf("projects/%s/secrets/%s/versions/latest", projectID, secretID)

	// Access the secret version
	accessRequest := &secretmanagerpb.AccessSecretVersionRequest{
		Name: name,
	}
	result, err := client.AccessSecretVersion(ctx, accessRequest)
	if err != nil {
		return "", fmt.Errorf("failed to access secret version: %v", err)
	}

	// Convert the secret payload to a string
	secretValue := string(result.Payload.Data)

	return secretValue, nil
}
