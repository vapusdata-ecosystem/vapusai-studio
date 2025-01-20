package authn

import (
	"context"

	oidc "github.com/coreos/go-oidc/v3/oidc"
	dmerrors "github.com/vapusdata-oss/aistudio/core/errors"
	dmlogger "github.com/vapusdata-oss/aistudio/core/logger"
	utils "github.com/vapusdata-oss/aistudio/core/utils"
	oauth2 "golang.org/x/oauth2"
)

type OIDCAuthenticator struct {
	*oidc.Provider
	*oidc.IDTokenVerifier
	oauth2.Config
	*OIDCSecrets
}

// NewOIDCAuthenticator function to create a new OIDC authenticator object using the OIDC secrets provided.
// The function returns an error if the OIDC provider cannot be initialized.
// The function returns a pointer to the OIDCAuthenticator object.
func NewOIDCAuthenticator(opts *OIDCSecrets) (*OIDCAuthenticator, error) {
	_, opts.Organization = utils.TrailingSlash(opts.Organization, true, true)
	dmlogger.CoreLogger.Info().Msgf("Initializing OIDC Authenticator with organization: %v", opts.Organization)
	provider, err := oidc.NewProvider(context.Background(), opts.Organization)
	if err != nil {
		return nil, dmerrors.DMError(ErrOIDCProviderInit, err)
	}

	conf := oauth2.Config{
		ClientID:     opts.ClientID,
		Endpoint:     provider.Endpoint(),
		RedirectURL:  opts.CallbackURI,
		ClientSecret: opts.ClientSecret,
		Scopes:       []string{oidc.ScopeOpenID, "email", "profile"},
	}

	return &OIDCAuthenticator{
		Provider:    provider,
		Config:      conf,
		OIDCSecrets: opts,
	}, nil

}

// VerifyIDToken verifies that an *oauth2.Token is a valid *oidc.IDToken.
// This functions uses the OIDC authenticator object to verify the token.
// The function returns an error if the token is invalid.
// The function returns a pointer to the *oidc.IDToken object.
func (a *OIDCAuthenticator) VerifyIDToken(ctx context.Context, token *oauth2.Token) (*oidc.IDToken, error) {
	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok {
		return nil, dmerrors.DMError(ErrOIDCInvalidToken, nil)
	}

	// default settings
	oidcConfig := &oidc.Config{ClientID: a.OIDCSecrets.ClientID}

	return a.Verifier(oidcConfig).Verify(ctx, rawIDToken)
}
