package services

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"

	"google.golang.org/api/option"
	"google.golang.org/api/playintegrity/v1"
)

type IntegrityService struct {
	client      *playintegrity.Service
	packageName string
}

func NewIntegrityService() (*IntegrityService, error) {
	ctx := context.Background()

	packageName := os.Getenv("APP_PACKAGE_NAME")
	if packageName == "" {
		return nil, fmt.Errorf("could not retrieve app package name")
	}

	b64creds := os.Getenv("GOOGLE_SERVICE_ACCOUNT_JSON")
	if b64creds == "" {
		return nil, fmt.Errorf("could not retrieve service account json")
	}
	jsonCreds, err := base64.StdEncoding.DecodeString(b64creds)
	if err != nil {
		return nil, err
	}

	client, err := playintegrity.NewService(ctx, option.WithAuthCredentialsJSON(
		option.ServiceAccount, jsonCreds,
	))
	if err != nil {
		return nil, err
	}

	return &IntegrityService{
		client:      client,
		packageName: packageName,
	}, nil
}

func (s *IntegrityService) DecodeToken(token string) (*playintegrity.DecodeIntegrityTokenResponse, error) {
	req := &playintegrity.DecodeIntegrityTokenRequest{
		IntegrityToken: token,
	}

	return s.client.V1.DecodeIntegrityToken(s.packageName, req).Do()
}
