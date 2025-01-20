package services

import (
	"fmt"
	"net/http"
	"slices"

	"github.com/labstack/echo/v4"
	"github.com/vapusdata-oss/aistudio/core/globals"
	serviceops "github.com/vapusdata-oss/aistudio/core/serviceops"
	"github.com/vapusdata-oss/aistudio/core/serviceops/grpcops"
	"github.com/vapusdata-oss/aistudio/webapp/models"
	pkgs "github.com/vapusdata-oss/aistudio/webapp/pkgs"
	"github.com/vapusdata-oss/aistudio/webapp/routes"
	"github.com/vapusdata-oss/aistudio/webapp/utils"
	mpb "github.com/vapusdata-oss/apis/protos/models/v1alpha1"
	pb "github.com/vapusdata-oss/apis/protos/vapusai-studio/v1alpha1"
)

func (x *WebappService) ManageAIModelNodesHandler(c echo.Context) error {
	response := &models.AIStudioResponse{
		AIModelNodes:    x.grpcClients.AIModelNodes(c),
		ActionParams:    &models.ResourceManagerParams{},
		BackListingLink: c.Request().URL.String(),
	}
	globalContext, err := x.getAiStudioSectionGlobals(c, routes.ManageAIModelNodesPage.String())
	if err != nil {
		x.logger.Err(err).Msg("error while getting aistudio section globals")
		return HandleGLobalContextError(c, err)
	}

	response.ActionParams.API = fmt.Sprintf("%s/api/v1alpha1/aistudio/models/nodes", pkgs.NetworkConfigManager.ExternalURL)
	return c.Render(http.StatusOK, "manageai-modelnodes.html", map[string]interface{}{
		"GlobalContext":  globalContext,
		"Response":       response,
		"SectionHeader":  "AI Models",
		"CreateTemplate": GetProtoYamlString(serviceops.AinodeConfiguratorRequest),
	})
}

func (x *WebappService) ManageAIModelNodesDetailHandler(c echo.Context) error {
	globalContext, err := x.getAiStudioSectionGlobals(c, routes.ManageAIModelNodesPage.String())
	if err != nil {
		x.logger.Err(err).Msg("error while getting aistudio section globals")
		return HandleGLobalContextError(c, err)
	}
	var updateSpec, deleteSpec, syncSpec string
	nodeId, err := GetUrlParams(c, globals.AIModelNodeId)
	if err != nil {
		return c.Render(http.StatusBadRequest, "400.html", map[string]interface{}{
			"GlobalContext": globalContext,
		})
	}
	response := &models.AIStudioResponse{
		AIModelNode:  x.grpcClients.AIModelNodesDetails(c, nodeId),
		ActionParams: &models.ResourceManagerParams{},
	}
	if response.AIModelNode != nil {
		mess := &pb.AIModelNodeConfiguratorRequest{
			Action: pb.AIModelNodeConfiguratorActions_PATCH_AIMODEL_NODES,
			Spec:   []*mpb.AIModelNode{response.AIModelNode},
		}
		bbytes, err := grpcops.ProtoYamlMarshal(mess)
		if err != nil {
			x.logger.Err(err).Msg("error while marshaling ai model node spec for update")
			return x.HandleError(c, err, http.StatusBadRequest, globalContext)
		}
		updateSpec = string(bbytes)
		mess.Action = pb.AIModelNodeConfiguratorActions_DELETE_AIMODEL_NODES
		bbytes, err = grpcops.ProtoYamlMarshal(mess)
		if err != nil {
			x.logger.Err(err).Msg("error while marshaling ai model node spec for delete")
			return x.HandleError(c, err, http.StatusBadRequest, globalContext)
		}
		deleteSpec = string(bbytes)
		mess.Action = pb.AIModelNodeConfiguratorActions_SYNC_AI_MODELS
		bbytes, err = grpcops.ProtoYamlMarshal(mess)
		if err != nil {
			x.logger.Err(err).Msg("error while marshaling ai model node spec for sync")
			return x.HandleError(c, err, http.StatusBadRequest, globalContext)
		}
		syncSpec = string(bbytes)
	} else {
		x.logger.Err(err).Msg("error while getting ai model node details")
		return x.HandleError(c, err, http.StatusNotFound, globalContext)
	}

	response.ActionParams.API = fmt.Sprintf("%s/api/v1alpha1/aistudio/models/nodes", pkgs.NetworkConfigManager.ExternalURL)
	response.ActionParams.ActionMap = map[string]string{}
	if response.AIModelNode.Org == globalContext.CurrentOrganization.OrgId && slices.Contains(response.AIModelNode.NodeOwners, globalContext.UserInfo.UserId) {
		response.ActionParams.ActionMap[utils.UPDATE] = updateSpec
		response.ActionParams.ActionMap[utils.DELETE] = deleteSpec
		response.ActionParams.ActionMap[utils.SYNC] = syncSpec
	}
	return c.Render(http.StatusOK, "manageai-modelnodes-detail.html", map[string]interface{}{
		"GlobalContext": globalContext,
		"Response":      response,
		"SectionHeader": response.AIModelNode.Name,
		"AIStudio":      routes.AIStudioHome,
		"ResourceBase":  response.AIModelNode.ResourceBase,
	})
}

