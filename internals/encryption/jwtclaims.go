package encrytion

import (
	"github.com/go-playground/validator/v10"
	jwt "github.com/golang-jwt/jwt/v5"
)

// type VapusDataStudioAccessClaims struct {
// 	Scope   string                 `json:"scope"`
// 	Profile map[string]interface{} `json:"profile"`
// 	jwt.RegisteredClaims
// }

type VapusDataResources struct {
	Name        string   `json:"name"`
	Identifiers []string `json:"identifiers"`
	Role        []string `json:"role"`
}

// TODO : add validators for below 2 structs

type StudioScope struct {
	UserId           string `validate:"required" json:"userId"`
	AccountId        string `validate:"required" json:"accountId"`
	OrganizationId   string `validate:"required" json:"organizationId"`
	OrganizationRole string `validate:"required" json:"organizationRole"`
	RoleScope        string `validate:"required" json:"roleScope"`
	StudioRole       string `validate:"required" json:"platformRole"`
	// MarketplaceId 	 string `validate:"required" json:"marketplaceId"`
}

type VapusDataStudioAccessClaims struct {
	jwt.RegisteredClaims
	Scope *StudioScope `validate:"required" json:"scope"`
}

func (x *StudioScope) Validate() error {
	validator := validator.New()
	return validator.Struct(x)
}

type VapusDataStudioRefreshTokenClaims struct {
	jwt.RegisteredClaims
	UserId         string `validate:"required" json:"userId"`
	OrganizationId string `validate:"required" json:"organizationId"`
}

func (x *VapusDataStudioRefreshTokenClaims) Validate() error {
	validator := validator.New()
	return validator.Struct(x)
}
