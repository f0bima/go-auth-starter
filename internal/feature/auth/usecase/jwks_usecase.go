package usecase

import (
	"github.com/f0bima/go-auth-starter/internal/feature/auth/domain/service"
)

func (s *authUseCase) GetJWKS() service.JWKS {
	return s.tokenGen.GetJWKS()
}
