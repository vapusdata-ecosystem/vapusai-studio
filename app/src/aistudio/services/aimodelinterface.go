package services

import (
	"context"
	"fmt"
	"log"
	"slices"
	"strings"
	"sync"

	"github.com/rs/zerolog"
	"github.com/vapusdata-oss/aistudio/aistudio/nabhiksvc"
	"github.com/vapusdata-oss/aistudio/aistudio/pkgs"
	"github.com/vapusdata-oss/aistudio/aistudio/utils"
	aicore "github.com/vapusdata-oss/aistudio/core/aistudio/core"
	"github.com/vapusdata-oss/aistudio/core/aistudio/prompts"
	aimodels "github.com/vapusdata-oss/aistudio/core/aistudio/providers"
	encryption "github.com/vapusdata-oss/aistudio/core/encryption"
	dmerrors "github.com/vapusdata-oss/aistudio/core/errors"
	models "github.com/vapusdata-oss/aistudio/core/models"
	coreutils "github.com/vapusdata-oss/aistudio/core/utils"
	mpb "github.com/vapusdata-oss/apis/protos/models/v1alpha1"
	pb "github.com/vapusdata-oss/apis/protos/vapusai-studio/v1alpha1"
)

type AIModelAgent struct {
	*models.VapusInterfaceAgentBase
	payload        *prompts.GenerativePrompterPayload
	resultMetaData []*models.Mapper
	logger         zerolog.Logger
	request        *pb.ChatRequest
	modelNode      *models.AIModelNode
	actions        []pb.AIModelNodeAction
	streamServer   pb.AIModelStudio_ChatStreamServer
	dbLog          *models.AIModelStudioLog
	*StudioServices
	ChatStatus string
	isChat     bool
	totalInput string
}

func (s *StudioServices) NewAIModelAgent(ctx context.Context,
	request *pb.ChatRequest,
	stream pb.AIModelStudio_ChatStreamServer) (*AIModelAgent, error) {
	var err error
	modelsNodeId := ""
	vapusStudioClaim, ok := encryption.GetCtxClaim(ctx)
	if !ok {
		s.logger.Error().Ctx(ctx).Msg("error while getting claim metadata from context")
		return nil, dmerrors.DMError(encryption.ErrInvalidJWTClaims, nil)
	}

	if request != nil {
		modelsNodeId = request.GetModelNodeId()
	} else {
		s.logger.Error().Err(err).Msg("error while getting model node")
		return nil, dmerrors.DMError(utils.ErrAIModelNode404, nil)
	}
	result, err := s.DMStore.GetAIModelNode(ctx, modelsNodeId, vapusStudioClaim)
	if err != nil || result == nil {
		s.logger.Error().Err(err).Msg("error while getting model node")
		return nil, dmerrors.DMError(utils.ErrAIModelNode404, nil)
	}

	agent := &AIModelAgent{
		VapusInterfaceAgentBase: &models.VapusInterfaceAgentBase{
			AgentId:  coreutils.GetUUID(),
			Ctx:      ctx,
			CtxClaim: vapusStudioClaim,
		},
		modelNode:    result,
		request:      request,
		streamServer: stream,
		dbLog: &models.AIModelStudioLog{
			Mode:         request.GetMode().String(),
			Organization: vapusStudioClaim[encryption.ClaimOrganizationKey],
		},
		StudioServices: s,
	}
	if request == nil {
		s.logger.Error().Err(err).Msg("error : no spec found in request for AIModel interface")
		return nil, dmerrors.DMError(utils.ErrAIModelManagerAction404, nil)
	}

	if request.Mode == pb.AIInterfaceMode_CHAT_MODE {
		agent.ChatStatus = mpb.AIInterfaceChatStatus_CHAT_STARTED.String()
		agent.isChat = true
	}
	agent.dbLog.ThreadId = agent.AgentId
	agent.logger = pkgs.GetSubDMLogger("AIModelPromptAgent", agent.AgentId)
	return agent, nil
}

