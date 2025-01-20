package encrytion

import "errors"

var (
	ErrParsingJWT                = errors.New("error while parsing JWT")
	ErrInvalidJWT                = errors.New("invalid JWT")
	ErrInvalidJWTClaims          = errors.New("invalid claims in the auth token")
	ErrInvalidUserAuthentication = errors.New("error while validating user's authentication, unautorized access")
	ErrOnlyPublicJWTKey          = errors.New("system has only public key is provided, cannot generate JWT")
	ErrUserOrg404                = errors.New("user does not have access to the organization")
)
