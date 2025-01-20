package services

import (
	"context"
	"slices"

	"github.com/rs/zerolog"
	"github.com/vapusdata-oss/aistudio/aistudio/nabhiksvc"
	"github.com/vapusdata-oss/aistudio/aistudio/pkgs"
	"github.com/vapusdata-oss/aistudio/aistudio/utils"
	"github.com/vapusdata-oss/aistudio/core/aistudio/prompts"
	aimodels "github.com/vapusdata-oss/aistudio/core/aistudio/providers"
	encryption "github.com/vapusdata-oss/aistudio/core/encryption"
	dmerrors "github.com/vapusdata-oss/aistudio/core/errors"
	models "github.com/vapusdata-oss/aistudio/core/models"
	coreutils "github.com/vapusdata-oss/aistudio/core/utils"
	mpb "github.com/vapusdata-oss/apis/protos/models/v1alpha1"
	pb "github.com/vapusdata-oss/apis/protos/vapusai-studio/v1alpha1"
)

type EmbeddingAgent struct {
	*models.VapusInterfaceAgentBase
	embeddings *mpb.Embeddings
	request    *pb.EmbeddingsInterface
	logger     zerolog.Logger
	modelNode  *models.AIModelNode
	*StudioServices
}

func (s *StudioServices) NewEmbeddingAgent(ctx context.Context,
	request *pb.EmbeddingsInterface) (*EmbeddingAgent, error) {
	var err error
	modelsNodeId := ""
	vapusStudioClaim, ok := encryption.GetCtxClaim(ctx)
	if !ok {
		s.logger.Error().Ctx(ctx).Msg("error while getting claim metadata from context")
		return nil, dmerrors.DMError(encryption.ErrInvalidJWTClaims, nil)
	}

	modelsNodeId = request.GetModelNodeId()
	result, err := s.DMStore.GetAIModelNode(ctx, modelsNodeId, vapusStudioClaim)
	if err != nil || result == nil {
		s.logger.Error().Err(err).Msg("error while getting model node")
		return nil, dmerrors.DMError(utils.ErrAIModelNode404, nil)
	}

	agent := &EmbeddingAgent{
		VapusInterfaceAgentBase: &models.VapusInterfaceAgentBase{
			AgentId:  coreutils.GetUUID(),
			Ctx:      ctx,
			CtxClaim: vapusStudioClaim,
		},
		modelNode: result,
		request:   request,
	}
	agent.logger = pkgs.GetSubDMLogger("NewEmbeddingAgent", agent.AgentId)
	return agent, nil
}

func (s *EmbeddingAgent) GetEmbeddings() *mpb.Embeddings {
	return s.embeddings
}

func (s *EmbeddingAgent) Act() error {
	var err error
	modelConn, err := nabhiksvc.AIModelNodeConnectionPoolManager.GetorSetConnection(s.modelNode, true)
	if err != nil {
		s.logger.Error().Err(err).Msg("error while getting model node connection")
		return dmerrors.DMError(utils.ErrAIModelNode404, err)
	}
	if s.modelNode != nil {
		if s.modelNode.GetScope() == mpb.ResourceScope_ORG_SCOPE.String() {
			if slices.Contains(s.modelNode.ApprovedOrganizations, s.CtxClaim[encryption.ClaimOrganizationKey]) == false {
				s.logger.Error().Msg("error while processing action, model not available for organization")
				return dmerrors.DMError(utils.ErrAIModelNode403, nil)
			}
		}
	} else {
		s.logger.Error().Msg("error while processing action, model not found")
		return dmerrors.DMError(utils.ErrAIModelNode404, nil)
	}
	if s.request != nil {
		err = s.generateEmbeddings(modelConn)
	} else {
		s.logger.Error().Msg("error while processing action, invalid action")
		return dmerrors.DMError(utils.ErrAIModelManagerAction404, nil)
	}

	return nil
}

func (s *EmbeddingAgent) generateEmbeddings(modelConn aimodels.AIModelNodeInterface) error {
	payload := &prompts.AIEmbeddingPayload{
		Dimensions:     int(s.request.GetDimension()),
		EmbeddingModel: s.request.GetAiModel(),
		Input:          s.request.GetInputText(),
	}
	err := modelConn.GenerateEmbeddings(s.Ctx, payload, s.request.GetAiModel())
	if err != nil {
		s.logger.Error().Err(err).Msgf("error while generating embeddings from model %v", s.request.GetAiModel())
		return err
	}
	if payload.Embeddings.Vectors32 != nil {
		s.embeddings = &mpb.Embeddings{
			Embeddings32: payload.Embeddings.Vectors32,
			Type:         mpb.EmbeddingType_FLOAT_32,
			Dimension:    int64(payload.Dimensions),
		}
	} else {
		s.embeddings = &mpb.Embeddings{
			Embeddings64: payload.Embeddings.Vectors64,
			Type:         mpb.EmbeddingType_FLOAT_64,
			Dimension:    int64(payload.Dimensions),
		}
	}
	return nil
}
