package services

import (
	"context"
	"fmt"
	"slices"

	dmstores "github.com/vapusdata-oss/aistudio/aistudio/datastoreops"
	pkgs "github.com/vapusdata-oss/aistudio/aistudio/pkgs"
	"github.com/vapusdata-oss/aistudio/aistudio/utils"
	encryption "github.com/vapusdata-oss/aistudio/core/encryption"
	dmerrors "github.com/vapusdata-oss/aistudio/core/errors"
	"github.com/vapusdata-oss/aistudio/core/globals"
	models "github.com/vapusdata-oss/aistudio/core/models"
	svcops "github.com/vapusdata-oss/aistudio/core/serviceops"
	coreutils "github.com/vapusdata-oss/aistudio/core/utils"
	nabhikutils "github.com/vapusdata-oss/aistudio/core/utils"
	mpb "github.com/vapusdata-oss/apis/protos/models/v1alpha1"
	pb "github.com/vapusdata-oss/apis/protos/vapusai-studio/v1alpha1"
)

type VapusAIAgentManagerAgent struct {
	*models.VapusInterfaceAgentBase
	managerRequest *pb.AgentManagerRequest
	getterRequest  *pb.AgentGetterRequest
	result         []*mpb.VapusAIAgent
	dmStore        *dmstores.DMStore
}

func (s *StudioServices) NewVapusAIAgentManagerAgent(ctx context.Context, managerRequest *pb.AgentManagerRequest, getterRequest *pb.AgentGetterRequest) (*VapusAIAgentManagerAgent, error) {
	vapusStudioClaim, ok := encryption.GetCtxClaim(ctx)
	if !ok {
		s.logger.Error().Ctx(ctx).Msg("error while getting claim metadata from context")
		return nil, dmerrors.DMError(encryption.ErrInvalidJWTClaims, nil)
	}
	agent := &VapusAIAgentManagerAgent{
		managerRequest: managerRequest,
		getterRequest:  getterRequest,
		result:         make([]*mpb.VapusAIAgent, 0),
		dmStore:        s.DMStore,
		VapusInterfaceAgentBase: &models.VapusInterfaceAgentBase{
			CtxClaim: vapusStudioClaim,
			Ctx:      ctx,
			Action:   managerRequest.GetAction().String(),
			InitAt:   coreutils.GetEpochTime(),
		},
	}
	agent.SetAgentId()
	if managerRequest != nil {
		agent.Action = managerRequest.GetAction().String()
	} else {
		agent.Action = ""
	}
	agent.Logger = pkgs.GetSubDMLogger(globals.AIPROMPTAGENT.String(), agent.AgentId)
	return agent, nil
}

func (v *VapusAIAgentManagerAgent) GetAgentId() string {
	return v.AgentId
}

func (v *VapusAIAgentManagerAgent) GetResult() []*mpb.VapusAIAgent {
	v.FinishAt = coreutils.GetEpochTime()
	v.FinalLog()
	return v.result
}

func (v *VapusAIAgentManagerAgent) Act() error {
	switch v.GetAction() {
	case pb.VapusAIAgentAction_CONFIGURE_AIAGENT.String():
		return v.configureAIAgents()
	case pb.VapusAIAgentAction_PATCH_AIAGENT.String():
		return v.patchAIAgents()
	default:
		if v.getterRequest != nil {
			if v.getterRequest.GetAgentId() != "" {
				return v.describeAIAgents()
			} else {
				return v.listAIAgents()
			}
		}
		v.Logger.Error().Msg("invalid action")
		return dmerrors.DMError(utils.ErrInvalidAction, nil)
	}
}

func (v *VapusAIAgentManagerAgent) validateAIModel(modelMaps []*models.AIModelMap) error {
	validModels, err := v.dmStore.ListAIModelNodes(v.Ctx,
		fmt.Sprintf("(status = 'ACTIVE' AND deleted_at IS NULL)  ORDER BY created_at DESC"),
		v.CtxClaim)
	if err != nil {
		v.Logger.Error().Err(err).Msg("error while getting ai model list from datastore")
		return dmerrors.DMError(utils.ErrInvalidAIModels, err)
	}
	if len(validModels) == 0 {
		v.Logger.Error().Msg("no valid ai models found")
		return dmerrors.DMError(utils.ErrInvalidAIModels, nil)
	}
	validModelMap := make(map[string][]string)
	for _, model := range validModels {
		validModelMap[model.VapusID] = func() []string {
			var v []string
			for _, m := range model.GetGenerativeModels() {
				v = append(v, m.ModelName)
			}
			return v
		}()
		validModelMap[model.VapusID] = append(validModelMap[model.VapusID], func() []string {
			var v []string
			for _, m := range model.GetEmbeddingModels() {
				v = append(v, m.ModelName)
			}
			return v
		}()...)
	}
	valid := false
	for _, modelMap := range modelMaps {
		if modelMap.ModelNodeId == "" || modelMap.ModelName == "" {
			v.Logger.Error().Msg("empty ai model node requested")
			valid = true
			continue
		}
		if models, ok := validModelMap[modelMap.ModelNodeId]; !ok {
			v.Logger.Error().Msg("invalid ai model node requested")
			return dmerrors.DMError(utils.ErrInvalidAIModelNodeRequested, nil)
		} else {
			if !slices.Contains(models, modelMap.ModelName) {
				v.Logger.Error().Msg("invalid ai model requested")
				return dmerrors.DMError(utils.ErrInvalidAIModelNodeRequested, nil)
			} else {
				valid = true
			}
		}
	}
	if !valid {
		v.Logger.Error().Msg("invalid ai model requested")
		return dmerrors.DMError(utils.ErrInvalidAIModelNodeRequested, nil)
	}
	return nil
}