func (s *AIModelAgent) GetResult() *pb.ChatResponse {
	if s.payload == nil {
		return &pb.ChatResponse{
			Event: aicore.StreamEventEnd.String(),
		}
	}
	return s.payload.Response
}

func (s *AIModelAgent) Act() error {
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

	if s.request != nil && s.streamServer != nil {
		err = s.generateContentStream(modelConn)
	} else if s.request != nil && s.streamServer == nil {
		err = s.generateContent(modelConn)
	} else {
		s.logger.Error().Msg("error while processing action, invalid action")
		return dmerrors.DMError(utils.ErrAIModelManagerAction404, nil)
	}

	if err != nil {
		s.dbLog.Error = err.Error()
		s.dbLog.ResponseStatus = "400"
		s.logger.Error().Err(err).Msg("error while processing action")
		return dmerrors.DMError(utils.ErrAIModelManagerAction404, err)
	} else {
		s.dbLog.ResponseStatus = "200"
	}
	go func() {
		s.dbLog.UpdatedAt = coreutils.GetEpochTime()
		_ = s.DMStore.SaveAIInterfaceLog(context.TODO(), s.dbLog, s.CtxClaim)
	}()
	return nil
}

func (s *AIModelAgent) buildPayload(stream pb.AIModelStudio_ChatStreamServer) (*prompts.GenerativePrompterPayload, error) {
	var err error
	var payload *prompts.GenerativePrompterPayload

	if stream != nil {
		payload = prompts.NewPrompter(s.request, nil, stream, s.logger)
	} else {
		payload = prompts.NewPrompter(s.request, nil, nil, s.logger)
	}
	if s.request.MaxOutputTokens < 1 {
		s.request.MaxOutputTokens = prompts.DefaultMaxOPTokenLength
	}
	var wg sync.WaitGroup
	var errChan = make(chan error, 2)
	wg.Add(2)
	go func() {
		defer wg.Done()
		if s.request.PromptId != "" {
			prompt, err := s.DMStore.GetAIPrompt(s.Ctx, s.request.PromptId, s.CtxClaim)
			if err != nil {
				s.logger.Error().Err(err).Msgf("error while getting prompt %v", s.request.PromptId)
				errChan <- err
			}
			payload.Prompt = prompt
		}
		payload.RenderPrompt()
	}()
	go func() {
		defer wg.Done()
		if s.isChat {
			sessionMessages, err := s.DMStore.ListAIInterfaceLogByUser(s.Ctx, 5, s.CtxClaim)
			if err != nil {
				s.logger.Error().Err(err).Msg("error while getting session messages for current user")
				return
			}
			for _, message := range sessionMessages {
				if message.Input == "" || message.Output == "" {
					continue
				}
				if strings.TrimSpace(message.Input) != "" {
					if len(message.Input) > 500 {
						message.Input = coreutils.StringSlicer(message.Input, 500)
					}
					payload.SessionContext = append(payload.SessionContext, &prompts.SessionMessage{
						Message: message.Input,
						Role:    prompts.USER,
					})
				}
				if strings.TrimSpace(message.Output) != "" {
					if len(message.Output) > 500 {
						message.Output = coreutils.StringSlicer(message.Output, 500)
					}
					payload.SessionContext = append(payload.SessionContext, &prompts.SessionMessage{
						Message: message.Output,
						Role:    prompts.ASSISTANT,
					})
				}
			}
		}
	}()
	wg.Wait()
	close(errChan)
	for err := range errChan {
		if err != nil {
			return nil, err
		}
	}

	s.totalInput = ""
	log.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>1111", s.request.Messages)
	for _, mess := range s.request.Messages {
		if mess.Role != pb.AIMessageRoles_SYSTEM {
			s.totalInput += mess.Content
		}
	}
	s.payload = payload
	log.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>2222", s.totalInput)
	if len(s.request.Messages) < 1 || len(s.totalInput) < 1 {
		payload.GuardrailsFailed = true
		payload.ParsedOutput = "Please provide a valid input, input cannot be empty"
		return payload, nil
	}
	err = s.guardrailChecks(payload)
	if err != nil || payload.GuardrailsFailed {
		return payload, nil
	}
	s.dbLog.Input = s.totalInput

	return payload, nil
}

