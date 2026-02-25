package services

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"strings"

	"google.golang.org/api/option"
	"google.golang.org/api/playintegrity/v1"
)

type IntegrityService struct {
	client      *playintegrity.Service
	packageName string
	sha256      string
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

	sha256, err := loadSigningCert()
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
		sha256:      sha256,
	}, nil
}

func loadSigningCert() (string, error) {
	cert := os.Getenv(("APP_SIGNING_CERT"))
	if cert == "" {
		return "", fmt.Errorf("could not retrieve app signing cert")
	}
	cert = strings.ReplaceAll(cert, ":", "")
	return cert, nil
}

func (s *IntegrityService) DecodeToken(token string) (*playintegrity.DecodeIntegrityTokenResponse, error) {
	req := &playintegrity.DecodeIntegrityTokenRequest{
		IntegrityToken: token,
	}

	return s.client.V1.DecodeIntegrityToken(s.packageName, req).Do()
}

func (s *IntegrityService) ValidateToken(decoded *playintegrity.DecodeIntegrityTokenResponse) (bool, string) {
	if ok, reason := s.VerifyAppIntegrity(decoded.TokenPayloadExternal.AppIntegrity); !ok {
		return false, reason
	}
	if ok, reason := s.VerifyDeviceIntegrity(decoded.TokenPayloadExternal.DeviceIntegrity); !ok {
		return false, reason
	}
	return true, ""
}

func (s *IntegrityService) VerifyAppIntegrity(appIntegrity *playintegrity.AppIntegrity) (bool, string) {
	ok := false
	switch appIntegrity.AppRecognitionVerdict {
	case "PLAY_RECOGNIZED":
		ok = true
	case "UNKNOWN":
	case "UNRECOGNIZED_VERSION":
	case "UNEVALUATED":
	}

	if !ok {
		return false, "unrecognized app"
	}

	expectedCert := strings.ToUpper(strings.ReplaceAll(s.sha256, ":", ""))
	for _, cert := range appIntegrity.CertificateSha256Digest {
		c := strings.ToUpper(strings.ReplaceAll(cert, ":", ""))
		if c == expectedCert {
			ok = true
			break
		}
	}

	if !ok {
		return false, "signing certificate mismatch"
	}

	return true, ""
}

func (s *IntegrityService) VerifyDeviceIntegrity(deviceIntegrity *playintegrity.DeviceIntegrity) (bool, string) {
	for _, v := range deviceIntegrity.DeviceRecognitionVerdict {
		if v == "MEETS_BASIC_INTEGRITY" || v == "MEETS_DEVICE_INTEGRITY" || v == "MEETS_STRONG_INTEGRITY" {
			return true, ""
		}
	}

	return false, "device integrity failed"
}