func (x *WebappService) ManageAIPromptsHandler(c echo.Context) error {
	response := &models.AIStudioResponse{
		AIPrompts:       x.grpcClients.AIModelPrompts(c),
		ActionParams:    &models.ResourceManagerParams{},
		BackListingLink: c.Request().URL.String(),
	}
	globalContext, err := x.getAiStudioSectionGlobals(c, routes.AIPromptsPage.String())
	if err != nil {
		x.logger.Err(err).Msg("error while getting aistudio section globals")
		return HandleGLobalContextError(c, err)
	}

	response.ActionParams.API = fmt.Sprintf("%s/api/v1alpha1/aistudio/prompts", pkgs.NetworkConfigManager.ExternalURL)
	return c.Render(http.StatusOK, "manageai-prompts.html", map[string]interface{}{
		"GlobalContext":  globalContext,
		"Response":       response,
		"SectionHeader":  "AI Model Prompts",
		"CreateTemplate": GetProtoYamlString(serviceops.AiPromptManagerRequest),
	})
}

func (x *WebappService) ManageAIPromptDetailHandler(c echo.Context) error {
	globalContext, err := x.getAiStudioSectionGlobals(c, routes.AIPromptsPage.String())
	if err != nil {
		x.logger.Err(err).Msg("error while getting aistudio section globals")
		return HandleGLobalContextError(c, err)
	}

	promptId, err := GetUrlParams(c, globals.PromptId)
	if err != nil {
		return c.Render(http.StatusBadRequest, "400.html", map[string]interface{}{
			"GlobalContext": globalContext,
		})
	}
	response := &models.AIStudioResponse{
		AIPrompt:     x.grpcClients.AIModelPromptDetails(c, promptId),
		ActionParams: &models.ResourceManagerParams{},
	}
	if response.AIPrompt != nil {
		if response.AIPrompt.Editable {
			var updateSpec string = ""
			upObj := &pb.PromptManagerRequest{
				Spec:   []*mpb.AIModelPrompt{response.AIPrompt},
				Action: pb.PromptAgentAction_PATCH_PROMPT,
			}
			bytess, err := grpcops.ProtoYamlMarshal(upObj)
			if err != nil {
				x.logger.Err(err).Msg("error while marshaling ai prompt spec for publish")
				return x.HandleError(c, err, http.StatusBadRequest, globalContext)
			}
			updateSpec = string(bytess)
			response.ActionParams.ActionMap = map[string]string{
				utils.UPDATE: updateSpec,
			}
		}
	} else {
		x.logger.Err(err).Msg("error while getting ai prompt details")
		return x.HandleError(c, err, http.StatusNotFound, globalContext)
	}

	response.ActionParams.API = fmt.Sprintf("%s/api/v1alpha1/aistudio/prompts", pkgs.NetworkConfigManager.ExternalURL)
	return c.Render(http.StatusOK, "manageai-prompts-detail.html", map[string]interface{}{
		"GlobalContext": globalContext,
		"Response":      response,
		"SectionHeader": response.AIPrompt.Name,
		"AIStudio":      routes.AIStudioHome,
		"ResourceBase":  response.AIPrompt.ResourceBase,
	})
}

