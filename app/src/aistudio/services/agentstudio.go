package services

import (
	"context"
	"fmt"
	"log"

	grpccodes "google.golang.org/grpc/codes"

	"github.com/rs/zerolog"
	dmstores "github.com/vapusdata-oss/aistudio/aistudio/datastoreops"
	pkgs "github.com/vapusdata-oss/aistudio/aistudio/pkgs"
	"github.com/vapusdata-oss/aistudio/aistudio/utils"
	encryption "github.com/vapusdata-oss/aistudio/core/encryption"
	dmerrors "github.com/vapusdata-oss/aistudio/core/errors"
	models "github.com/vapusdata-oss/aistudio/core/models"
	grpcops "github.com/vapusdata-oss/aistudio/core/serviceops/grpcops"
	coreutils "github.com/vapusdata-oss/aistudio/core/utils"
	mpb "github.com/vapusdata-oss/apis/protos/models/v1alpha1"
	pb "github.com/vapusdata-oss/apis/protos/vapusai-studio/v1alpha1"
)

var VapusAIAgentRunnerMap = map[string]func(*VapusAIAgentThread, *VapusAgentChainLink) error{
	mpb.VapusAiAgentTypes_EMAILER.String():       (*VapusAIAgentThread).EmailerAgent,
	mpb.VapusAiAgentTypes_FILE_UPLOADER.String(): (*VapusAIAgentThread).FileUploader,
}

type VapusAIAgentThread struct {
	*models.VapusInterfaceAgentBase
	request         *pb.AgentInvokeRequest
	result          []*mpb.VapusContentObject
	streamServer    pb.AIAgentStudio_ChatStreamServer
	isStream        bool
	VapusAgentChain map[int]*VapusAgentChainLink
	*StudioServices
	logger       zerolog.Logger
	threadLog    *models.AIAgentThread
	currentIndex int
}

type VapusAgentChainLink struct {
	vapusAgent                   *models.VapusAIAgent
	linkId                       string
	userParams                   *pb.AgentInvokeLink
	requestAgentMap              map[string]string
	renderedInput                string
	resultData                   []map[string]interface{}
	dataFiles                    []*mpb.FileData
	dataFields                   []string
	dataproductId                string
	linkSteps                    map[string]*models.Steps
	failedStep                   string
	dbLog                        *models.AIAgentThreadLog
	statusCode                   grpccodes.Code
	eolReason                    mpb.EOSReasons
	NabhikServerResponseMetadata []*mpb.Mapper
	dataproductname              string
}

