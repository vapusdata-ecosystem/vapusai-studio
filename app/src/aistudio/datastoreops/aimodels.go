package dmstores

import (
	"context"
	"fmt"
	"log"

	"github.com/vapusdata-oss/aistudio/aistudio/utils"
	encrytion "github.com/vapusdata-oss/aistudio/core/encryption"
	models "github.com/vapusdata-oss/aistudio/core/models"
	svcops "github.com/vapusdata-oss/aistudio/core/serviceops"
	mpb "github.com/vapusdata-oss/apis/protos/models/v1alpha1"
)

func (ds *DMStore) ConfigureGetAIModelNode(ctx context.Context, obj *models.AIModelNode, ctxClaim map[string]string) error {
	_, err := ds.Db.PostgresClient.DB.NewInsert().Model(obj).ModelTableExpr(svcops.AIModelsNodesTable).Exec(ctx)
	if err != nil {
		logger.Err(err).Ctx(ctx).Msg("error while saving ai model node metadata in datastore")
		return err
	}
	go func() {
		mCtx := context.TODO()
		_ = svcops.AddResourceArn(mCtx, ds.Db, &models.VapusResourceArn{
			ResourceId:   obj.VapusID,
			ResourceName: mpb.RequestObjects_VAPUS_AIMODEL_NODES.String(),
			VapusBase: models.VapusBase{
				Editors: obj.Editors,
			},
		}, logger, ctxClaim)
	}()
	return nil
}

func (ds *DMStore) SaveAIInterfaceLog(ctx context.Context, obj *models.AIModelStudioLog, ctxClaim map[string]string) error {
	obj.PreSaveCreate(ctxClaim)
	_, err := ds.Db.PostgresClient.DB.NewInsert().
		Model(obj).
		ModelTableExpr(svcops.AIModelStudioLogsTable).
		Returning("id").
		Exec(ctx)
	if err != nil {
		logger.Err(err).Ctx(ctx).Msg("error while saving ai model studio log in datastore")
		return err
	}
	upQ := fmt.Sprintf("UPDATE %s SET search = to_tsvector('english', parsed_input) WHERE id = %d", svcops.AIModelStudioLogsTable, obj.ID)
	_, err = ds.Db.PostgresClient.DB.NewRaw(upQ).Exec(ctx)
	if err != nil {
		logger.Err(err).Ctx(ctx).Msgf("error while updating %s FTS field in datastore", svcops.AIModelStudioLogsTable)
		return err
	}
	return nil
}

func (ds *DMStore) ListAIInterfaceLogByUser(ctx context.Context, limit int, ctxClaim map[string]string) ([]*models.AIModelStudioLog, error) {
	resp := []*models.AIModelStudioLog{}
	query := fmt.Sprintf("SELECT * FROM %s WHERE created_by='%s' AND domain='%s' AND owner_account='%s' ORDER BY id DESC LIMIT %v",
		svcops.AIModelStudioLogsTable,
		ctxClaim[encrytion.ClaimUserIdKey],
		ctxClaim[encrytion.ClaimOrganizationKey],
		ctxClaim[encrytion.ClaimAccountKey],
		limit)
	log.Println("query---------------------->>>>>>>>>>>>>>>>>>>>>>>>>}}}}}}}}", query)
	err := ds.Db.PostgresClient.SelectInApp(ctx, &query, &resp)
	if err != nil {
		logger.Err(err).Ctx(ctx).Msg("error while fetching ai model studio logs from datastore")
		return nil, err
	}
	return resp, nil
}

func (ds *DMStore) PutAIModelNode(ctx context.Context, obj *models.AIModelNode) error {
	_, err := ds.Db.PostgresClient.DB.NewUpdate().Model(obj).ModelTableExpr(svcops.AIModelsNodesTable).Where(VapusIdFilter(), obj.VapusID).Exec(ctx)
	if err != nil {
		logger.Err(err).Ctx(ctx).Msg("error while updating ai model node metadata in datastore")
		return err
	}
	return nil
}

func (ds *DMStore) ListAIModelNodes(ctx context.Context, condition string, ctxClaim map[string]string) ([]*models.AIModelNode, error) {
	if ctxClaim == nil {
		condition = GetAccountFilter(ctxClaim, condition)
	}
	if condition == "" {
		condition = "deleted_at IS NULL AND status = 'ACTIVE'"
	}
	query := fmt.Sprintf("SELECT * FROM %s WHERE %s", svcops.AIModelsNodesTable, condition)
	log.Println("query---------------------->>>>>>>>>>>>>>>>>>>>>>>>>", query)
	result := []*models.AIModelNode{}
	err := ds.Db.PostgresClient.SelectInApp(ctx, &query, &result)
	// _, err = ds.Db.PostgresClient.DB.NewSelect().Model(&result).ModelTableExpr(svcops.AIModelsNodesTable).Where(condition).Exec(ctx)
	if err != nil {
		logger.Err(err).Ctx(ctx).Msg("error while fetching ai model nodes from datastore")
		return nil, err
	}
	return result, nil
}

func (ds *DMStore) GetAIModelNode(ctx context.Context, iden string, ctxClaim map[string]string) (*models.AIModelNode, error) {
	if iden == "" {
		return nil, utils.ErrAIModelNode404
	}
	result := []*models.AIModelNode{}
	query := fmt.Sprintf("SELECT * FROM %s WHERE %s", svcops.AIModelsNodesTable, GetByIdFilter("", iden, ctxClaim))
	// err := ds.Db.Select(ctx, &datasvcpkgs.QueryOpts{
	// 	DataCollection: svcops.DataSourcesTable,
	// 	RawQuery:       query,
	// }, &result, logger)
	log.Println("query---------------------->>>>>>>>>>>>>>>>>>>>>>>>>", query)
	err := ds.Db.PostgresClient.SelectInApp(ctx, &query, &result)
	if err != nil || len(result) == 0 {
		logger.Err(err).Ctx(ctx).Msg("error while getting ai models from datastore")
		return nil, err
	}
	return result[0], err
}
