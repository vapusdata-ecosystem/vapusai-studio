package services

import (
	"context"
	"slices"

	dmstores "github.com/vapusdata-oss/aistudio/aistudio/datastoreops"
	"github.com/vapusdata-oss/aistudio/aistudio/nabhiksvc"
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

type VapusAIGuardrailManagerAgent struct {
	*models.VapusInterfaceAgentBase
	managerRequest *pb.GuardrailsManagerRequest
	getterRequest  *pb.GuardrailsGetterRequest
	result         []*mpb.AIGuardrails
	dmStore        *dmstores.DMStore
}

func (s *StudioServices) NewVapusAIGuardrailManagerAgent(ctx context.Context, managerRequest *pb.GuardrailsManagerRequest, getterRequest *pb.GuardrailsGetterRequest) (*VapusAIGuardrailManagerAgent, error) {
	vapusStudioClaim, ok := encryption.GetCtxClaim(ctx)
	if !ok {
		s.logger.Error().Ctx(ctx).Msg("error while getting claim metadata from context")
		return nil, dmerrors.DMError(encryption.ErrInvalidJWTClaims, nil)
	}
	guardrail := &VapusAIGuardrailManagerAgent{
		managerRequest: managerRequest,
		getterRequest:  getterRequest,
		result:         make([]*mpb.AIGuardrails, 0),
		dmStore:        s.DMStore,
		VapusInterfaceAgentBase: &models.VapusInterfaceAgentBase{
			CtxClaim: vapusStudioClaim,
			Ctx:      ctx,
			Action:   managerRequest.GetAction().String(),
			InitAt:   coreutils.GetEpochTime(),
		},
	}
	guardrail.SetAgentId()
	if managerRequest != nil {
		guardrail.Action = managerRequest.GetAction().String()
	} else {
		guardrail.Action = ""
	}
	guardrail.Logger = pkgs.GetSubDMLogger(globals.AIPROMPTAGENT.String(), guardrail.AgentId)
	return guardrail, nil
}

func (v *VapusAIGuardrailManagerAgent) GetAgentId() string {
	return v.AgentId
}

func (v *VapusAIGuardrailManagerAgent) GetResult() []*mpb.AIGuardrails {
	v.FinishAt = coreutils.GetEpochTime()
	v.FinalLog()
	return v.result
}

func (v *VapusAIGuardrailManagerAgent) Act() error {
	switch v.GetAction() {
	case pb.VapusAIGuardrailsAction_CONFIGURE_GUARDRAIL.String():
		return v.configureAIGuardrails()
	case pb.VapusAIGuardrailsAction_PATCH_GUARDRAIL.String():
		return v.patchAIGuardrails()
	default:
		if v.getterRequest != nil {
			if v.getterRequest.GetGuardrailId() != "" {
				return v.describeAIGuardrails()
			} else {
				return v.listAIGuardrails()
			}
		}
		v.Logger.Error().Msg("invalid action")
		return utils.ErrInvalidAction
	}
}

func (v *VapusAIGuardrailManagerAgent) updateCachePool(guardrail *models.AIGuardrails) {
	nabhiksvc.GuardrailPoolManager.UpdateGuardrailPool(guardrail)
	return
}

func (v *VapusAIGuardrailManagerAgent) configureAIGuardrails() error {
	aiGuardrails := (&models.AIGuardrails{}).ConvertFromPb(v.managerRequest.GetSpec())
	aiGuardrails.PreSaveCreate(v.CtxClaim)
	aiGuardrails.Status = mpb.CommonStatus_ACTIVE.String()
	aiGuardrails.Organization = v.CtxClaim[encryption.ClaimOrganizationKey]
	if aiGuardrails.Editors == nil {
		aiGuardrails.Editors = []string{v.CtxClaim[encryption.ClaimUserIdKey]}
	} else if !slices.Contains(aiGuardrails.Editors, v.CtxClaim[encryption.ClaimUserIdKey]) {
		aiGuardrails.Editors = append(aiGuardrails.Editors, v.CtxClaim[encryption.ClaimUserIdKey])
	}
	aiGuardrails.Schema = BuildGuardrailSchema(aiGuardrails)
	err := v.dmStore.ConfigureAIGuardrails(v.Ctx, aiGuardrails, v.CtxClaim)
	if err != nil {
		v.Logger.Error().Err(err).Msg("error while configuring ai guardrail")
		return dmerrors.DMError(utils.ErrAIGuardrailCreate400, err)
	}
	v.updateCachePool(aiGuardrails)
	v.result = []*mpb.AIGuardrails{aiGuardrails.ConvertToPb()}
	return nil
}

func (v *VapusAIGuardrailManagerAgent) patchAIGuardrails() error {
	existingAIGuardrail, err := v.dmStore.GetAIGuardrail(v.Ctx, v.managerRequest.GetSpec().GetGuardrailId(), v.CtxClaim)
	if err != nil {
		v.Logger.Error().Err(err).Msg("error while getting ai guardrail from datastore")
		return dmerrors.DMError(utils.ErrAIGuardrail404, err)
	}
	if !(existingAIGuardrail.CreatedBy == v.CtxClaim[encryption.ClaimUserIdKey] && existingAIGuardrail.Organization == v.CtxClaim[encryption.ClaimOrganizationKey]) {
		v.Logger.Error().Msg("error while validating user access to ai guardrail")
		return dmerrors.DMError(utils.ErrAIGuardrail403, nil)
	}
	aiGuardrails := (&models.AIGuardrails{}).ConvertFromPb(v.managerRequest.GetSpec())
	existingAIGuardrail.PreSaveUpdate(v.CtxClaim[encryption.ClaimUserIdKey])

	existingAIGuardrail.Description = aiGuardrails.Description
	if aiGuardrails.Editors != nil {
		existingAIGuardrail.Editors = aiGuardrails.Editors
	}
	existingAIGuardrail.Contents = aiGuardrails.Contents
	existingAIGuardrail.Topics = aiGuardrails.Topics
	existingAIGuardrail.Words = aiGuardrails.Words
	existingAIGuardrail.SensitiveDataset = aiGuardrails.SensitiveDataset
	existingAIGuardrail.Editors = aiGuardrails.Editors
	existingAIGuardrail.Status = mpb.CommonStatus_ACTIVE.String()
	existingAIGuardrail.Organization = v.CtxClaim[encryption.ClaimOrganizationKey]
	existingAIGuardrail.Description = aiGuardrails.Description
	existingAIGuardrail.FailureMessage = aiGuardrails.FailureMessage
	existingAIGuardrail.Schema = BuildGuardrailSchema(existingAIGuardrail)
	existingAIGuardrail.ScanMode = aiGuardrails.ScanMode
	existingAIGuardrail.GuardModel = aiGuardrails.GuardModel
	err = v.dmStore.PutAIGuardrails(v.Ctx, existingAIGuardrail, v.CtxClaim)
	if err != nil {
		v.Logger.Error().Err(err).Msg("error while configuring ai guardrail")
		return dmerrors.DMError(utils.ErrAIGuardrailPatch400, err)
	}
	v.updateCachePool(existingAIGuardrail)
	v.result = []*mpb.AIGuardrails{existingAIGuardrail.ConvertToPb()}
	return nil
}

func (v *VapusAIGuardrailManagerAgent) describeAIGuardrails() error {
	v.Logger.Info().Msg("getting ai guardrail describe from datastore")
	result, err := v.dmStore.GetAIGuardrail(v.Ctx, v.getterRequest.GetGuardrailId(), v.CtxClaim)
	if err != nil {
		v.Logger.Error().Err(err).Msg("error while getting ai guardrail describe from datastore")
		return dmerrors.DMError(utils.ErrAIGuardrail404, err)
	}
	v.result = []*mpb.AIGuardrails{result.ConvertToPb()}
	return nil
}

func (v *VapusAIGuardrailManagerAgent) listAIGuardrails() error {
	v.Logger.Info().Msg("getting ai guardrail list from datastore")
	result, err := v.dmStore.ListAIGuardrails(v.Ctx,
		vapusorms.ListResourceWithGovernance(v.CtxClaim), v.CtxClaim)
	if err != nil {
		v.Logger.Error().Err(err).Msg("error while getting ai guardrails list from datastore")
		return dmerrors.DMError(utils.ErrAIGuardrail404, err)
	}
	v.result = utils.AIGDObjToPb(result)
	return nil
}