func (v *VapusAIAgentManagerAgent) configureAIAgents() error {
	aiAgent := (&models.VapusAIAgent{}).ConvertFromPb(v.managerRequest.GetSpec())
	aiAgent.PreSaveCreate(v.CtxClaim)
	aiAgent.Status = mpb.CommonStatus_ACTIVE.String()
	aiAgent.Organization = v.CtxClaim[encryption.ClaimOrganizationKey]
	aiAgent.AgentVersion = nabhikutils.GetVersionNumber("", "")
	if aiAgent.Editors == nil {
		aiAgent.Editors = []string{v.CtxClaim[encryption.ClaimUserIdKey]}
	} else if !slices.Contains(aiAgent.Editors, v.CtxClaim[encryption.ClaimUserIdKey]) {
		aiAgent.Editors = append(aiAgent.Editors, v.CtxClaim[encryption.ClaimUserIdKey])
	}
	err := v.validateAIModel(aiAgent.AIModelMap)
	if err != nil {
		v.Logger.Error().Err(err).Msg("error while validating ai model map")
		return dmerrors.DMError(utils.ErrInvalidAIModelNodeRequested, err)
	}
	aiAgent.Settings = &models.AIAgentSettings{
		ToolCallSchema: BuildAgentFuntioncallRenderer(aiAgent),
	}
	err = v.dmStore.ConfigureAIAgents(v.Ctx, aiAgent, v.CtxClaim)
	if err != nil {
		v.Logger.Error().Err(err).Msg("error while configuring ai agent")
		return dmerrors.DMError(utils.ErrAIAgentCreate500, err)
	}
	v.result = []*mpb.VapusAIAgent{aiAgent.ConvertToPb()}
	return nil
}

func (v *VapusAIAgentManagerAgent) patchAIAgents() error {
	existingAIAgent, err := v.dmStore.GetAIAgent(v.Ctx, v.managerRequest.GetSpec().GetAgentId(), v.CtxClaim)
	if err != nil {
		v.Logger.Error().Err(err).Msg("error while getting ai agent from datastore")
		return dmerrors.DMError(utils.ErrAIAgent404, err)
	}
	if !existingAIAgent.Editable {
		v.Logger.Error().Msg("ai agent is not editable")
		return dmerrors.DMError(utils.ErrAIAgentNotEditable, nil)
	}
	if !(existingAIAgent.CreatedBy == v.CtxClaim[encryption.ClaimUserIdKey] && existingAIAgent.Organization == v.CtxClaim[encryption.ClaimOrganizationKey]) {
		v.Logger.Error().Msg("error while validating user access to ai guardrail")
		return dmerrors.DMError(utils.ErrPrompt403, nil)
	}
	aiAgent := (&models.VapusAIAgent{}).ConvertFromPb(v.managerRequest.GetSpec())
	existingAIAgent.PreSaveUpdate(v.CtxClaim[encryption.ClaimUserIdKey])
	err = v.validateAIModel(aiAgent.AIModelMap)
	if err != nil {
		v.Logger.Error().Err(err).Msg("error while validating ai model map")
		return dmerrors.DMError(utils.ErrInvalidAIModelNodeRequested, err)
	}

	existingAIAgent.Steps = aiAgent.Steps
	existingAIAgent.Description = aiAgent.Description
	if aiAgent.Editors != nil {
		existingAIAgent.Editors = aiAgent.Editors
	}
	existingAIAgent.Editors = aiAgent.Editors
	existingAIAgent.Editable = aiAgent.Editable
	existingAIAgent.AIModelMap = aiAgent.AIModelMap
	existingAIAgent.Status = mpb.CommonStatus_ACTIVE.String()
	existingAIAgent.Organization = v.CtxClaim[encryption.ClaimOrganizationKey]
	existingAIAgent.Labels = aiAgent.Labels
	existingAIAgent.AgentVersion = nabhikutils.GetVersionNumber(existingAIAgent.AgentVersion, mpb.VersionBumpType_MINOR.String())
	if existingAIAgent.Settings == nil {
		existingAIAgent.Settings = &models.AIAgentSettings{}
	}
	existingAIAgent.Settings.ToolCallSchema = BuildAgentFuntioncallRenderer(existingAIAgent)
	err = v.dmStore.PutAIAgents(v.Ctx, existingAIAgent, v.CtxClaim)
	if err != nil {
		v.Logger.Error().Err(err).Msg("error while configuring ai agent")
		return dmerrors.DMError(utils.ErrAIAgentPatch400, err)
	}
	v.result = []*mpb.VapusAIAgent{existingAIAgent.ConvertToPb()}
	return nil
}

func (v *VapusAIAgentManagerAgent) describeAIAgents() error {
	v.Logger.Info().Msg("getting ai agent describe from datastore")
	result, err := v.dmStore.GetAIAgent(v.Ctx, v.getterRequest.GetAgentId(), v.CtxClaim)
	if err != nil {
		v.Logger.Error().Err(err).Msg("error while getting ai agent describe from datastore")
		return dmerrors.DMError(utils.ErrAIAgent404, err)
	}
	v.result = []*mpb.VapusAIAgent{result.ConvertToPb()}
	return nil
}

func (v *VapusAIAgentManagerAgent) listAIAgents() error {
	v.Logger.Info().Msg("getting ai agent list from datastore")
	result, err := v.dmStore.ListAIAgents(v.Ctx,
		svcops.ListResourceWithGovernance(v.CtxClaim), v.CtxClaim)
	if err != nil {
		v.Logger.Error().Err(err).Msg("error while getting ai agents list from datastore")
		return dmerrors.DMError(utils.ErrAIAgent404, err)
	}
	v.result = utils.AIAGPb2Obj(result)
	return nil
}
