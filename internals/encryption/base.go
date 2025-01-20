package encrytion

import (
	"context"
	"encoding/json"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog"
	dmlogger "github.com/vapusdata-oss/aistudio/core/logger"
	utils "github.com/vapusdata-oss/aistudio/core/utils"
	mpb "github.com/vapusdata-oss/apis/protos/models/v1alpha1"
	"google.golang.org/grpc/metadata"
)

var encryptLogger zerolog.Logger

// TO:DO Make it more generic to handle different type of claims
type JwtAuthService interface {
	GenerateVDPAJWT(claims *VapusDataStudioAccessClaims) (string, error)
	GenerateVDPARefreshJWT(claims *VapusDataStudioRefreshTokenClaims) (string, error)
	ParseAndValidateAT(tokenString string) (*VapusDataStudioAccessClaims, error)
	ValidateAccessToken(tokenString string) (map[string]string, error)
}

type JWTAuthn struct {
	PublicJWTKey        string `validate:"required" yaml:"publicJwtKey" json:"publicJwtKey"`
	PrivateJWTKey       string `yaml:"privateJwtKey" json:"privateJwtKey"`
	SigningAlgorithm    string `validate:"required" yaml:"signingAlgorithm" json:"signingAlgorithm"`
	ForPublicValidation bool   `default:"false" yaml:"forPublicValidation" json:"forPublicValidation"`
	TokenIssuer         string `default:"VapusData" yaml:"tokenIssuer" json:"tokenIssuer"`
	TokenAudience       string `default:"*.vapusdata.com" yaml:"tokenAudience" json:"tokenAudience"`
}

type jwtAuthOpts func(jo *JWTAuthn)

type VapusDataJwtAuthn struct {
	Opts *JWTAuthn
	JwtAuthService
}

var JwtTokenIssuer = "VapusData"
var JwtTokenAudience = "*.vapusdata.com"
var VapusStudioTokenSubject = "VapusData access token"
var JwtOrganizationScope = "OrganizationScope"
var JwtStudioScope = "platformScope"
var JwtCtxClaimKey = "vapusStudioJwtClaim"
var JwtDPCtxClaimKey = "vapusStudioJwtClaim"
var JwtClaimRoleSeparator = "|"
var JWTParser *jwt.Parser

func SetCtxClaim(ctx context.Context, claim map[string]string) context.Context {
	return context.WithValue(ctx, JwtCtxClaimKey, claim)
}

func GetCtxClaim(ctx context.Context) (map[string]string, bool) {
	val, ok := ctx.Value(JwtCtxClaimKey).(map[string]string)
	if !ok {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, ok
		}
		val = make(map[string]string)
		if len(md.Get(JwtDPCtxClaimKey)) != 1 {
			return nil, false
		}
		strval := md.Get(JwtDPCtxClaimKey)[0]
		err := json.Unmarshal([]byte(strval), &val)
		if err != nil {
			return nil, false
		}
		return val, true

	}
	return val, ok
}

func GetDPtxClaim(ctx context.Context) (map[string]string, bool) {
	val, ok := ctx.Value(JwtDPCtxClaimKey).(map[string]string)
	return val, ok
}

func NewVapusDataJwtAuthnWithConfig(path string) (*VapusDataJwtAuthn, error) {
	encryptLogger = dmlogger.GetSubDMLogger(dmlogger.CoreLogger, "pkgs", "encryption")
	jwtAuthnSecrets, err := LoadJwtAuthnSecrets(path)
	if err != nil {
		return nil, err
	}
	return NewVapusDataJwtAuthn(jwtAuthnSecrets)
}

func LoadJwtAuthnSecrets(path string) (*JWTAuthn, error) {
	cf, err := utils.ReadBasicConfig(utils.GetConfFileType(path), path, &JWTAuthn{})
	if err != nil {
		encryptLogger.Info().Msgf("Error loading jwt authn secrets: %v", err)
		return nil, err
	}
	return cf.(*JWTAuthn), err
}

func NewVapusDataJwtAuthn(opts *JWTAuthn) (*VapusDataJwtAuthn, error) {
	encryptLogger = dmlogger.GetSubDMLogger(dmlogger.CoreLogger, "pkgs", "encryption")
	obj := &VapusDataJwtAuthn{
		Opts: opts,
	}
	JWTParser = jwt.NewParser(jwt.WithLeeway(2 * time.Second))
	switch opts.SigningAlgorithm {
	case mpb.EncryptionAlgo_ECDSA.String():
		val, err := NewECDSAJwtAuthn(opts)
		if err != nil {
			return nil, err
		}
		obj.JwtAuthService = val.(*ECDSAManager)
		return obj, nil
	case mpb.EncryptionAlgo_RSA.String():
		val, err := NewRSAJwtAuthn(opts)
		if err != nil {
			return nil, err
		}
		obj.JwtAuthService = val.(*RSAManager)
		return obj, nil
	default:
		return nil, ErrInvalidJWT
	}

}
