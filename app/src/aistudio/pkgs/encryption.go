package pkgs

import (
	"fmt"
	"strings"
	"time"

	validator "github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	encryption "github.com/vapusdata-oss/aistudio/core/encryption"
	models "github.com/vapusdata-oss/aistudio/core/models"
)

var VapusAuth *encryption.VapusDataJwtAuthn

var JwtAuthnParams *encryption.JWTAuthn

func NewVapusAuth(params *encryption.JWTAuthn, validator *validator.Validate) (*encryption.VapusDataJwtAuthn, error) {
	obj, err := encryption.NewVapusDataJwtAuthn(params)
	if err != nil {
		pkgLogger.Err(err).Msg("Error while loading jwt authnetication config")
		return nil, err
	}
	if err := validator.Struct(obj.Opts); err != nil {
		pkgLogger.Err(err).Msg("Error while validating jwt config")
		return nil, err
	}

	return obj, nil
}

func InitAIStudioAuth(params *encryption.JWTAuthn, validator *validator.Validate) {
	var err error
	if VapusAuth == nil {
		VapusAuth, err = NewVapusAuth(params, validator)
		if err != nil {
			pkgLogger.Err(err).Msg("error initializing authn")
			panic(err)
		}
	}
}

func BuildVDPAClaim(userObj *models.Users, domainId string, validTill time.Time) (*encryption.VapusDataStudioAccessClaims, error) {
	var roleScope string
	var domainRoles string
	if domainId != "" {
		roleScope = encryption.JwtOrganizationScope
		r := userObj.GetOrganizationRole(domainId)
		if len(r) < 1 && domainId != "" {
			return nil, fmt.Errorf("user does not have any role in the organization")
		} else {
			for _, val := range r[0].RoleArns {
				if domainRoles == "" {
					domainRoles = val
				} else {
					domainRoles = domainRoles + encryption.JwtClaimRoleSeparator + val
				}
			}
		}
	} else {
		roleScope = encryption.JwtStudioScope
	}

	return &encryption.VapusDataStudioAccessClaims{
		Scope: &encryption.StudioScope{
			UserId:           userObj.UserId,
			OrganizationId:   domainId,
			StudioRole:       strings.Join(userObj.GetStudioRole(), ","),
			RoleScope:        roleScope,
			AccountId:        userObj.OwnerAccount,
			OrganizationRole: domainRoles,
		},
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   encryption.VapusStudioTokenSubject,
			Audience:  []string{NetworkConfigManager.ExternalURL},
			ExpiresAt: jwt.NewNumericDate(validTill), // configurable tokens
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now().Add(time.Second * 1)),
			Issuer:    NetworkConfigManager.ExternalURL,
		},
	}, nil
}

func BuildVDPRTClaim(userObj *models.Users, orgId string, validTill time.Time) (*encryption.VapusDataStudioRefreshTokenClaims, error) {
	return &encryption.VapusDataStudioRefreshTokenClaims{
		UserId:         userObj.UserId,
		OrganizationId: orgId,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   encryption.VapusStudioTokenSubject,
			Audience:  []string{NetworkConfigManager.ExternalURL},
			ExpiresAt: jwt.NewNumericDate(validTill), // configurable tokens
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now().Add(time.Second * 1)),
			Issuer:    NetworkConfigManager.ExternalURL,
		},
	}, nil
}
