package authtypes

import (
	"encoding/json"

	"github.com/SigNoz/signoz/pkg/errors"
	"github.com/SigNoz/signoz/pkg/valuer"
)

type LDAPConfig struct {
	// LDAP server host (e.g. "mail.setyazilim.com.tr")
	Host string `json:"host"`

	// LDAP server port (default 389 for LDAP, 636 for LDAPS)
	Port int `json:"port"`

	// Use LDAPS (TLS) instead of plain LDAP
	UseTLS bool `json:"useTLS"`

	// Skip TLS certificate verification (for self-signed certs)
	InsecureSkipVerify bool `json:"insecureSkipVerify"`

	// LDAP domain names used for bind (e.g. ["SETYAZILIM", "SETSOFTWARE"]).
	// Users authenticate as DOMAIN\username.
	// The first domain is used as the default.
	Domains []string `json:"domains"`

	// Base DN for user search (e.g. "DC=setyazilim,DC=com,DC=tr").
	// If empty, user search is skipped and only bind authentication is performed.
	BaseDN string `json:"baseDN"`

	// User search filter. Use %s as placeholder for username.
	// Default: "(sAMAccountName=%s)"
	UserSearchFilter string `json:"userSearchFilter"`

	// Attribute for user's email. Default: "mail"
	EmailAttribute string `json:"emailAttribute"`

	// Attribute for user's display name. Default: "displayName"
	NameAttribute string `json:"nameAttribute"`
}

func (c *LDAPConfig) GetPort() int {
	if c.Port == 0 {
		if c.UseTLS {
			return 636
		}
		return 389
	}
	return c.Port
}

func (c *LDAPConfig) GetUserSearchFilter() string {
	if c.UserSearchFilter == "" {
		return "(sAMAccountName=%s)"
	}
	return c.UserSearchFilter
}

func (c *LDAPConfig) GetEmailAttribute() string {
	if c.EmailAttribute == "" {
		return "mail"
	}
	return c.EmailAttribute
}

func (c *LDAPConfig) GetNameAttribute() string {
	if c.NameAttribute == "" {
		return "displayName"
	}
	return c.NameAttribute
}

type PostableLDAPSession struct {
	Email    valuer.Email `json:"email"`
	Password string       `json:"password"`
	OrgID    valuer.UUID  `json:"orgId"`
}

func (typ *PostableLDAPSession) UnmarshalJSON(data []byte) error {
	type Alias PostableLDAPSession
	var temp Alias

	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	if temp.Email.IsZero() {
		return errors.New(errors.TypeInvalidInput, errors.CodeInvalidInput, "email is required")
	}

	if temp.Password == "" {
		return errors.New(errors.TypeInvalidInput, errors.CodeInvalidInput, "password is required")
	}

	if temp.OrgID.IsZero() {
		return errors.New(errors.TypeInvalidInput, errors.CodeInvalidInput, "orgID is required")
	}

	*typ = PostableLDAPSession(temp)
	return nil
}
