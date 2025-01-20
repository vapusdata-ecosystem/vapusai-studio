package services

import (
	"context"

	dmstores "github.com/vapusdata-oss/aistudio/aistudio/datastoreops"
	pkgs "github.com/vapusdata-oss/aistudio/aistudio/pkgs"
	"github.com/vapusdata-oss/aistudio/aistudio/utils"
	encryption "github.com/vapusdata-oss/aistudio/core/encryption"
	dmerrors "github.com/vapusdata-oss/aistudio/core/errors"
	"github.com/vapusdata-oss/aistudio/core/globals"
	models "github.com/vapusdata-oss/aistudio/core/models"
	vapusorms "github.com/vapusdata-oss/aistudio/core/serviceops"
	coreutils "github.com/vapusdata-oss/aistudio/core/utils"
	mpb "github.com/vapusdata-oss/apis/protos/models/v1alpha1"
	pb "github.com/vapusdata-oss/apis/protos/vapusai-studio/v1alpha1"
)

type VapusAIPromptManagerAgent struct {
	*models.VapusInterfaceAgentBase
	managerRequest *pb.PromptManagerRequest
	getterRequest  *pb.PromptGetterRequest
	result         []*mpb.AIModelPrompt
	dmStore        *dmstores.DMStore
}

func (s *StudioServices) NewVapusAIPromptManagerAgent(ctx context.Context, managerRequest *pb.PromptManagerRequest, getterRequest *pb.PromptGetterRequest) (*VapusAIPromptManagerAgent, error) {
	vapusStudioClaim, ok := encryption.GetCtxClaim(ctx)
	if !ok {
		s.logger.Error().Ctx(ctx).Msg("error while getting claim metadata from context")
		return nil, dmerrors.DMError(encryption.ErrInvalidJWTClaims, nil)
	}
	agent := &VapusAIPromptManagerAgent{
		managerRequest: managerRequest,
		getterRequest:  getterRequest,
		result:         make([]*mpb.AIModelPrompt, 0),
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

func (v *VapusAIPromptManagerAgent) GetAgentId() string {
	return v.AgentId
}

func (v *VapusAIPromptManagerAgent) GetResult() []*mpb.AIModelPrompt {
	v.FinishAt = coreutils.GetEpochTime()
	v.FinalLog()
	return v.result
}

func (v *VapusAIPromptManagerAgent) Act() error {
	switch v.GetAction() {
	case pb.PromptAgentAction_CONFIGURE_PROMPT.String():
		return v.configureAIPrompts()
	case pb.PromptAgentAction_PATCH_PROMPT.String():
		return v.patchAIPrompts()
	default:
		if v.getterRequest != nil {
			if v.getterRequest.GetPromptId() != "" {
				return v.describeAIPrompts()
			} else {
				return v.listAIModelPrompts()
			}
		}
		v.Logger.Error().Msg("invalid action")
		return utils.ErrInvalidAction
	}
}

func (v *VapusAIPromptManagerAgent) configureAIPrompts() error {
	aiPrompts := utils.AIPRPb2Obj(v.managerRequest.GetSpec())
	if len(aiPrompts) == 0 {
		v.Logger.Error().Msg("error while unmarshalling ai model prompt from managerRequest")
		return dmerrors.DMError(utils.ErrInvalidAIModelPromptRequestSpec, nil)
	}
	var err error
	var result []*models.AIModelPrompt
	for _, prompt := range aiPrompts {
		prompt.SetPromptId()
		prompt.PreSaveCreate(v.CtxClaim)
		prompt.Status = mpb.CommonStatus_ACTIVE.String()
		prompt.Organization = v.CtxClaim[encryption.ClaimOrganizationKey]
		if len(prompt.Prompt.Tools) > 0 {
			BuildPromptSchema(prompt)
		}
		BuildAIPromptTemplate(prompt)
		err = v.dmStore.ConfigureAIPrompts(v.Ctx, prompt, v.CtxClaim)
		if err != nil {
			v.Logger.Error().Err(err).Msg("error while saving ai prompt to datastore")
			return dmerrors.DMError(utils.ErrSavingAIModelPrompt, err)
		}
		result = append(result, prompt)
	}
	v.result = utils.AIPRObj2Pb(result)
	return nil
}

func (v *VapusAIPromptManagerAgent) patchAIPrompts() error {
	aiPrompts := utils.AIPRPb2Obj(v.managerRequest.GetSpec())
	if len(aiPrompts) == 0 {
		v.Logger.Error().Msg("error while unmarshalling ai model prompt from managerRequest")
		return dmerrors.DMError(utils.ErrInvalidAIModelPromptRequestSpec, nil)
	}
	var result []*models.AIModelPrompt
	for _, pbObj := range v.managerRequest.GetSpec() {
		exprompt, err := v.dmStore.GetAIPrompt(v.Ctx, pbObj.PromptId, v.CtxClaim)
		if err != nil {
			v.Logger.Error().Err(err).Msgf("error while getting ai model prompt with id - %s from datastore for patching", pbObj.PromptId)
			return dmerrors.DMError(utils.ErrAIPrompt404, err)
		}
		if !(exprompt.CreatedBy == v.CtxClaim[encryption.ClaimUserIdKey] && exprompt.Organization == v.CtxClaim[encryption.ClaimOrganizationKey]) {
			v.Logger.Error().Msg("error while validating user access")
			return dmerrors.DMError(utils.ErrPrompt403, nil)
		}
		if !exprompt.Editable {
			v.Logger.Error().Msgf("ai model prompt with id - %s is not editable", pbObj.PromptId)
			return dmerrors.DMError(utils.ErrAIPromptNotEditable, nil)
		}
		newPrompt := (&models.AIModelPrompt{}).ConvertFromPb(pbObj)
		newPrompt.VapusBase = exprompt.VapusBase
		newPrompt.PreSaveUpdate(v.CtxClaim[encryption.ClaimUserIdKey])
		newPrompt.Organization = exprompt.Organization
		newPrompt.Status = mpb.CommonStatus_ACTIVE.String()
		if len(newPrompt.Prompt.Tools) > 0 {
			BuildPromptSchema(newPrompt)
		}
		BuildAIPromptTemplate(newPrompt)
		err = v.dmStore.PutAIPrompts(v.Ctx, newPrompt, v.CtxClaim)
		if err != nil {
			v.Logger.Error().Err(err).Msg("error while updating ai prompt to datastore")
			return dmerrors.DMError(utils.ErrSavingAIModelPrompt, err)
		}
		result = append(result, newPrompt)
	}
	v.result = utils.AIPRObj2Pb(result)
	return nil
}

func (v *VapusAIPromptManagerAgent) describeAIPrompts() error {
	v.Logger.Info().Msg("getting ai model prompt describe from datastore")
	result, err := v.dmStore.GetAIPrompt(v.Ctx, v.getterRequest.GetPromptId(), v.CtxClaim)
	if err != nil {
		v.Logger.Error().Err(err).Msg("error while getting ai model prompt describe from datastore")
		return dmerrors.DMError(utils.ErrAIPrompt404, err)
	}
	v.result = []*mpb.AIModelPrompt{result.ConvertToPb()}
	return nil
}

func (v *VapusAIPromptManagerAgent) listAIModelPrompts() error {
	v.Logger.Info().Msg("getting ai model prompt list from datastore")
	result, err := v.dmStore.ListAIPrompts(v.Ctx,
		vapusorms.ListResourceWithGovernance(v.CtxClaim), v.CtxClaim)
	if err != nil {
		v.Logger.Error().Err(err).Msg("error while getting ai model prompt list from datastore")
		return dmerrors.DMError(utils.ErrAIPrompt404, err)
	}
	v.result = utils.AIPRObj2Pb(result)
	return nil
}
