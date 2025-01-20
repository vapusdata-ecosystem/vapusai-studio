package dmstores

import (
	"context"
	"fmt"
	"log"

	"github.com/vapusdata-oss/aistudio/aistudio/utils"
	models "github.com/vapusdata-oss/aistudio/core/models"
	svcops "github.com/vapusdata-oss/aistudio/core/serviceops"
	mpb "github.com/vapusdata-oss/apis/protos/models/v1alpha1"
)

func (ds *DMStore) ConfigureAIGuardrails(ctx context.Context, obj *models.AIGuardrails, ctxClaim map[string]string) error {
	_, err := ds.Db.PostgresClient.DB.NewInsert().Model(obj).ModelTableExpr(svcops.VapusGuardrailsTable).Exec(ctx)
	if err != nil {
		logger.Err(err).Ctx(ctx).Msg("error while saving ai guardrail to datastore")
		return err
	}
	go func() {
		mCtx := context.TODO()
		_ = svcops.AddResourceArn(mCtx, ds.Db, &models.VapusResourceArn{
			ResourceId:   obj.VapusID,
			ResourceName: mpb.RequestObjects_VAPUS_AIAGENTS.String(),
			VapusBase: models.VapusBase{
				Editors: obj.Editors,
			},
		}, logger, ctxClaim)
	}()
	return nil
}

func (ds *DMStore) SaveGuardrailThread(ctx context.Context, obj *models.AIGuardrails, ctxClaim map[string]string) error {
	_, err := ds.Db.PostgresClient.DB.NewInsert().Model(obj).ModelTableExpr(svcops.VapusGuardrailsTable).Exec(ctx)
	if err != nil {
		logger.Err(err).Ctx(ctx).Msg("error while saving ai guardrail thread to datastore")
		return err
	}
	return nil
}

func (ds *DMStore) ListAIGuardrails(ctx context.Context, condition string, ctxClaim map[string]string) ([]*models.AIGuardrails, error) {
	if ctxClaim == nil {
		condition = GetAccountFilter(ctxClaim, condition)
	}
	if condition == "" {
		condition = "deleted_at IS NULL AND status = 'ACTIVE'"
	}
	query := fmt.Sprintf("SELECT * FROM %s WHERE %s", svcops.VapusGuardrailsTable, condition)
	log.Println(query)
	result := []*models.AIGuardrails{}
	err := ds.Db.PostgresClient.SelectInApp(ctx, &query, &result)
	if err != nil {
		logger.Err(err).Ctx(ctx).Msg("error while fetching ai guardrail list from datastore")
		return nil, err
	}
	return result, nil
}

func (ds *DMStore) GetAIGuardrail(ctx context.Context, iden string, ctxClaim map[string]string) (*models.AIGuardrails, error) {
	if iden == "" {
		return nil, utils.ErrAIModelNode404
	}
	result := []*models.AIGuardrails{}
	query := fmt.Sprintf("SELECT * FROM %s WHERE %s", svcops.VapusGuardrailsTable, GetByIdFilter("", iden, ctxClaim))
	log.Println(query)
	err := ds.Db.PostgresClient.SelectInApp(ctx, &query, &result)
	if err != nil || len(result) == 0 {
		logger.Err(err).Ctx(ctx).Msg("error while getting ai guardrail details from datastore")
		return nil, err
	}
	return result[0], err
}

func (ds *DMStore) PutAIGuardrails(ctx context.Context, obj *models.AIGuardrails, ctxClaim map[string]string) error {
	_, err := ds.Db.PostgresClient.DB.NewUpdate().Model(obj).ModelTableExpr(svcops.VapusGuardrailsTable).Where(VapusIdFilter(), obj.VapusID).Exec(ctx)
	if err != nil {
		logger.Err(err).Ctx(ctx).Msg("error while updating ai guardrail to datastore")
		return err
	}
	return nil
}
