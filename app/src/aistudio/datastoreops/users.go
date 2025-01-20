package dmstores

import (
	"context"
	"fmt"
	"log"
	"time"

	pkgs "github.com/vapusdata-oss/aistudio/aistudio/pkgs"
	utils "github.com/vapusdata-oss/aistudio/aistudio/utils"
	encrytion "github.com/vapusdata-oss/aistudio/core/encryption"
	dmerrors "github.com/vapusdata-oss/aistudio/core/errors"
	"github.com/vapusdata-oss/aistudio/core/models"
	svcops "github.com/vapusdata-oss/aistudio/core/serviceops"
	dmutils "github.com/vapusdata-oss/aistudio/core/utils"
	mpb "github.com/vapusdata-oss/apis/protos/models/v1alpha1"
)

func (ds *DMStore) GetOrUpdateUser(ctx context.Context, lu *pkgs.LocalUserM, createIfInvited bool, useDefaultOrganization bool, ctxClaim map[string]string) (*models.Users, error) {
	result := []*models.Users{}
	var err error
	log.Println("Email", lu.Email)
	query := fmt.Sprintf("SELECT * FROM %s WHERE email = '%s'", svcops.UsersTable, lu.Email)
	err = ds.Db.PostgresClient.SelectInApp(ctx, &query, &result)
	if err != nil || len(result) == 0 {
		logger.Error().Msgf("User Not found with userEmail - %v, please get an invite.", lu.Email)
		return nil, dmerrors.DMError(utils.ErrUser404, err)
	}
	userObj := result[0]
	log.Println("UserObj", userObj)
	if userObj.Status == mpb.CommonStatus_INVITED.String() && createIfInvited && userObj.InviteExpiresOn > time.Now().Unix() {
		validOrganization := false
		domainId := ""
		if !useDefaultOrganization {
			for _, val := range userObj.OrganizationRoles {
				if val.OrganizationId == lu.Organization {
					validOrganization = true
					domainId = val.OrganizationId
				}
			}
		} else {
			if len(userObj.OrganizationRoles) > 0 {
				validOrganization = true
				domainId = userObj.OrganizationRoles[0].OrganizationId
			}
		}

		if !validOrganization {
			logger.Error().Msgf("error: domain %v is not attached to user %v", lu.Organization, lu.Email)
			return nil, utils.ErrUserOrganization404
		}
		userObj.Status = mpb.CommonStatus_ACTIVE.String()
		userObj.FirstName = lu.FirstName
		userObj.LastName = lu.LastName
		userObj.DisplayName = lu.DisplayName
		if userObj.Profile == nil {
			userObj.Profile = &models.UserProfile{}
		}
		userObj.Profile.Avatar = lu.ProfileImage
		userObj.SetUserId()
		userObj.PreSaveCreate(ctxClaim)
		userObj.SetDefaultOrganization(domainId)
		err = ds.PutUser(ctx, userObj, ctxClaim)
		if err != nil {
			logger.Error().Msgf("error: user %v is not attached to any domain", lu.Email)
			return nil, err
		}
		return userObj, nil
	}
	if len(userObj.OrganizationRoles) == 0 {
		logger.Error().Msgf("error: user %v is not attached to any domain", lu.Email)
		return nil, utils.ErrUserOrganization404
	}
	if userObj.Status == mpb.CommonStatus_ACTIVE.String() {
		return userObj, nil
	}
	return nil, utils.ErrUser404
}

