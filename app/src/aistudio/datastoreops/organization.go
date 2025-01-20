package dmstores

import (
	"context"
	"fmt"

	// mpb "github.com/vapusdata-oss/apis/protos/models/v1alpha1"
	utils "github.com/vapusdata-oss/aistudio/aistudio/utils"
	encrytion "github.com/vapusdata-oss/aistudio/core/encryption"
	dmerrors "github.com/vapusdata-oss/aistudio/core/errors"
	"github.com/vapusdata-oss/aistudio/core/models"
	serviceops "github.com/vapusdata-oss/aistudio/core/serviceops"
	svcops "github.com/vapusdata-oss/aistudio/core/serviceops"
	dmutils "github.com/vapusdata-oss/aistudio/core/utils"
	mpb "github.com/vapusdata-oss/apis/protos/models/v1alpha1"
)

// AddOrganization adds a new domain to the data domain
// It will create a new domain and attach it to the data domain
func (ds *DMStore) ConfigureOrganization(ctx context.Context, domain *models.Organization, ctxClaim map[string]string) error {
	domain.SetAccountId(ctxClaim[encrytion.ClaimAccountKey])
	if ds.Cacher != nil {
		_, err := ds.BeDataStore.Cacher.RedisClient.WrtiteData(ctx, domain.VapusID, dmutils.EMPTYSTR, domain)
		if err != nil {
			logger.Err(err).Ctx(ctx).Msg(utils.ErrOrganizationInitialization.Error())
			return err
		}
	}
	_, err := ds.Db.PostgresClient.DB.NewInsert().ModelTableExpr(serviceops.OrganizationTable).Model(domain).Exec(ctx)
	if err != nil {
		logger.Err(err).Ctx(ctx).Msg("error while saving datamarketplace in datastore")
		return err
	}
	go func() {
		mCtx := context.TODO()
		_ = svcops.AddResourceArn(mCtx, ds.Db, &models.VapusResourceArn{
			ResourceId:   domain.VapusID,
			ResourceName: mpb.RequestObjects_VAPUS_DOMAINS.String(),
			VapusBase: models.VapusBase{
				Editors: domain.Editors,
			},
		}, logger, ctxClaim)
	}()
	return nil
}

func (ds *DMStore) PatchOrganization(ctx context.Context, data, conditions map[string]interface{}, ctxClaim map[string]string) error {
	pq := ds.Db.PostgresClient.DB.NewUpdate().Model(&data).ModelTableExpr(serviceops.OrganizationTable)

	for key, value := range conditions {
		pq = pq.Where(key, value)
	}
	_, err := pq.Exec(ctx)
	if err != nil {
		logger.Err(err).Ctx(ctx).Msg("error while saving domain info in datastore")
		return err
	}
	return nil
}

func (ds *DMStore) PutOrganization(ctx context.Context, obj *models.Organization, ctxClaim map[string]string) error {
	_, err := ds.Db.PostgresClient.DB.NewUpdate().Model(obj).ModelTableExpr(serviceops.OrganizationTable).Where(VapusIdFilter(), obj.VapusID).Exec(ctx)
	if err != nil {
		logger.Err(err).Ctx(ctx).Msg("error while updating domain in datastore")
		return err
	}
	return nil
}

// GetOrganization gets the domain object from the data domain store based on the key identifier i.e. domainid
// GetMarketplaceOrganizations retrieves the domains associated with a given key from the DMStore.
// It returns a slice of *models.Organization, a map[string]interface{} for custom messages, and an error.
// The key parameter specifies the key to query the DMStore.
// If the retrieval is successful, the function returns the domains, custom messages, and a nil error.
// If an error occurs during retrieval, the function returns nil for domains, custom messages, and an error indicating the cause.
func (ds *DMStore) ListOrganizations(ctx context.Context, condition string, ctxClaim map[string]string) ([]*models.Organization, error) {
	result := []*models.Organization{}
	condition = GetAccountFilter(ctxClaim, condition)
	query := fmt.Sprintf("SELECT * FROM %s WHERE %s ORDER BY created_at DESC ", serviceops.OrganizationTable, condition)
	err := ds.Db.PostgresClient.SelectInApp(ctx, &query, &result)
	if err != nil {
		logger.Err(err).Ctx(ctx).Msgf("error while getting domain for the request")
		return nil, err
	}
	return result, nil
}

func (ds *DMStore) GetOrganization(ctx context.Context, iden string, ctxClaim map[string]string) (*models.Organization, error) {
	result := []*models.Organization{}
	// condition = GetAccountFilter(ctxClaim, condition)
	query := fmt.Sprintf("SELECT * FROM %s WHERE %s", serviceops.OrganizationTable, GetByIdFilter("", iden, ctxClaim))
	err := ds.Db.PostgresClient.SelectInApp(ctx, &query, &result)
	if err != nil || len(result) == 0 {
		logger.Err(err).Ctx(ctx).Msgf("error while getting domain for the request")
		return nil, dmerrors.DMError(utils.ErrInvalidOrganizationRequested, err)
	}
	return result[0], nil
}
