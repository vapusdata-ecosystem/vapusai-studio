package dmstores

import (
	"context"
	"fmt"

	"github.com/vapusdata-oss/aistudio/aistudio/utils"
	models "github.com/vapusdata-oss/aistudio/core/models"
	svcops "github.com/vapusdata-oss/aistudio/core/serviceops"
	mpb "github.com/vapusdata-oss/apis/protos/models/v1alpha1"
)

func (ds *DMStore) ConfigureAIPrompts(ctx context.Context, obj *models.AIModelPrompt, ctxClaim map[string]string) error {
	_, err := ds.Db.PostgresClient.DB.NewInsert().Model(obj).ModelTableExpr(svcops.AIModelPromptTable).Exec(ctx)
	if err != nil {
		logger.Err(err).Ctx(ctx).Msg("error while saving ai model prompt to datastore")
		return err
	}
	go func() {
		mCtx := context.TODO()
		_ = svcops.AddResourceArn(mCtx, ds.Db, &models.VapusResourceArn{
			ResourceId:   obj.VapusID,
			ResourceName: mpb.RequestObjects_VAPUS_AIPROMPTS.String(),
			VapusBase: models.VapusBase{
				Editors: obj.Editors,
			},
		}, logger, ctxClaim)
	}()
	return nil
}

func (ds *DMStore) ListAIPrompts(ctx context.Context, condition string, ctxClaim map[string]string) ([]*models.AIModelPrompt, error) {
	condition = GetAccountFilter(ctxClaim, condition)
	query := fmt.Sprintf("SELECT * FROM %s WHERE %s", svcops.AIModelPromptTable, condition)
	result := []*models.AIModelPrompt{}
	err := ds.Db.PostgresClient.SelectInApp(ctx, &query, &result)
	if err != nil {
		logger.Err(err).Ctx(ctx).Msg("error while fetching ai model prompt list from datastore")
		return nil, err
	}
	return result, nil
}

func (ds *DMStore) GetAIPrompt(ctx context.Context, iden string, ctxClaim map[string]string) (*models.AIModelPrompt, error) {
	if iden == "" {
		return nil, utils.ErrAIModelNode404
	}
	result := []*models.AIModelPrompt{}
	query := fmt.Sprintf("SELECT * FROM %s WHERE %s", svcops.AIModelPromptTable, GetByIdFilter("", iden, ctxClaim))
	err := ds.Db.PostgresClient.SelectInApp(ctx, &query, &result)
	if err != nil || len(result) == 0 {
		logger.Err(err).Ctx(ctx).Msg("error while getting ai model prompt details from datastore")
		return nil, err
	}
	return result[0], err
}

func (ds *DMStore) PutAIPrompts(ctx context.Context, obj *models.AIModelPrompt, ctxClaim map[string]string) error {
	_, err := ds.Db.PostgresClient.DB.NewUpdate().Model(obj).ModelTableExpr(svcops.AIModelPromptTable).Where(VapusIdFilter(), obj.VapusID).Exec(ctx)
	if err != nil {
		logger.Err(err).Ctx(ctx).Msg("error while updating ai model prompt to datastore")
		return err
	}
	return nil
}
