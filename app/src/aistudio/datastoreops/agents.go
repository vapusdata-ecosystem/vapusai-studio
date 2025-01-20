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

func (ds *DMStore) ConfigureAIAgents(ctx context.Context, obj *models.VapusAIAgent, ctxClaim map[string]string) error {
	_, err := ds.Db.PostgresClient.DB.NewInsert().Model(obj).ModelTableExpr(svcops.VapusAIAgentsTable).Exec(ctx)
	if err != nil {
		logger.Err(err).Ctx(ctx).Msg("error while saving ai agent to datastore")
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

func (ds *DMStore) SaveAIAgentThread(ctx context.Context, obj *models.AIAgentThread, ctxClaim map[string]string) error {
	_, err := ds.Db.PostgresClient.DB.NewInsert().Model(obj).ModelTableExpr(svcops.AIAgentThreadsTable).Exec(ctx)
	if err != nil {
		logger.Err(err).Ctx(ctx).Msg("error while saving ai agent thread to datastore")
		return err
	}
	return nil
}

func (ds *DMStore) ListAIAgents(ctx context.Context, condition string, ctxClaim map[string]string) ([]*models.VapusAIAgent, error) {
	condition = GetAccountFilter(ctxClaim, condition)
	query := fmt.Sprintf("SELECT * FROM %s WHERE %s", svcops.VapusAIAgentsTable, condition)
	result := []*models.VapusAIAgent{}
	err := ds.Db.PostgresClient.SelectInApp(ctx, &query, &result)
	if err != nil {
		logger.Err(err).Ctx(ctx).Msg("error while fetching ai agent list from datastore")
		return nil, err
	}
	return result, nil
}

func (ds *DMStore) GetAIAgent(ctx context.Context, iden string, ctxClaim map[string]string) (*models.VapusAIAgent, error) {
	if iden == "" {
		return nil, utils.ErrAIModelNode404
	}
	result := []*models.VapusAIAgent{}
	query := fmt.Sprintf("SELECT * FROM %s WHERE %s", svcops.VapusAIAgentsTable, GetByIdFilter("", iden, ctxClaim))
	log.Println(query)
	err := ds.Db.PostgresClient.SelectInApp(ctx, &query, &result)
	if err != nil || len(result) == 0 {
		logger.Err(err).Ctx(ctx).Msg("error while getting ai agent details from datastore")
		return nil, err
	}
	return result[0], err
}

func (ds *DMStore) PutAIAgents(ctx context.Context, obj *models.VapusAIAgent, ctxClaim map[string]string) error {
	_, err := ds.Db.PostgresClient.DB.NewUpdate().Model(obj).ModelTableExpr(svcops.VapusAIAgentsTable).Where(VapusIdFilter(), obj.VapusID).Exec(ctx)
	if err != nil {
		logger.Err(err).Ctx(ctx).Msg("error while updating ai agent to datastore")
		return err
	}
	return nil
}
