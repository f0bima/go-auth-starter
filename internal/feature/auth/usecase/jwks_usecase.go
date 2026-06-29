package usecase

import (
	auth "github.com/f0bima/go-auth-starter/internal/feature/auth/domain/repository"
)

func (s *authUseCase) GetJWKS() auth.JWKS {
	return s.tokenGen.GetJWKS()
}
