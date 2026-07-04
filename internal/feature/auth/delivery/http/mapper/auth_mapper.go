package mapper

import (
	"github.com/f0bima/go-auth-starter/internal/feature/auth/delivery/http/dto"
	"github.com/f0bima/go-auth-starter/internal/feature/auth/domain/entity"
	"github.com/f0bima/go-auth-starter/internal/feature/auth/domain/service"
)

// ToUserResponse converts a domain User to an HTTP response DTO.
func ToUserResponse(user *entity.User) dto.UserResponse {
	return dto.UserResponse{
		ID:    user.ID,
		Email: user.Email,
	}
}

// ToTokenResponse creates a token response DTO from access and refresh tokens.
func ToTokenResponse(accessToken, refreshToken string) dto.TokenResponse {
	return dto.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
}

// ToJWKSResponse converts domain JWKS to an HTTP response DTO.
func ToJWKSResponse(jwks service.JWKS) dto.JWKSResponse {
	keys := make([]dto.JWKResponse, len(jwks.Keys))
	for i, k := range jwks.Keys {
		keys[i] = dto.JWKResponse{
			Kty: k.Kty,
			Alg: k.Alg,
			Use: k.Use,
			Kid: k.Kid,
			N:   k.N,
			E:   k.E,
		}
	}
	return dto.JWKSResponse{Keys: keys}
}
