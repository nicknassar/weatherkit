package weatherkit

import (
	"crypto/ecdsa"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

var TokenDuration time.Duration = time.Hour
var MaxClockSkew time.Duration = 5 * time.Minute

type credentials struct {
	KeyId            string
	TeamId           string
	ServiceId        string
	privateKey       *ecdsa.PrivateKey
	lastToken        string
	lastTokenCreated time.Time
}

func newCredentials(keyId, teamId, serviceId, privateKey string) (*credentials, error) {
	if len(keyId) == 0 || len(teamId) == 0 || len(serviceId) == 0 || len(privateKey) == 0 {
		return nil, fmt.Errorf("invalid credentials")
	}
	if pk, err := jwt.ParseECPrivateKeyFromPEM([]byte(privateKey)); err != nil {
		return nil, fmt.Errorf("invalid private key: %w", err)
	} else {
		return &credentials{
			KeyId:      keyId,
			TeamId:     teamId,
			ServiceId:  serviceId,
			privateKey: pk,
		}, nil
	}
}

func (creds *credentials) Token() (string, error) {
	if !creds.hasValidToken() {
		if token, err := creds.generateToken(); err != nil {
			return "", err
		} else {
			creds.lastToken = token
			creds.lastTokenCreated = time.Now()
		}
	}
	return creds.lastToken, nil
}

func (creds *credentials) generateToken() (string, error) {
	now := time.Now()

	t := jwt.Token{
		Header: map[string]interface{}{
			"alg": jwt.SigningMethodES256.Alg(),
			"kid": creds.KeyId,
			"id":  creds.TeamId + "." + creds.ServiceId,
		},
		Claims: jwt.MapClaims{
			"iss": creds.TeamId,
			"sub": creds.ServiceId,
			"iat": now.Unix(),
			"exp": now.Add(TokenDuration + MaxClockSkew).Unix(),
		},
		Method: jwt.SigningMethodES256,
	}

	return t.SignedString(creds.privateKey)
}

func (creds *credentials) hasValidToken() bool {
	return len(creds.lastToken) > 0 && creds.lastTokenCreated.Add(TokenDuration).After(time.Now())
}
