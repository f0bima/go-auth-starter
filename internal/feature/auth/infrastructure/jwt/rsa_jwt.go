package jwt

import (
	"crypto/rsa"
	"encoding/base64"
	"fmt"
	"math/big"
	"os"
	"time"

	"github.com/f0bima/go-auth-starter/internal/feature/auth/domain/service"
	"github.com/golang-jwt/jwt/v5"
)

// Compile-time check that RSAKeys implements domain.TokenGenerator
var _ service.TokenGenerator = (*RSAKeys)(nil)

// RSAKeys holds the RSA key pair for token signing and verification.
type RSAKeys struct {
	PrivateKey *rsa.PrivateKey
	PublicKey  *rsa.PublicKey
	KeyID      string
}

// LoadKeys reads RSA key pair from PEM files.
func LoadKeys(privateKeyPath, publicKeyPath string) (*RSAKeys, error) {
	privBytes, err := os.ReadFile(privateKeyPath)
	if err != nil {
		return nil, fmt.Errorf("could not read private key: %w", err)
	}

	privKey, err := jwt.ParseRSAPrivateKeyFromPEM(privBytes)
	if err != nil {
		return nil, fmt.Errorf("could not parse private key: %w", err)
	}

	pubBytes, err := os.ReadFile(publicKeyPath)
	if err != nil {
		return nil, fmt.Errorf("could not read public key: %w", err)
	}

	pubKey, err := jwt.ParseRSAPublicKeyFromPEM(pubBytes)
	if err != nil {
		return nil, fmt.Errorf("could not parse public key: %w", err)
	}

	return &RSAKeys{
		PrivateKey: privKey,
		PublicKey:  pubKey,
		KeyID:      "auth-key-1",
	}, nil
}

// GetJWKS returns the JSON Web Key Set for public key distribution.
func (k *RSAKeys) GetJWKS() service.JWKS {
	eBytes := big.NewInt(int64(k.PublicKey.E)).Bytes()

	jwk := service.JWK{
		Kty: "RSA",
		Alg: "RS256",
		Use: "sig",
		Kid: k.KeyID,
		N:   base64.RawURLEncoding.EncodeToString(k.PublicKey.N.Bytes()),
		E:   base64.RawURLEncoding.EncodeToString(eBytes),
	}

	return service.JWKS{
		Keys: []service.JWK{jwk},
	}
}

// GenerateToken creates a signed JWT token with the given claims.
func (k *RSAKeys) GenerateToken(userID string, email string, tokenType string, expiration time.Duration) (string, error) {
	claims := jwtClaims{
		UserID: userID,
		Email:  email,
		Type:   tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	token.Header["kid"] = k.KeyID
	return token.SignedString(k.PrivateKey)
}

// ValidateToken parses and validates a JWT token, returning the decoded claims.
func (k *RSAKeys) ValidateToken(tokenString string) (*service.TokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwtClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return k.PublicKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*jwtClaims); ok && token.Valid {
		return &service.TokenClaims{
			UserID: claims.UserID,
			Email:  claims.Email,
			Type:   claims.Type,
		}, nil
	}

	return nil, fmt.Errorf("invalid token")
}

// jwtClaims is an internal struct for JWT signing/parsing.
// Not exported - infrastructure implementation detail.
type jwtClaims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Type   string `json:"type"`
	jwt.RegisteredClaims
}
