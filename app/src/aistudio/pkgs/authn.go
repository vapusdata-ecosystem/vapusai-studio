package pkgs

import (
	"github.com/vapusdata-oss/aistudio/core/authn"
)

type AuthnService struct {
	*authn.Authenticator
	Auth, Callback string
}

var AuthnManager *AuthnService

var AuthnParams *authn.AuthnSecrets

func InitAuthnManager(params *authn.AuthnSecrets) {
	if AuthnManager == nil {
		AuthnManager = &AuthnService{
			Authenticator: SvcPackageManager.AuthnManager,
		}
	}
}
