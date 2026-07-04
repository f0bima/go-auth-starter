package service

import (
	"context"
)

// LDAPService defines the interface for LDAP authentication.
type LDAPService interface {
	Authenticate(ctx context.Context, username, password string) (bool, error)
}
