package router

import (
	"github.com/labstack/echo/v4"
	"github.com/vapusdata-oss/aistudio/webapp/routes"
	services "github.com/vapusdata-oss/aistudio/webapp/services"
)

func manageAIRoutes(e *echo.Group) {
	aistudio := e.Group(routes.ManageAIGroup)
	aistudio.GET(routes.ManageAIModelNodes, services.WebappServiceManager.ManageAIModelNodesHandler)
	aistudio.GET(routes.ManageAIManageModelNodeResource, services.WebappServiceManager.ManageAIModelNodesDetailHandler)
	aistudio.GET(routes.ManageAIPrompts, services.WebappServiceManager.ManageAIPromptsHandler)
	aistudio.GET(routes.ManageAIPromptResource, services.WebappServiceManager.ManageAIPromptDetailHandler)
	aistudio.GET(routes.ManageAIAgents, services.WebappServiceManager.ManageAIAgentsHandler)
	aistudio.GET(routes.ManageAIAgentsResource, services.WebappServiceManager.ManageAIAgentDetailHandler)
	aistudio.GET(routes.ManageAIGuardrails, services.WebappServiceManager.ManageAIGuardrailsHandler)
	aistudio.GET(routes.ManageAIGuardrailResource, services.WebappServiceManager.ManageAIGuardrailDetailsHandler)
}
