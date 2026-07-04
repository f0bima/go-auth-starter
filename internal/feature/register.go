package feature

import (
	"log"

	"github.com/f0bima/go-auth-starter/internal/feature/auth/delivery/http/handler"
	"github.com/f0bima/go-auth-starter/internal/feature/auth/delivery/http/route"
	"github.com/f0bima/go-auth-starter/internal/feature/auth/infrastructure/config"
	"github.com/f0bima/go-auth-starter/internal/feature/auth/infrastructure/jwt"
	"github.com/f0bima/go-auth-starter/internal/feature/auth/infrastructure/ldap"
	"github.com/f0bima/go-auth-starter/internal/feature/auth/infrastructure/persistence/repository"
	"github.com/f0bima/go-auth-starter/internal/feature/auth/usecase"
	"github.com/f0bima/go-core/bootstrap"
)

// Register wires up all auth dependencies and registers routes on the core App.
func Register(app *bootstrap.App) {
	authCfg := config.LoadAuthConfig()

	// Load RSA Keys
	rsaKeys, err := jwt.LoadKeys(authCfg.PrivateKeyPath, authCfg.PublicKeyPath)
	if err != nil {
		log.Fatalf("Failed to load RSA keys: %v", err)
	}

	// Initialize Repositories
	userRepo := repository.NewUserRepository(app.DB)
	refreshTokenRepo := repository.NewRefreshTokenRepository(app.DB)

	// Initialize LDAP Client
	ldapClient := ldap.NewLDAPClient(authCfg)

	// Initialize UseCase
	authUseCase := usecase.NewAuthUseCase(userRepo, refreshTokenRepo, authCfg, rsaKeys, app.Tracer)
	ldapLoginUseCase := usecase.NewLDAPLoginUseCase(ldapClient, userRepo, refreshTokenRepo, rsaKeys)

	// Initialize Handlers
	loginHandler := handler.NewLoginHandler(authUseCase)
	ldapLoginHandler := handler.NewLDAPLoginHandler(ldapLoginUseCase)
	registerHandler := handler.NewRegisterHandler(authUseCase)
	refreshHandler := handler.NewRefreshHandler(authUseCase)
	jwksHandler := handler.NewJwksHandler(authUseCase)
	meHandler := handler.NewMeHandler(authUseCase)
	panicHandler := handler.NewPanicHandler()

	// Register routes
	route.RegisterRoutes(
		app.Router,
		rsaKeys.PublicKey,
		registerHandler,
		loginHandler,
		ldapLoginHandler,
		refreshHandler,
		jwksHandler,
		meHandler,
		panicHandler,
	)
}
