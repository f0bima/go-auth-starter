package usecase

import (
	"os"
)

// AuthConfig holds auth-specific configuration that does not belong in core.
type AuthConfig struct {
	JWTExpire      string
	RefreshExpire  string
	PrivateKeyPath string
	PublicKeyPath  string
}

// LoadAuthConfig reads auth-specific config from environment variables.
func LoadAuthConfig() *AuthConfig {
	return &AuthConfig{
		JWTExpire:      getEnv("JWT_EXPIRE", "15m"),
		RefreshExpire:  getEnv("REFRESH_EXPIRE", "168h"),
		PrivateKeyPath: getEnv("JWT_PRIVATE_KEY_PATH", "keys/private.pem"),
		PublicKeyPath:  getEnv("JWT_PUBLIC_KEY_PATH", "keys/public.pem"),
	}
}

func getEnv(key, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}
