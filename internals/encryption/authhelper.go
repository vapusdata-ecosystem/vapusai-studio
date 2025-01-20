package encrytion

import (
	jwt "github.com/golang-jwt/jwt/v5"
)

const (
	ClaimUserIdKey            = "userId"
	ClaimOrganizationKey      = "organization"
	ClaimOrganizationRolesKey = "organizationRole"
	ClaimRoleScopeKey         = "roleScope"
	ClaimRoleKey              = "role"
	ClaimAccountKey           = "account"
	ClaimUserNameKey          = "userName"
)

func ParseUnValidatedJWT(tokenString string) (map[string]interface{}, error) {
	claims := jwt.MapClaims{}
	newParser := jwt.NewParser()
	_, _, err := newParser.ParseUnverified(tokenString, claims)
	if err != nil {
		return nil, err
	}
	return claims, nil
}

func FlatATClaims(claims *VapusDataStudioAccessClaims, separator string) map[string]string {
	if err := claims.Scope.Validate(); err != nil {
		encryptLogger.Err(err).Msg("error while validating vapusdata platform access claims")
		return nil
	}
	if claims != nil {

		val := map[string]string{ClaimUserIdKey: claims.Scope.UserId, ClaimRoleScopeKey: claims.Scope.RoleScope, ClaimRoleKey: claims.Scope.StudioRole, ClaimAccountKey: claims.Scope.AccountId}
		if claims.Scope.OrganizationId != "" {
			val[ClaimOrganizationKey] = claims.Scope.OrganizationId
			val[ClaimOrganizationRolesKey] = claims.Scope.OrganizationRole
		}

		return val
	}
	return nil
}
