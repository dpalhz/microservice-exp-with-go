package token

import (
	"crypto/rsa"
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JWTConfig struct {
	PrivateKeyPath string
	PublicKeyPath  string
	ExpireMinutes  int
	Issuer         string
}

type JWTGenerator struct {
	config      JWTConfig
	privateKey  *rsa.PrivateKey
	publicKey   *rsa.PublicKey
}

// NewJWTGenerator memuat RSA key dari path dan inisialisasi JWTGenerator
func NewJWTGenerator(config JWTConfig) (*JWTGenerator, error) {
	privKey, err := loadPrivateKey(config.PrivateKeyPath)
	if err != nil {
		return nil, err
	}

	pubKey, err := loadPublicKey(config.PublicKeyPath)
	if err != nil {
		return nil, err
	}

	return &JWTGenerator{
		config:     config,
		privateKey: privKey,
		publicKey:  pubKey,
	}, nil
}

func loadPrivateKey(path string) (*rsa.PrivateKey, error) {
	keyData, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return jwt.ParseRSAPrivateKeyFromPEM(keyData)
}

func loadPublicKey(path string) (*rsa.PublicKey, error) {
	keyData, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return jwt.ParseRSAPublicKeyFromPEM(keyData)
}

// GenerateToken membuat JWT dengan RS256
func (j *JWTGenerator) GenerateTokens(userID uuid.UUID) (string, string, error) {
	now := time.Now()

	// Access Token (short-lived)
	accessClaims := jwt.MapClaims{
		"sub": userID.String(),
		"exp": now.Add(time.Duration(j.config.ExpireMinutes) * time.Minute).Unix(),
		"iss": j.config.Issuer,
		"iat": now.Unix(),
		"type": "access",
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodRS256, accessClaims)
	accessSigned, err := accessToken.SignedString(j.privateKey)
	if err != nil {
		return "", "", err
	}

	// Refresh Token (longer-lived, e.g., 7 days)
	refreshClaims := jwt.MapClaims{
		"sub": userID.String(),
		"exp": now.Add(7 * 24 * time.Hour).Unix(),
		"iss": j.config.Issuer,
		"iat": now.Unix(),
		"type": "refresh",
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodRS256, refreshClaims)
	refreshSigned, err := refreshToken.SignedString(j.privateKey)
	if err != nil {
		return "", "", err
	}

	return accessSigned, refreshSigned, nil
}


// ValidateToken memverifikasi dan mengembalikan klaim JWT
func (j *JWTGenerator) ValidateToken(tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return j.publicKey, nil
	})

	if err != nil || !token.Valid {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("could not parse claims")
	}
	return claims, nil
}
