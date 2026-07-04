package ldap

import (
	"context"
	"fmt"

	"github.com/f0bima/go-auth-starter/internal/feature/auth/infrastructure/config"
	"github.com/go-ldap/ldap/v3"
)

// Client represents the LDAP client.
type Client struct {
	config *config.AuthConfig
}

// NewLDAPClient creates a new LDAP client.
func NewLDAPClient(cfg *config.AuthConfig) *Client {
	return &Client{
		config: cfg,
	}
}

// Authenticate verifies the user's credentials against the LDAP server.
func (c *Client) Authenticate(ctx context.Context, username, password string) (bool, error) {
	// Connect to LDAP
	l, err := ldap.DialURL(c.config.LDAPURL)
	if err != nil {
		return false, fmt.Errorf("failed to connect to LDAP: %w", err)
	}
	defer l.Close()

	// Bind as admin to search for the user
	err = l.Bind(c.config.LDAPBindDN, c.config.LDAPBindPassword)
	if err != nil {
		return false, fmt.Errorf("failed to bind to LDAP: %w", err)
	}

	// Search for the user
	searchRequest := ldap.NewSearchRequest(
		c.config.LDAPBaseDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf("(&(objectClass=*)(|(uid=%[1]s)(cn=%[1]s)(mail=%[1]s)))", ldap.EscapeFilter(username)),
		[]string{"dn"},
		nil,
	)

	sr, err := l.Search(searchRequest)
	if err != nil {
		return false, fmt.Errorf("failed to search LDAP: %w", err)
	}

	if len(sr.Entries) != 1 {
		return false, fmt.Errorf("user does not exist or too many entries returned")
	}

	userDN := sr.Entries[0].DN

	// Bind as the user to verify their password
	err = l.Bind(userDN, password)
	if err != nil {
		return false, nil // Invalid password
	}

	return true, nil
}
