package nabhiksvc

import (
	"context"
	"log"
	"sync"

	"github.com/rs/zerolog"
	dmstores "github.com/vapusdata-oss/aistudio/aistudio/datastoreops"
	"github.com/vapusdata-oss/aistudio/aistudio/utils"
	aimodels "github.com/vapusdata-oss/aistudio/core/aistudio/providers"
	dmerrors "github.com/vapusdata-oss/aistudio/core/errors"
	"github.com/vapusdata-oss/aistudio/core/models"
)

type AIModelNodeConnectionPool struct {
	connectionPool map[string]aimodels.AIModelNodeInterface
	logger         zerolog.Logger
	dmStore        *dmstores.DMStore
	errPool        map[string]error
}

type aiPoolOpts func(*AIModelNodeConnectionPool)

var AIModelNodeConnectionPoolManager *AIModelNodeConnectionPool

func WithLogger(logger zerolog.Logger) aiPoolOpts {
	return func(a *AIModelNodeConnectionPool) {
		a.logger = logger
	}
}

func WithDMStore(dmStore *dmstores.DMStore) aiPoolOpts {
	return func(a *AIModelNodeConnectionPool) {
		a.dmStore = dmStore
	}
}

func InitAIModelNodeConnectionPool(opts ...aiPoolOpts) *AIModelNodeConnectionPool {
	if AIModelNodeConnectionPoolManager != nil {
		return AIModelNodeConnectionPoolManager
	}
	obj := &AIModelNodeConnectionPool{}
	for _, opt := range opts {
		opt(obj)
	}
	obj.connectionPool = map[string]aimodels.AIModelNodeInterface{}
	ctx, cancel := context.WithCancel(context.Background())
	obj.bootConnectionPool(ctx)
	defer cancel()
	if len(obj.errPool) > 0 {
		obj.logger.Error().Msg("error while booting connection pool")
	}
	log.Println(obj.connectionPool)
	return obj
}

func (a *AIModelNodeConnectionPool) AddConnection(model *models.AIModelNode, connection aimodels.AIModelNodeInterface) {
	a.connectionPool[model.VapusID] = connection
}

func (a *AIModelNodeConnectionPool) GetConnectionById(nodeId string) aimodels.AIModelNodeInterface {
	val, ok := a.connectionPool[nodeId]
	if !ok {
		a.logger.Info().Msgf("Connection not found in pool for %s", nodeId)
		return nil
	}
	return val
}

func (a *AIModelNodeConnectionPool) GetorSetConnection(model *models.AIModelNode, setIfNotPresent bool) (aimodels.AIModelNodeInterface, error) {
	val, ok := a.connectionPool[model.VapusID]
	if !ok && setIfNotPresent {
		a.logger.Info().Msgf("Connection not found in pool for %s , creating new connection", model.VapusID)
		if !setIfNotPresent {
			return nil, dmerrors.DMError(utils.ErrAIModelConn, nil)
		} else {
			a.createModelConnection(context.Background(), model)
			val, ok = a.connectionPool[model.VapusID]
			if !ok {
				return nil, dmerrors.DMError(utils.ErrAIModelConn, nil)
			}
			return val, nil
		}
	}
	return val, nil
}

func (a *AIModelNodeConnectionPool) RemoveConnection(model *models.AIModelNode) {
	delete(a.connectionPool, model.VapusID)
}

func (a *AIModelNodeConnectionPool) bootConnectionPool(ctx context.Context) error {
	result, err := a.dmStore.ListAIModelNodes(ctx, "status = 'ACTIVE' AND deleted_at IS NULL ORDER BY created_at DESC", nil)
	if err != nil {
		a.logger.Error().Err(err).Msg("error while fetching models from datastore")
		return err
	}
	var wg sync.WaitGroup
	for _, model := range result {
		wg.Add(1)
		go func(model *models.AIModelNode) {
			defer wg.Done()
			ctx := context.Background()
			a.createModelConnection(ctx, model)
		}(model)
	}
	wg.Wait()
	return nil
}

func (a *AIModelNodeConnectionPool) createModelConnection(ctx context.Context, model *models.AIModelNode) {
	netParam, err := a.dmStore.GetAIModelNodeNetworkParams(ctx, model)
	if err != nil {
		a.logger.Err(err).Ctx(ctx).Msg("error while getting network params for AI model")
		a.errPool[model.VapusID] = dmerrors.DMError(utils.ErrGetAIModelNetParams, err)
	}
	// secretStr, err := a.dmStore.ReadSecret(ctx, model.GetNetworkParams().SecretName)
	// gCred := &models.GenericCredentialModel{}
	// if err != nil {
	// 	a.logger.Err(err).Ctx(ctx).Msg("error while reading the secrets from secret store")
	// 	a.errPool[model.VapusID] = dmerrors.DMError(utils.ErrGetCredsSecret, err)
	// }
	// err = json.Unmarshal([]byte(secretStr), gCred)
	// if err != nil {
	// 	a.logger.Err(err).Ctx(ctx).Msg("error while unmarshaling the secrets from secret store")
	// 	a.errPool[model.VapusID] = dmerrors.DMError(utils.ErrGetCredsSecret, err)
	// }
	// // creds := &models.GenericCredentialModel{}
	// // model.GetNetworkParams().Credentials = creds.ConvertFromPb(gCred)
	model.NetworkParams = netParam
	conn, err := aimodels.NewAIModelNode(aimodels.WithAIModelNode(model), aimodels.WithLogger(a.logger))
	if err != nil {
		a.logger.Err(err).Ctx(ctx).Msg("error while creating model connection")
		a.errPool[model.VapusID] = dmerrors.DMError(utils.ErrAIModelConn, err)
	}
	a.logger.Info().Ctx(ctx).Msgf("Connection created for AI model %s", model.VapusID)
	a.AddConnection(model, conn)
}
