package dto

// AuthPayload is the request body for register and login.
type AuthPayload struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// RefreshPayload is the request body for token refresh.
type RefreshPayload struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// LDAPLoginPayload is the request body for LDAP login.
type LDAPLoginPayload struct {
	Email    string `json:"email" binding:"required"` // Could be a username or email
	Password string `json:"password" binding:"required,min=1"`
}
