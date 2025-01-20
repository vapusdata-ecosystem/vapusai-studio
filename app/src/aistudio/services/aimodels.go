package services

import (
	"context"
	"log"
	"slices"
	"strconv"

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

type VapusAINodeManagerAgent struct {
	*models.VapusInterfaceAgentBase
	managerRequest *pb.AIModelNodeConfiguratorRequest
	getterRequest  *pb.AIModelNodeGetterRequest
	result         *pb.AIModelNodeResponse
	dmStore        *dmstores.DMStore
}

func (s *StudioServices) NewVapusAINodeManagerAgent(ctx context.Context, managerRequest *pb.AIModelNodeConfiguratorRequest, getterRequest *pb.AIModelNodeGetterRequest) (*VapusAINodeManagerAgent, error) {
	vapusStudioClaim, ok := encryption.GetCtxClaim(ctx)
	if !ok {
		s.logger.Error().Ctx(ctx).Msg("error while getting claim metadata from context")
		return nil, dmerrors.DMError(encryption.ErrInvalidJWTClaims, nil)
	}
	agent := &VapusAINodeManagerAgent{
		managerRequest: managerRequest,
		getterRequest:  getterRequest,
		result:         &pb.AIModelNodeResponse{Output: &pb.AIModelNodeResponse_AIModelNodeResponse{}},
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
		// } else if getterRequest != nil {
		// 	if getterRequest.GetAiModelNodeId() == "" {
		// 		agent.Action = pb.AIModelNodeConfiguratorActions_LIST_AIMODEL_NODES.String()
		// 	} else {
		// 		agent.Action = pb.AIModelNodeConfiguratorActions_DESCRIBE_AIMODEL_NODE.String()
		// 	}
	} else {
		agent.Action = ""
	}
	agent.Logger = pkgs.GetSubDMLogger(globals.AISTUDIONODE.String(), agent.AgentId)
	return agent, nil
}

func (v *VapusAINodeManagerAgent) GetAgentId() string {
	return v.AgentId
}

func (v *VapusAINodeManagerAgent) GetResult() *pb.AIModelNodeResponse {
	return v.result
}

func (v *VapusAINodeManagerAgent) Act() error {
	switch v.GetAction() {
	case pb.AIModelNodeConfiguratorActions_CONFIGURE_AIMODEL_NODES.String():
		return v.configureAIModelNode()
	case pb.AIModelNodeConfiguratorActions_SYNC_AI_MODELS.String():
		return v.syncAIModelNode()
	case pb.AIModelNodeConfiguratorActions_PATCH_AIMODEL_NODES.String():
		return v.updateAIModelNode()
	case pb.AIModelNodeConfiguratorActions_DELETE_AIMODEL_NODES.String():
		return v.archiveAIModelNode()
	default:
		if v.getterRequest != nil {
			if v.getterRequest.GetAiModelNodeId() != "" {
				log.Println("describeAIModelNode", v.getterRequest.GetAiModelNodeId())
				return v.describeAIModelNode()
			} else {
				return v.listAIModelNodes()
			}
		}
		v.Logger.Error().Msg("invalid action")
		return utils.ErrInvalidAction
	}
}

func (v *VapusAINodeManagerAgent) updateCachePool(modelNode *models.AIModelNode) {
	if modelNode.SecurityGuardrails != nil {
		for _, guardrail := range modelNode.SecurityGuardrails.Guardrails {
			nabhiksvc.GuardrailPoolManager.AddGuardrails(modelNode.VapusID, guardrail)
		}
		_, _ = nabhiksvc.AIModelNodeConnectionPoolManager.GetorSetConnection(modelNode, true)
	}
	return
}

func (v *VapusAINodeManagerAgent) configureAIModelNode() error {
	aiModelNodes := utils.AIMNListPb2Obj(v.managerRequest.GetSpec())
	if len(aiModelNodes) == 0 {
		v.Logger.Error().Msg("error while unmarshalling ai model node from managerRequest")
		return dmerrors.DMError(utils.ErrInvalidAIModelNodeRequestSpec, nil)
	}
	// TODO: Track the failed nodes and return the failed nodes in the response of DMresponse.
	var err error
	var result []*models.AIModelNode
	for _, node := range aiModelNodes {
		node.SetAINodeId()
		if len(node.Editors) == 0 {
			node.Editors = []string{v.CtxClaim[encryption.ClaimUserIdKey]}
		} else {
			node.Editors = append(node.Editors, v.CtxClaim[encryption.ClaimUserIdKey])
		}
		node.PreSaveCreate(v.CtxClaim)
		node.Status = mpb.CommonStatus_ACTIVE.String()
		if node.Scope == mpb.ResourceScope_ORG_SCOPE.String() {
			if len(node.ApprovedOrganizations) == 0 {
				node.ApprovedOrganizations = []string{v.CtxClaim[encryption.ClaimOrganizationKey]}
			} else {
				node.ApprovedOrganizations = append(node.ApprovedOrganizations, v.CtxClaim[encryption.ClaimOrganizationKey])
			}
		}
		if node.ApprovedOrganizations == nil {
			node.ApprovedOrganizations = []string{v.CtxClaim[encryption.ClaimOrganizationKey]}
		} else {
			node.ApprovedOrganizations = append(node.ApprovedOrganizations, v.CtxClaim[encryption.ClaimOrganizationKey])
		}
		node.Organization = v.CtxClaim[encryption.ClaimOrganizationKey]
		newCtx := context.Background()
		if node.DiscoverModels {
			err = crawlAIModels(newCtx, node, v.Logger) // nolint:errcheck,gosec //
			if err != nil {
				v.Logger.Error().Err(err).Msg("error while crawling ai models")
				return dmerrors.DMError(utils.ErrCrawlingAIModels, err)
			}
		}
		err = v.setAiModelNodeCredentials(v.Ctx, coreutils.GetSecretName("aistudio", node.VapusID, "aiModelNode"), node, v.Logger)
		if err != nil {
			v.Logger.Error().Err(err).Msg("error while setting ai model node credentials")
			return dmerrors.DMError(utils.ErrSettingAIModelNodeCredentials, err)
		}
		err = v.dmStore.ConfigureGetAIModelNode(v.Ctx, node, v.CtxClaim)
		if err != nil {
			v.Logger.Error().Err(err).Msg("error while saving ai model node metadata in datastore")
			return dmerrors.DMError(utils.ErrUpdatingAIModelNode, err)
		}
		v.updateCachePool(node)
		result = append(result, node)
	}
	v.result.Output.AiModelNodes = utils.AIMNListM2Pb(result)
	return nil
}

func (v *VapusAINodeManagerAgent) syncAIModelNode() error {
	if len(v.managerRequest.GetSpec()) == 0 {
		v.Logger.Error().Msg("error while unmarshalling ai model nodes from managerRequest")
		return dmerrors.DMError(utils.ErrInvalidAIModelNodeRequestSpec, nil)
	}
	// TODO: Track the failed nodes and return the failed nodes in the response of DMresponse.
	var result []*models.AIModelNode
	for _, n := range v.managerRequest.GetSpec() {
		nodeObj, err := v.dmStore.GetAIModelNode(v.Ctx, n.ModelNodeId, v.CtxClaim)
		if err != nil {
			v.Logger.Error().Err(err).Msg("error while getting ai model node from datastore")
			return dmerrors.DMError(utils.ErrAIModelNode404, err)
		}

		if !slices.Contains(nodeObj.Editors, v.CtxClaim[encryption.ClaimUserIdKey]) {
			v.Logger.Error().Msg("error while validating user access")
			return dmerrors.DMError(utils.ErrAIModelNode403, nil)
		}

		nodeObj.PreSaveUpdate(v.CtxClaim[encryption.ClaimUserIdKey])
		netParams, err := v.dmStore.GetAIModelNodeNetworkParams(v.Ctx, nodeObj)
		if err != nil {
			v.Logger.Error().Err(err).Msg("error while getting ai model node network params")
			return dmerrors.DMError(utils.ErrGetAIModelNetParams, err)
		}
		nodeObj.NetworkParams = netParams

		err = crawlAIModels(v.Ctx, nodeObj, v.Logger) // nolint:errcheck,gosec //
		if err != nil {
			v.Logger.Error().Err(err).Msg("error while crawling ai models")
			return dmerrors.DMError(utils.ErrCrawlingAIModels, err)
		}
		nodeObj.NetworkParams.Credentials = nil
		err = v.dmStore.PutAIModelNode(v.Ctx, nodeObj)
		if err != nil {
			v.Logger.Error().Err(err).Msg("error while saving ai model node metadata in datastore")
			return dmerrors.DMError(utils.ErrSavingAIModelNode, err)
		}
		result = append(result, nodeObj)
	}
	v.result.Output.AiModelNodes = utils.AIMNListM2Pb(result)
	return nil
}

func (v *VapusAINodeManagerAgent) updateAIModelNode() error {
	if len(v.managerRequest.GetSpec()) == 0 {
		v.Logger.Error().Msg("error while unmarshalling ai model nodes from updateAIModelNode")
		return dmerrors.DMError(utils.ErrInvalidAIModelNodeRequestSpec, nil)
	}
	// TODO: Track the failed nodes and return the failed nodes in the response of DMresponse.
	var result []*models.AIModelNode
	for _, n := range v.managerRequest.GetSpec() {

		nodeObj, err := v.dmStore.GetAIModelNode(v.Ctx, n.ModelNodeId, v.CtxClaim)
		if err != nil {
			v.Logger.Error().Err(err).Msg("error while getting ai model node from datastore")
			return dmerrors.DMError(utils.ErrAIModelNode404, err)
		}
		if !slices.Contains(nodeObj.Editors, v.CtxClaim[encryption.ClaimUserIdKey]) {
			v.Logger.Error().Msg("error while validating user access")
			return dmerrors.DMError(utils.ErrAIModelNode403, nil)
		}
		if nodeObj.Organization != v.CtxClaim[encryption.ClaimOrganizationKey] {
			v.Logger.Error().Msg("error while validating organization access")
			return dmerrors.DMError(utils.ErrAIModelNode403, nil)
		}
		mConvertor := func(m []*mpb.AIModelBase) []*models.AIModelBase {
			r := make([]*models.AIModelBase, 0)
			if m != nil {
				for _, gm := range m {
					r = append(r, (&models.AIModelBase{}).ConvertFromPb(gm))
				}
			}
			return r
		}
		nodeObj.PreSaveUpdate(v.CtxClaim[encryption.ClaimUserIdKey])
		if nodeObj.Scope == mpb.ResourceScope_ORG_SCOPE.String() {

			nodeObj.ApprovedOrganizations = n.Attributes.GetApprovedOrgs()
		}
		nodeObj.SecurityGuardrails = (&models.SecurityGuardrails{}).ConvertFromPb(n.GetSecurityGuardrails())
		nodeObj.GenerativeModels = mConvertor(n.Attributes.GetGenerativeModels())
		nodeObj.EmbeddingModels = mConvertor(n.Attributes.GetEmbeddingModels())
		if n.Attributes.GetNetworkParams().Credentials != nil {
			err = v.setAiModelNodeCredentials(v.Ctx, coreutils.GetSecretName("aistudio", nodeObj.VapusID, strconv.Itoa(int(coreutils.GetEpochTime()))), nodeObj, v.Logger)
			if err != nil {
				v.Logger.Error().Err(err).Msg("error while updating ai model node credentials")
				return dmerrors.DMError(utils.ErrSettingAIModelNodeCredentials, err)
			}
		}
		err = v.dmStore.PutAIModelNode(v.Ctx, nodeObj)
		if err != nil {
			v.Logger.Error().Err(err).Msg("error while updating ai model node metadata in datastore")
			return dmerrors.DMError(utils.ErrUpdatingAIModelNode, err)
		}
		v.updateCachePool(nodeObj)
		result = append(result, nodeObj)
	}
	v.result.Output.AiModelNodes = utils.AIMNListM2Pb(result)
	return nil
}

func (v *VapusAINodeManagerAgent) describeAIModelNode() error {
	result, err := v.dmStore.GetAIModelNode(v.Ctx, v.getterRequest.GetAiModelNodeId(), v.CtxClaim)
	if err != nil {
		v.Logger.Error().Err(err).Msg("error while getting ai model node describe from datastore")
		return dmerrors.DMError(err, nil)
	}

	v.result.Output.AiModelNodes = utils.AIMNListM2Pb([]*models.AIModelNode{result})
	return nil
}

func (v *VapusAINodeManagerAgent) archiveAIModelNode() error {
	if len(v.managerRequest.GetSpec()) == 0 {
		v.Logger.Error().Msg("error while unmarshalling ai model nodes from updateAIModelNode")
		return dmerrors.DMError(utils.ErrInvalidAIModelNodeRequestSpec, nil)
	}
	// TODO: Track the failed nodes and return the failed nodes in the response of DMresponse.
	for _, n := range v.managerRequest.GetSpec() {
		result, err := v.dmStore.GetAIModelNode(v.Ctx, n.GetModelNodeId(), v.CtxClaim)
		if err != nil {
			v.Logger.Error().Err(err).Msg("error while getting ai model node describe from datastore")
			return dmerrors.DMError(err, nil)
		}
		result.Status = mpb.CommonStatus_DELETED.String()
		result.PreSaveDelete(v.CtxClaim)
		err = v.dmStore.PutAIModelNode(v.Ctx, result)
		if err != nil {
			v.Logger.Error().Err(err).Msg("error while archiving ai model node")
			return dmerrors.DMError(err, nil)
		}
	}
	v.result.Output.AiModelNodes = []*mpb.AIModelNode{}
	return nil
}

func (v *VapusAINodeManagerAgent) listAIModelNodes() error {
	// condition := fmt.Sprintf("(status = 'ACTIVE' AND deleted_at IS NULL AND scope = 'PLATFORM_SCOPE' OR (scope = 'ORG_SCOPE' AND '%s' = ANY(approved_organizations))) OR organization='%s' ORDER BY created_at DESC", v.CtxClaim[encryption.ClaimOrganizationKey], v.CtxClaim[encryption.ClaimOrganizationKey])
	result, err := v.dmStore.ListAIModelNodes(v.Ctx, vapusorms.ListResourceWithGovernance(v.CtxClaim), v.CtxClaim)
	if err != nil {
		v.Logger.Error().Err(err).Msg("error while getting ai model node list from datastore")
		return dmerrors.DMError(err, nil)
	}
	v.result.Output.AiModelNodes = utils.AIMNListM2Pb(result)
	return nil
}