func (s *AIModelAgent) guardrailChecks(payload *prompts.GenerativePrompterPayload) error {
	log.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>3333", s.modelNode.VapusID)
	gdClients, ok := nabhiksvc.GuardrailPoolManager.ModelGuardRails[s.modelNode.VapusID]
	if !ok {
		return nil
	}
	log.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>4444", gdClients)
	for _, guardId := range gdClients {
		guard, ok := nabhiksvc.GuardrailPoolManager.GuardrailClientMap[guardId]
		if !ok {
			continue
		}
		scanResult := guard.Scan(s.Ctx, s.totalInput, s.logger)
		log.Println("scanResult", scanResult.ContentGuard, scanResult.TopicGuard, scanResult.WordGuard)
		if len(scanResult.WordGuard) > 0 || len(scanResult.TopicGuard) > 0 || len(scanResult.ContentGuard) > 0 {
			payload.ParsedOutput = guard.Guardrail.FailureMessage
			payload.GuardrailsFailed = true
			return fmt.Errorf("guardrail failed for your input %v", scanResult)
		}
	}
	return nil
}

func (s *AIModelAgent) generateContent(modelConn aimodels.AIModelNodeInterface) error {
	var err error
	payload, err := s.buildPayload(nil)
	if err != nil {
		s.logger.Error().Err(err).Msg("error while building payload")
		return err
	}
	if payload.GuardrailsFailed {
		payload.BuildResponseOP(aicore.StreamGuardrailFailed.String(), &prompts.PayloadgenericResponse{
			FinishReason: aicore.StreamGuardrailFailed.String(),
			Data:         payload.ParsedOutput,
			Role:         pb.AIMessageRoles_VAPUSGUARD.String(),
		}, false)
		return nil
	}
	log.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>", payload.Params.MaxOutputTokens, s.request.MaxOutputTokens)
	s.dbLog.Input = coreutils.StringSlicer(s.totalInput, 3700)
	err = modelConn.GenerateContent(s.Ctx, payload)
	if err != nil {
		s.logger.Error().Err(err).Msgf("error while generating content from model %v", s.request.Model)
		return err
	}
	log.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>> payload.Response", payload.Response)
	if len(payload.ParsedOutput) > 1 {
		s.dbLog.Output = coreutils.StringSlicer(payload.ParsedOutput, 3300)
	}
	return nil
}

func (s *AIModelAgent) generateContentStream(modelConn aimodels.AIModelNodeInterface) error {
	var err error
	payload, err := s.buildPayload(s.streamServer)
	if err != nil {
		s.logger.Error().Err(err).Msg("error while building payload")
		return err
	}
	log.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>", payload.Params.MaxOutputTokens, s.request.MaxOutputTokens)
	if payload.GuardrailsFailed {
		s.streamServer.Send(
			payload.BuildResponseOP(aicore.StreamGuardrailFailed.String(), &prompts.PayloadgenericResponse{
				Data: payload.ParsedOutput,
				Role: pb.AIMessageRoles_VAPUSGUARD.String(),
			}, true),
		)
		s.streamServer.SendMsg(&pb.StreamChatResponse{
			Output: &pb.ChatResponse{
				Event: aicore.StreamEventEnd.String(),
			},
		})
		return nil
	}
	err = modelConn.GenerateContentStream(s.Ctx, payload)
	if err != nil {
		s.logger.Error().Err(err).Msgf("error while generating content from model %v", s.request.Model)
		return err
	}
	if len(payload.ParsedOutput) > 1 {
		s.dbLog.Output = coreutils.StringSlicer(payload.ParsedOutput, 3300)
	}
	return nil
}
