package ldapauthn

import (
	"context"
	"crypto/tls"
	"fmt"
	"log/slog"
	"strings"

	"github.com/SigNoz/signoz/pkg/authn"
	"github.com/SigNoz/signoz/pkg/errors"
	"github.com/SigNoz/signoz/pkg/types"
	"github.com/SigNoz/signoz/pkg/types/authtypes"
	"github.com/SigNoz/signoz/pkg/valuer"
	ldap "github.com/go-ldap/ldap/v3"
)

var _ authn.PasswordAuthN = (*AuthN)(nil)

type AuthN struct {
	store           authtypes.AuthNStore
	authDomainStore authtypes.AuthDomainStore
}

func New(store authtypes.AuthNStore, authDomainStore authtypes.AuthDomainStore) *AuthN {
	return &AuthN{
		store:           store,
		authDomainStore: authDomainStore,
	}
}

func (a *AuthN) Authenticate(ctx context.Context, email string, password string, orgID valuer.UUID) (*authtypes.Identity, error) {
	username := email
	if idx := strings.Index(email, "@"); idx != -1 {
		username = email[:idx]
	}

	emailDomain := ""
	if idx := strings.Index(email, "@"); idx != -1 {
		emailDomain = email[idx+1:]
	}

	if emailDomain == "" {
		return nil, errors.New(errors.TypeInvalidInput, errors.CodeInvalidInput, "email must contain a domain")
	}

	// Look up the LDAP config from the auth domain for this email domain
	ldapConfig, err := a.getLDAPConfig(ctx, emailDomain, orgID)
	if err != nil {
		return nil, err
	}

	// Connect to LDAP server
	conn, err := dialLDAP(ldapConfig)
	if err != nil {
		slog.ErrorContext(ctx, "failed to connect to ldap server",
			slog.String("host", ldapConfig.Host),
			slog.Int("port", ldapConfig.GetPort()),
			slog.Any("error", err),
		)
		return nil, errors.New(errors.TypeInternal, errors.CodeInternal, "ldap authentication failed: could not connect to server")
	}
	defer conn.Close()

	// Try bind with each configured domain (e.g. SETYAZILIM\username, SETSOFTWARE\username)
	var bindErr error
	for _, domain := range ldapConfig.Domains {
		bindDN := domain + `\` + username
		bindErr = conn.Bind(bindDN, password)
		if bindErr == nil {
			break
		}
	}

	if bindErr != nil {
		slog.InfoContext(ctx, "ldap bind failed for all domains",
			slog.String("username", username),
		)
		return nil, errors.New(errors.TypeUnauthenticated, types.ErrCodeIncorrectPassword, "invalid credentials")
	}

	// LDAP bind succeeded. Check if user exists locally.
	user, _, userRoles, err := a.store.GetActiveUserAndFactorPasswordByEmailAndOrgID(ctx, email, orgID)
	if err != nil {
		if errors.Ast(err, errors.TypeNotFound) {
			// User authenticated via LDAP but doesn't exist locally.
			// Return identity with zero UserID to signal auto-provisioning is needed.
			return authtypes.NewIdentity(valuer.UUID{}, orgID, valuer.MustNewEmail(email), types.RoleViewer, authtypes.IdentNProviderTokenizer), nil
		}
		return nil, err
	}

	if len(userRoles) == 0 {
		return nil, errors.New(errors.TypeUnexpected, authtypes.ErrCodeUserRolesNotFound, "no user roles entries found")
	}

	role := authtypes.SigNozManagedRoleToExistingLegacyRole[userRoles[0].Role.Name]
	return authtypes.NewIdentity(user.ID, orgID, user.Email, role, authtypes.IdentNProviderTokenizer), nil
}

func (a *AuthN) getLDAPConfig(ctx context.Context, emailDomain string, orgID valuer.UUID) (*authtypes.LDAPConfig, error) {
	authDomain, err := a.authDomainStore.GetByNameAndOrgID(ctx, emailDomain, orgID)
	if err != nil {
		return nil, errors.New(errors.TypeUnauthenticated, errors.CodeUnauthenticated, "ldap authentication failed: could not find auth domain configuration")
	}

	ldapConfig := authDomain.AuthDomainConfig().LDAP
	if ldapConfig == nil {
		return nil, errors.New(errors.TypeUnauthenticated, errors.CodeUnauthenticated, "ldap is not configured for this domain")
	}

	return ldapConfig, nil
}

func dialLDAP(config *authtypes.LDAPConfig) (*ldap.Conn, error) {
	address := fmt.Sprintf("%s:%d", config.Host, config.GetPort())

	if config.UseTLS {
		return ldap.DialTLS("tcp", address, &tls.Config{
			InsecureSkipVerify: config.InsecureSkipVerify, //nolint:gosec // configurable by admin
		})
	}

	return ldap.Dial("tcp", address)
}