func (ds *DMStore) CreateUser(ctx context.Context, lu *pkgs.LocalUserM, role []string, uo *models.Users, ctxClaim map[string]string) (*models.Users, error) {
	userObj := &models.Users{}
	if uo == nil {
		logger.Info().Msgf("Creating user object for '%v'", lu.Email)
		userObj.Email = lu.Email
		if userObj.Status == "" {
			userObj.Status = mpb.CommonStatus_ACTIVE.String()
		}
		userObj.FirstName = lu.FirstName
		userObj.DisplayName = lu.DisplayName
		userObj.LastName = lu.LastName
		userObj.StudioRoles = role
		userObj.InvitedType = mpb.UserInviteType_INVITE_ACCESS.String()
		userObj.InvitedOn = dmutils.GetEpochTime()
		userObj.OwnerAccount = ctxClaim[encrytion.ClaimAccountKey]
		if lu.Organization != "" {
			userObj.OrganizationRoles = []*models.UserOrganizationRole{
				{
					OrganizationId: lu.Organization,
					RoleArns:       lu.OrganizationRoles,
				},
			}
		}
		userObj.Profile = &models.UserProfile{}
		userObj.SetUserId()
		userObj.PreSaveCreate(ctxClaim)
		logger.Info().Msgf("New User object for '%v'", userObj)
	} else {
		userObj = uo
		logger.Info().Msgf("User object for '%v'", userObj)
	}
	userObj.SetAccountId(ctxClaim[encrytion.ClaimAccountKey])
	_, err := ds.Db.PostgresClient.DB.NewInsert().Model(userObj).ModelTableExpr(svcops.UsersTable).Exec(ctx)
	if err != nil {
		logger.Err(err).Ctx(ctx).Msg("error while saving datamarketplace in datastore")
		return nil, err
	}
	go func() {
		mCtx := context.TODO()
		_ = svcops.AddResourceArn(mCtx, ds.Db, &models.VapusResourceArn{
			ResourceId:   userObj.VapusID,
			ResourceName: "USER",
			VapusBase: models.VapusBase{
				Editors: []string{userObj.UserId},
			},
		}, logger, ctxClaim)
	}()
	return userObj, nil
}

func (ds *DMStore) UserInviteExists(ctx context.Context, userId string, ctxClaim map[string]string) bool {
	query := fmt.Sprintf("SELECT * FROM %s WHERE user_id = '%s'", svcops.UsersTable, userId)
	var user *models.Users
	err := ds.Db.PostgresClient.SelectInApp(ctx, &query, user)
	if err != nil || user == nil {
		logger.Info().Msgf("User Not found with userEmail - %v, please get an invite.", userId)
		return false
	}
	return true
}

func (ds *DMStore) LogStudioRTinfo(ctx context.Context, obj *models.RefreshTokenLog, ctxClaim map[string]string) error {
	obj.SetAccountId(ctxClaim[encrytion.ClaimAccountKey])
	_, err := ds.Db.PostgresClient.DB.NewInsert().Model(obj).ModelTableExpr(svcops.RefreshTokenLogsTable).Exec(ctx)
	if err != nil {
		logger.Err(err).Ctx(ctx).Msg("error while saving refresh token log in datastore")
		return err
	}
	ctx.Done()
	return nil
}

func (ds *DMStore) LogStudioJwtinfo(ctx context.Context, obj *models.JwtLog, ctxClaim map[string]string) error {
	obj.SetAccountId(ctxClaim[encrytion.ClaimAccountKey])
	_, err := ds.Db.PostgresClient.DB.NewInsert().Model(obj).ModelTableExpr(svcops.JwtLogsTable).Exec(ctx)
	if err != nil {
		logger.Err(err).Ctx(ctx).Msg("error while saving jwt log in datastore")
		return err
	}
	ctx.Done()
	return nil
}

func (ds *DMStore) GetStudioRTinfo(ctx context.Context, token string, ctxClaim map[string]string) (*models.RefreshTokenLog, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE token_hash = '%s'", svcops.RefreshTokenLogsTable, token)
	var user *models.RefreshTokenLog
	err := ds.Db.PostgresClient.SelectInApp(ctx, &query, user)
	if err != nil {
		logger.Info().Msg("error while getting refresh token log from datastore")
		return nil, err
	}
	return user, nil
}

