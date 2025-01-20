package services

import (
	"fmt"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/vapusdata-oss/aistudio/core/globals"
	"github.com/vapusdata-oss/aistudio/webapp/models"
	pkgs "github.com/vapusdata-oss/aistudio/webapp/pkgs"
	routes "github.com/vapusdata-oss/aistudio/webapp/routes"
	mpb "github.com/vapusdata-oss/apis/protos/models/v1alpha1"
)

func (x *WebappService) AIStudioHandler(c echo.Context) error {
	globalContext, err := x.getStudioSectionGlobals(c, routes.AIStudioPage.String())
	if err != nil {
		Logger.Err(err).Msg("error while getting aistudio section globals")
		return HandleGLobalContextError(c, err)
	}
	promptId := c.QueryParam(globals.PromptId)
	aiModelNode := c.QueryParam("aiModelNode")
	log.Println("PomptId ID: ", promptId)
	log.Println("AI Model Node: ", aiModelNode)
	response := &models.AIStudioResponse{
		AIModelNodes: x.grpcClients.AIModelNodes(c),
		AIPrompts:    x.grpcClients.AIModelPrompts(c),
		AIPrompt: &mpb.AIModelPrompt{
			PromptId: promptId,
		},
		AIModelNode: &mpb.AIModelNode{
			Name: aiModelNode,
		},
		ActionParams: &models.ResourceManagerParams{
			API:       fmt.Sprintf("%s/api/v1alpha1/aistudio/chat", pkgs.NetworkConfigManager.ExternalURL),
			StreamAPI: fmt.Sprintf("%s/api/v1alpha1/aistudio/chat-stream", pkgs.NetworkConfigManager.ExternalURL),
		},
	}

	return c.Render(http.StatusOK, "aistudio.html", map[string]interface{}{
		"GlobalContext": globalContext,
		"Response":      response,
		"SectionHeader": "AI Model Interface",
	})
}