func (x *WebappService) ManageAIAgentsHandler(c echo.Context) error {
	response := &models.AIStudioResponse{
		AIAgents:        x.grpcClients.ListAIAgents(c),
		ActionParams:    &models.ResourceManagerParams{},
		BackListingLink: c.Request().URL.String(),
	}
	globalContext, err := x.getAiStudioSectionGlobals(c, routes.ManageAIAgentsPage.String())
	if err != nil {
		x.logger.Err(err).Msg("error while getting aistudio section globals")
		return HandleGLobalContextError(c, err)
	}

	response.ActionParams.API = fmt.Sprintf("%s/api/v1alpha1/aistudio/agents", pkgs.NetworkConfigManager.ExternalURL)
	return c.Render(http.StatusOK, "manageai-agents.html", map[string]interface{}{
		"GlobalContext":  globalContext,
		"Response":       response,
		"SectionHeader":  "AI Agents",
		"CreateTemplate": GetProtoYamlString(serviceops.AiAgentManagerRequest),
	})
}

func (x *WebappService) ManageAIAgentDetailHandler(c echo.Context) error {
	globalContext, err := x.getAiStudioSectionGlobals(c, routes.ManageAIAgentsPage.String())
	if err != nil {
		x.logger.Err(err).Msg("error while getting aistudio section globals")
		return HandleGLobalContextError(c, err)
	}
	var updateSpec, yamlSpec string
	agentId, err := GetUrlParams(c, globals.AgentId)
	if err != nil {
		return c.Render(http.StatusBadRequest, "400.html", map[string]interface{}{
			"GlobalContext": globalContext,
		})
	}
	response := &models.AIStudioResponse{
		AIAgent:      x.grpcClients.AIAgentDetails(c, agentId),
		ActionParams: &models.ResourceManagerParams{},
	}
	if response.AIAgent != nil {
		upObj := &pb.AgentManagerRequest{
			Spec:   response.AIAgent,
			Action: pb.VapusAIAgentAction_PATCH_AIAGENT,
		}
		bytess, err := grpcops.ProtoYamlMarshal(upObj)
		if err != nil {
			x.logger.Err(err).Msg("error while marshaling ai agent spec for publish")
			return x.HandleError(c, err, http.StatusBadRequest, globalContext)
		}
		updateSpec = string(bytess)
		yBytes, err := grpcops.ProtoYamlMarshal(response.AIAgent)
		if err != nil {
			x.logger.Err(err).Msg("error while marshaling ai agent spec for publish")
		}
		yamlSpec = string(yBytes)
	} else {
		x.logger.Err(err).Msg("error while getting ai agents details")
		return x.HandleError(c, err, http.StatusNotFound, globalContext)
	}
	response.ActionParams.API = fmt.Sprintf("%s/api/v1alpha1/aistudio/agents", pkgs.NetworkConfigManager.ExternalURL)
	if response.AIAgent.Editable && response.AIAgent.Org == globalContext.CurrentOrganization.OrgId && slices.Contains(response.AIAgent.Owners, globalContext.UserInfo.UserId) {
		response.ActionParams.ActionMap = map[string]string{
			utils.UPDATE: updateSpec,
		}
	}
	return c.Render(http.StatusOK, "manageai-agents-detail.html", map[string]interface{}{
		"GlobalContext": globalContext,
		"Response":      response,
		"YamlSpec":      yamlSpec,
		"SectionHeader": response.AIAgent.Name,
		"AgentStudio":   routes.AgentStudioHome,
		"ResourceBase":  response.AIAgent.ResourceBase,
	})
}

