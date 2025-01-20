package svcops

import (
	"github.com/vapusdata-oss/aistudio/core/models"
)

var ClientRetryLimit = 3

var ClientRetryStart = 0

var (
	AccountsTable          = "accounts"
	UsersTable             = "users"
	AIModelsNodesTable     = "aimodels_nodes"
	OrganizationTable      = "organizations"
	JwtLogsTable           = "jwt_logs"
	UserTeamsTable         = "teams"
	RagTextQueryLogsTable  = "rag_text_query_logs"
	AIModelPromptTable     = "ai_model_prompts"
	RefreshTokenLogsTable  = "refresh_token_logs"
	VapusAIAgentsTable     = "vapus_ai_agents"
	UpDownVotesTable       = "up_down_votes"
	StarReviewsTable       = "star_reviews"
	VapusResourceArnTable  = "vapus_resource_arns"
	PluginsTable           = "plugins"
	AIModelStudioLogsTable = "ai_model_studio_logs"
	AIAgentThreadsTable    = "ai_agent_threads"
	FabricChatLogTable     = "fabric_chat_logs"
	VapusGuardrailsTable   = "vapus_guardrails"
)

var DBTablesMap = map[string]interface{}{
	AccountsTable:          &models.Account{},
	UsersTable:             &models.Users{},
	AIModelsNodesTable:     &models.AIModelNode{},
	JwtLogsTable:           &models.JwtLog{},
	UserTeamsTable:         &models.Team{},
	RagTextQueryLogsTable:  &models.RagTextQueryLog{},
	AIModelPromptTable:     &models.AIModelPrompt{},
	RefreshTokenLogsTable:  &models.RefreshTokenLog{},
	VapusAIAgentsTable:     &models.VapusAIAgent{},
	UpDownVotesTable:       &models.UpDownVote{},
	StarReviewsTable:       &models.StarReview{},
	VapusResourceArnTable:  &models.VapusResourceArn{},
	PluginsTable:           &models.Plugin{},
	AIModelStudioLogsTable: &models.AIModelStudioLog{},
	AIAgentThreadsTable:    &models.AIAgentThread{},
	VapusGuardrailsTable:   &models.AIGuardrails{},
	OrganizationTable:      &models.Organization{},
}