func (ds *DMStore) PatchUser(ctx context.Context, userId string, data, conditions map[string]interface{}, ctxClaim map[string]string) error {
	// Convert the script query to JSON
	pq := ds.Db.PostgresClient.DB.NewUpdate().Model(&data).ModelTableExpr(svcops.UsersTable)

	for key, value := range conditions {
		pq = pq.Where(key, value)
	}
	_, err := pq.Exec(ctx)
	if err != nil {
		logger.Err(err).Ctx(ctx).Msg("error while patching user info in datastore")
		return err
	}
	logger.Debug().Ctx(ctx).Msgf("User '%v' successfully updated", userId)
	return nil
}

func (ds *DMStore) PutUser(ctx context.Context, obj *models.Users, ctxClaim map[string]string) error {
	_, err := ds.Db.PostgresClient.DB.NewUpdate().Model(obj).ModelTableExpr(svcops.UsersTable).Where("user_id = ?", obj.UserId).Exec(ctx)
	if err != nil {
		logger.Err(err).Ctx(ctx).Msg("error while updating user in datastore")
		return err
	}
	return nil
}

func (ds *DMStore) GetOrganizationUsers(ctx context.Context, domain string, ctxClaim map[string]string) ([]*models.Users, error) {
	result := []*models.Users{}
	query := fmt.Sprintf(`SELECT * FROM %s WHERE domain_roles @> '[{"domainId": "%s"}]'`, svcops.UsersTable, domain)
	err := ds.Db.PostgresClient.SelectInApp(ctx, &query, &result)
	if err != nil {
		logger.Err(err).Ctx(ctx).Msg("error while getting users from datastore")
		return result, err
	}
	return result, nil
}

func (ds *DMStore) ListUsers(ctx context.Context, condition string, ctxClaim map[string]string) ([]*models.Users, error) {
	result := make([]*models.Users, 0)
	query := fmt.Sprintf("SELECT * FROM %s WHERE %s", svcops.UsersTable, GetAccountFilter(ctxClaim, condition))
	err := ds.Db.PostgresClient.SelectInApp(ctx, &query, &result)
	if err != nil || len(result) == 0 {
		logger.Err(err).Ctx(ctx).Msg("error while getting users from datastore")
		return nil, err
	}
	return result, err
}

func (ds *DMStore) CountUsers(ctx context.Context, condition string, ctxClaim map[string]string) (int64, error) {
	var result int64
	condition = GetAccountFilter(ctxClaim, condition)
	query := fmt.Sprintf("SELECT count(*) FROM %s WHERE %s", svcops.UsersTable, condition)
	log.Println("query", query)
	err := ds.Db.PostgresClient.SelectInApp(ctx, &query, &result)
	if err != nil {
		logger.Err(err).Ctx(ctx).Msg("error while getting users count from datastore")
		return 0, err
	}
	log.Println("result--------->>>>>>>>>>>>>>>>", result)
	return result, err
}

func (ds *DMStore) CustomListUsers(ctx context.Context, fieldQuery, condition, postFilterForamtting string, ctxClaim map[string]string) ([]map[string]interface{}, error) {
	result := []map[string]interface{}{}
	condition = GetAccountFilter(ctxClaim, condition)
	query := fmt.Sprintf("SELECT %s FROM %s WHERE %s %s", fieldQuery, svcops.UsersTable, condition, postFilterForamtting)
	log.Println("query CustomListUsers --- >", query)
	err := ds.Db.PostgresClient.SelectInApp(ctx, &query, &result)
	if err != nil || len(result) == 0 {
		logger.Err(err).Ctx(ctx).Msg("error while getting dataproducts from datastore")
		return nil, err
	}
	return result, err
}

func (ds *DMStore) GetUser(ctx context.Context, userId string, ctxClaim map[string]string) (*models.Users, error) {
	result := make([]*models.Users, 0)
	query := fmt.Sprintf("SELECT * FROM %s WHERE %s", svcops.UsersTable, GetByIdFilter("user_id", userId, ctxClaim))
	err := ds.Db.PostgresClient.SelectInApp(ctx, &query, &result)
	if err != nil || len(result) == 0 {
		logger.Err(err).Ctx(ctx).Msg("error while getting users from datastore")
		return nil, utils.ErrUser404
	}
	return result[0], err
}
