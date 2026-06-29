package feature

import (
	"log"

	delivhttp "github.com/f0bima/go-auth-starter/internal/feature/auth/delivery/http"
	"github.com/f0bima/go-auth-starter/internal/feature/auth/infrastructure/jwt"
	"github.com/f0bima/go-auth-starter/internal/feature/auth/infrastructure/persistence/repository"
	"github.com/f0bima/go-auth-starter/internal/feature/auth/usecase"
	"github.com/f0bima/go-core/bootstrap"
)

// Register wires up all auth dependencies and registers routes on the core App.
func Register(app *bootstrap.App) {
	authCfg := usecase.LoadAuthConfig()

	// Load RSA Keys
	rsaKeys, err := jwt.LoadKeys(authCfg.PrivateKeyPath, authCfg.PublicKeyPath)
	if err != nil {
		log.Fatalf("Failed to load RSA keys: %v", err)
	}

	// Initialize Repositories
	userRepo := repository.NewUserRepository(app.DB)
	refreshTokenRepo := repository.NewRefreshTokenRepository(app.DB)

	// Initialize UseCase
	useCase := usecase.NewAuthUseCase(userRepo, refreshTokenRepo, authCfg, rsaKeys, app.Tracer)

	// Initialize Controller and register routes
	authController := delivhttp.NewAuthController(useCase)
	authController.RegisterRoutes(app.Router)
}
