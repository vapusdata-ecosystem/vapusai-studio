package services

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/labstack/echo/v4"
	"github.com/vapusdata-oss/aistudio/core/globals"
	"github.com/vapusdata-oss/aistudio/webapp/models"
	pkgs "github.com/vapusdata-oss/aistudio/webapp/pkgs"
	routes "github.com/vapusdata-oss/aistudio/webapp/routes"
	mpb "github.com/vapusdata-oss/apis/protos/models/v1alpha1"
)

func (x *WebappService) AgentStudioHandler(c echo.Context) error {
	globalContext, err := x.getStudioSectionGlobals(c, routes.AgentStudioPage.String())
	if err != nil {
		Logger.Err(err).Msg("error while getting agent Studio section globals")
		return c.Render(http.StatusBadRequest, "400.html", map[string]interface{}{
			"GlobalContext": globalContext,
		})
	}
	agentId := c.QueryParam(globals.AgentId)
	// aiModelNode := c.QueryParam(models.AIModelNodeId)

	response := &models.AgentStudioResponse{
		ActionParams: &models.ResourceManagerParams{
			API:       fmt.Sprintf("%s/api/v1alpha1/aistudio/agents/run", pkgs.NetworkConfigManager.ExternalURL),
			StreamAPI: fmt.Sprintf("%s/api/v1alpha1/aistudio/agents/run", pkgs.NetworkConfigManager.ExternalURL),
		},
		// AIModelNode: &mpb.AIModelNode{
		// 	ModelNodeId: aiModelNode,
		// },
		AIAgent: &mpb.VapusAIAgent{
			AgentId: agentId,
		},
	}
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		response.AIAgents = x.grpcClients.ListAIAgents(c)
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		response.AIModelNodes = x.grpcClients.AIModelNodes(c)
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		response.AIPrompts = x.grpcClients.AIModelPrompts(c)
	}()
	wg.Wait()

	return c.Render(http.StatusOK, "agent-studio.html", map[string]interface{}{
		"GlobalContext": globalContext,
		"Response":      response,
		"SectionHeader": "Agent Studio",
		"StepsEnum":     mpb.AgentStepEnum_name,
	})
}
