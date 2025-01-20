package services

import (
	"context"
	"encoding/json"

	"github.com/rs/zerolog"
	dmstores "github.com/vapusdata-oss/aistudio/aistudio/datastoreops"
	pkgs "github.com/vapusdata-oss/aistudio/aistudio/pkgs"
	utils "github.com/vapusdata-oss/aistudio/aistudio/utils"
	encryption "github.com/vapusdata-oss/aistudio/core/encryption"
	dmerrors "github.com/vapusdata-oss/aistudio/core/errors"
	"github.com/vapusdata-oss/aistudio/core/models"
	svcops "github.com/vapusdata-oss/aistudio/core/serviceops"
	coreutils "github.com/vapusdata-oss/aistudio/core/utils"
	mpb "github.com/vapusdata-oss/apis/protos/models/v1alpha1"
)

func organizationConfigureTool(ctx context.Context, organization *models.Organization, dbStore *dmstores.DMStore, logger zerolog.Logger, ctxClaim map[string]string) (*models.Organization, error) {
	var err error
	organization.VapusID = ""
	organization.SetOrganizationId()
	organization.Status = mpb.CommonStatus_ACTIVE.String()
	organization.PreSaveCreate(ctxClaim)
	jwtSecretName := coreutils.GetSecretName("organization", organization.VapusID, "authJwtParams")
	if organization.AuthnJwtParams != nil {
		if !organization.AuthnJwtParams.IsAlreadyInSecretBs {
			jwtParam, err := setJWTAuthzParams(ctx, jwtSecretName, dbStore.SecretStore, false, organization.AuthnJwtParams)
			if err != nil {
				return nil, dmerrors.DMError(utils.ErrSavingOrganizationAuthJwt, err)
			}
			organization.AuthnJwtParams.Reset()
			organization.AuthnJwtParams = jwtParam
		}
	} else {
		jwtParam, err := setJWTAuthzParams(ctx, jwtSecretName, dbStore.SecretStore, true, nil)
		if err != nil {
			return nil, dmerrors.DMError(utils.ErrSavingOrganizationAuthJwt, err)
		}
		organization.AuthnJwtParams = jwtParam
	}

	organization.SecretPasscode = coreutils.GenerateRandomString(16)

	if organization.BackendSecretStorage == nil {
		// TO:DO add logic to store account bestorage
		// secName := getSecretName("organization", organization.VapusID, "organizationBeSecretStore")
	}

	_ = dbStore.BlobOps.CreateBucket(ctx, organization.VapusID)

	organization.Users = []string{ctxClaim[encryption.ClaimUserIdKey]}
	organization.Editors = organization.Users
	err = dbStore.ConfigureOrganization(ctx, organization, ctxClaim)
	if err != nil {
		logger.Err(err).Ctx(ctx).Msgf("error while configuring organization %v", organization)
		return nil, dmerrors.DMError(utils.ErrCreateOrganization, err) //nolint:wrapcheck
	}
	return organization, nil
}

func setJWTAuthzParams(ctx context.Context, secretName string, secretStoreClient *svcops.SecretStore, useStudio bool, jwtparam *models.JWTParams) (*models.JWTParams, error) {
	if jwtparam == nil || jwtparam.PublicJwtKey == "" || jwtparam.PrivateJwtKey == "" {
		useStudio = true
		helperLogger.Info().Msg("using platform jwt secrets because jwt secrets are not provided in request")
	}
	var err error
	if useStudio {
		err = secretStoreClient.WriteSecret(ctx, pkgs.SvcPackageManager.VapusJwtAuth.Opts, secretName)
		if err != nil {
			helperLogger.Err(err).Msgf("error while swapping default platform JWT keys for given resource - %v", secretName)
			return nil, err
		}
	} else {
		err = secretStoreClient.WriteSecret(ctx, &encryption.JWTAuthn{
			PublicJWTKey:     jwtparam.PublicJwtKey,
			PrivateJWTKey:    jwtparam.PrivateJwtKey,
			SigningAlgorithm: jwtparam.SigningAlgorithm,
		}, secretName)
		if err != nil {
			helperLogger.Err(err).Msgf("error while swapping JWT keys for given resource - %v", secretName)
			return nil, err
		}
	}
	return &models.JWTParams{
		VId:                 secretName,
		Name:                secretName,
		SigningAlgorithm:    pkgs.SvcPackageManager.VapusJwtAuth.Opts.SigningAlgorithm,
		IsAlreadyInSecretBs: true,
		Status:              mpb.CommonStatus_ACTIVE.String(),
	}, nil
}

func getsecretPassCode(resource, resourceId string) string {
	return coreutils.SlugifyBase(resource) + "_" + coreutils.SlugifyBase(resourceId)
}

func getOrganizationAuthn(ctx context.Context, organization *models.Organization, store *dmstores.DMStore, forSignValidation bool) (*encryption.JWTAuthn, error) {
	authnObj := organization.GetAuthnJwtParams()
	helperLogger.Info().Msgf("authnObj - %v", authnObj)
	secretStr, err := store.SecretStore.ReadSecret(ctx, authnObj.GetName())
	if err != nil {
		helperLogger.Err(err).Msgf("error while fetching organization authn secrets")
		return nil, err
	}
	helperLogger.Info().Msgf("authnObj - %v", secretStr)
	jwtParams := &encryption.JWTAuthn{}
	err = json.Unmarshal([]byte(secretStr), jwtParams)
	if err != nil {
		helperLogger.Err(err).Ctx(ctx).Msg("error while unmarshaling the organization jwt from secret store")
		return nil, err
	}

	if forSignValidation {
		jwtParams.PrivateJWTKey = ""
		jwtParams.ForPublicValidation = true
	}
	return jwtParams, nil
}