func (s *StudioServices) NewVapusAIAgentThread(ctx context.Context,
	request *pb.AgentInvokeRequest,
	streamServer pb.AIAgentStudio_ChatStreamServer) (*VapusAIAgentThread, error) {
	var threadId string = ""
	vapusStudioClaim, ok := encryption.GetCtxClaim(ctx)
	if !ok {
		s.logger.Error().Ctx(ctx).Msg("error while getting claim metadata from context")
		return nil, dmerrors.DMError(encryption.ErrInvalidJWTClaims, nil)
	}
	agentChain := make(map[int]*VapusAgentChainLink)
	for index, link := range request.Chain {
		s.logger.Info().Ctx(ctx).Msgf("Agent %s is ready for linking in chain", link.GetAgentId())
		vapusAgent, err := s.DMStore.GetAIAgent(ctx, link.GetAgentId(), vapusStudioClaim)
		if err != nil || vapusAgent == nil {
			s.logger.Error().Ctx(ctx).Err(err).Msg("error while getting agent")
			if streamServer == nil {
				return nil, err
			} else {
				streamServer.Send(&pb.AgentInvokeStreamResponse{
					StatusCode: int64(grpccodes.NotFound),
					Output: grpcops.BuildVapusStreamResponse(
						mpb.VapusStreamEvents_ABORTED,
						mpb.ContentFormats_PLAIN_TEXT,
						"",
						mpb.EOSReasons_INVALID_PARAMS,
						err,
					),
				})
			}
			return nil, err

		} else {
			agentChain[index] = &VapusAgentChainLink{
				vapusAgent:      vapusAgent,
				linkId:          coreutils.GetUUID(),
				userParams:      link,
				linkSteps:       make(map[string]*models.Steps),
				dataFiles:       make([]*mpb.FileData, 0),
				requestAgentMap: make(map[string]string),
				dbLog: &models.AIAgentThreadLog{
					AgentId:   vapusAgent.VapusID,
					AgentType: vapusAgent.AgentType,
				},
			}
			log.Println("======>>>>>>>>>>>>>>>>>>>>>>>", link.Steps)
			threadId = coreutils.GetUUID()
			if streamServer != nil {
				streamServer.Send(&pb.AgentInvokeStreamResponse{
					StatusCode: int64(grpccodes.OK),
					Output: grpcops.BuildVapusStreamResponse(
						mpb.VapusStreamEvents_START,
						mpb.ContentFormats_JSON,
						&models.Mapper{
							Key:   "Agent Invoked with thread id",
							Value: fmt.Sprintf("%s", threadId),
						},
						0,
						nil,
					),
				})
			}
		}
	}
	agent := &VapusAIAgentThread{
		request:         request,
		result:          make([]*mpb.VapusContentObject, 0),
		VapusAgentChain: agentChain,
		StudioServices:  s,
		VapusInterfaceAgentBase: &models.VapusInterfaceAgentBase{
			CtxClaim: vapusStudioClaim,
			Ctx:      ctx,
			InitAt:   coreutils.GetEpochTime(),
			AgentId:  threadId,
		},
	}

	agent.SetAgentId()
	if streamServer != nil {
		agent.isStream = true
		agent.streamServer = streamServer
	} else {
		agent.isStream = false
	}
	// for _, step := range agent.request.Steps {
	// 	agent.requestAgentMap[step.StepId] = step.Input
	// }
	agent.threadLog = &models.AIAgentThread{
		ThreadId: agent.AgentId,
		Log:      []*models.AIAgentThreadLog{},
	}
	agent.threadLog.PreSaveCreate(vapusStudioClaim)
	agent.threadLog.Organization = vapusStudioClaim[encryption.ClaimOrganizationKey]
	agent.logger = pkgs.GetSubDMLogger(agent.AgentType, agent.AgentId)
	for _, link := range agent.VapusAgentChain {
		if link.userParams.GetModelNodeId() == "" {
			if len(link.vapusAgent.AIModelMap) > 0 {
				link.userParams.ModelNodeId = link.vapusAgent.AIModelMap[0].ModelNodeId
			} else {
				link.userParams.ModelNodeId = dmstores.Account.AIAttributes.GenerativeModelNode
			}
		}
		if link.userParams.GetModelName() == "" {
			if len(link.vapusAgent.AIModelMap) > 0 {
				link.userParams.ModelName = link.vapusAgent.AIModelMap[0].ModelName
			} else {
				link.userParams.ModelName = dmstores.Account.AIAttributes.GenerativeModel
			}
		}
	}
	return agent, nil
}

func (s *VapusAIAgentThread) SetStepResponse(step, status string) map[string]string {
	return map[string]string{
		"key":   step,
		"value": status,
	}
}
func (s *VapusAIAgentThread) GetResult() []*mpb.VapusContentObject {
	return s.result
}

func (s *VapusAIAgentThread) Act() error {
	var err error
	if s == nil {
		return dmerrors.DMError(utils.ErrInvalidAgentChainParams, nil)
	}
	for _, link := range s.VapusAgentChain {
		err = s.RunAgent(link)
		if err != nil {
			s.logger.Error().Err(err).Msg("error while running agent")
			s.threadLog.Error = err.Error()
			break
		}
		link.dbLog.EOLReason = link.eolReason.String()
		link.dbLog.FailedStep = link.failedStep
		s.threadLog.Log = append(s.threadLog.Log, link.dbLog)
	}
	_ = s.DMStore.SaveAIAgentThread(s.Ctx, s.threadLog, s.CtxClaim)
	return err
}

func (s *VapusAIAgentThread) RunAgent(link *VapusAgentChainLink) error {
	var err error
	// switch link.vapusAgent.AgentType {
	// case mpb.VapusAiAgentTypes_EMAILER.String():
	// 	err = s.EmailerAgent(link)
	// }
	runner, ok := VapusAIAgentRunnerMap[link.vapusAgent.AgentType]
	if !ok {
		err = dmerrors.DMError(utils.ErrInvalidAgentType, nil)
		link.failedStep = "Invalid agent type"
		link.eolReason = mpb.EOSReasons_INVALID_PARAMS
		link.statusCode = grpccodes.InvalidArgument
		return err
	}
	err = runner(s, link)
	s.FinishAt = coreutils.GetEpochTime()
	s.FinalLog()
	if err != nil {
		s.logger.Error().Err(err).Msg("error while running agent")
		s.threadLog.Error = err.Error()
	}
	if s.isStream {
		s.streamServer.Send(&pb.AgentInvokeStreamResponse{
			StatusCode: int64(link.statusCode),
			Output: grpcops.BuildVapusStreamResponse(
				mpb.VapusStreamEvents_END,
				mpb.ContentFormats_JSON,
				s.GetAgentLogs(),
				link.eolReason,
				err,
			),
		})
	}
	return nil
}
