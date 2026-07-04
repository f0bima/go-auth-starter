package config

import (
	"os"
)

// AuthConfig holds auth-specific configuration that does not belong in core.
type AuthConfig struct {
	JWTExpire      string
	RefreshExpire  string
	PrivateKeyPath string
	PublicKeyPath  string

	// LDAP Config
	LDAPURL          string
	LDAPBaseDN       string
	LDAPBindDN       string
	LDAPBindPassword string
}

// LoadAuthConfig reads auth-specific config from environment variables.
func LoadAuthConfig() *AuthConfig {
	return &AuthConfig{
		JWTExpire:        getEnv("JWT_EXPIRE", "15m"),
		RefreshExpire:    getEnv("REFRESH_EXPIRE", "168h"),
		PrivateKeyPath:   getEnv("JWT_PRIVATE_KEY_PATH", "keys/private.pem"),
		PublicKeyPath:    getEnv("JWT_PUBLIC_KEY_PATH", "keys/public.pem"),
		LDAPURL:          getEnv("LDAP_URL", "ldap://localhost:389"),
		LDAPBaseDN:       getEnv("LDAP_BASE_DN", "dc=example,dc=com"),
		LDAPBindDN:       getEnv("LDAP_BIND_DN", "cn=admin,dc=example,dc=com"),
		LDAPBindPassword: getEnv("LDAP_BIND_PASSWORD", "admin"),
	}
}

func getEnv(key, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}