func (x *WebappService) ManageAIGuardrailsHandler(c echo.Context) error {
	response := &models.AIStudioResponse{
		AIGuardrails:    x.grpcClients.ListAIGuardrails(c),
		ActionParams:    &models.ResourceManagerParams{},
		BackListingLink: c.Request().URL.String(),
	}
	globalContext, err := x.getAiStudioSectionGlobals(c, routes.ManageAIGuardrailsPage.String())
	if err != nil {
		x.logger.Err(err).Msg("error while getting aistudio section globals")
		return HandleGLobalContextError(c, err)
	}
	response.ActionParams.API = fmt.Sprintf("%s/api/v1alpha1/aistudio/guardrails", pkgs.NetworkConfigManager.ExternalURL)
	return c.Render(http.StatusOK, "manageai-guardrails.html", map[string]interface{}{
		"GlobalContext":  globalContext,
		"Response":       response,
		"SectionHeader":  "AI Guardrails",
		"CreateTemplate": GetProtoYamlString(serviceops.AiGuardrailManagerRequest),
	})
}

func (x *WebappService) ManageAIGuardrailDetailsHandler(c echo.Context) error {
	globalContext, err := x.getAiStudioSectionGlobals(c, routes.ManageAIGuardrailsPage.String())
	if err != nil {
		x.logger.Err(err).Msg("error while getting aistudio section globals")
		return HandleGLobalContextError(c, err)
	}
	var updateSpec, yamlSpec string
	guradrailId, err := GetUrlParams(c, globals.GuardrailId)
	if err != nil {
		return c.Render(http.StatusBadRequest, "400.html", map[string]interface{}{
			"GlobalContext": globalContext,
		})
	}
	response := &models.AIStudioResponse{
		AIGuardrail:  x.grpcClients.DescribeAIGuardrail(c, guradrailId),
		ActionParams: &models.ResourceManagerParams{},
	}
	if response.AIGuardrail != nil {
		upObj := &pb.GuardrailsManagerRequest{
			Spec:   response.AIGuardrail,
			Action: pb.VapusAIGuardrailsAction_PATCH_GUARDRAIL,
		}
		bytess, err := grpcops.ProtoYamlMarshal(upObj)
		if err != nil {
			x.logger.Err(err).Msg("error while marshaling ai guraail spec for publish")
			return x.HandleError(c, err, http.StatusBadRequest, globalContext)
		}
		updateSpec = string(bytess)
		yBytes, err := grpcops.ProtoYamlMarshal(response.AIAgent)
		if err != nil {
			x.logger.Err(err).Msg("error while marshaling ai guardrail spec for publish")
		}
		yamlSpec = string(yBytes)
	} else {
		x.logger.Err(err).Msg("error while getting ai guardrail details")
		return x.HandleError(c, err, http.StatusNotFound, globalContext)
	}
	response.ActionParams.API = fmt.Sprintf("%s/api/v1alpha1/aistudio/guardrails", pkgs.NetworkConfigManager.ExternalURL)
	if response.AIGuardrail.Base.Org == globalContext.CurrentOrganization.OrgId && slices.Contains(response.AIGuardrail.Base.Owners, globalContext.UserInfo.UserId) {
		response.ActionParams.ActionMap = map[string]string{
			utils.UPDATE: updateSpec,
		}
	}
	return c.Render(http.StatusOK, "manageai-guardrail-detail.html", map[string]interface{}{
		"GlobalContext": globalContext,
		"Response":      response,
		"YamlSpec":      yamlSpec,
		"SectionHeader": response.AIGuardrail.Name,
		"ResourceBase":  response.AIGuardrail.ResourceBase,
	})
}
