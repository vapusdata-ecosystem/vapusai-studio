package svcops

import (
	"strings"

	"github.com/vapusdata-oss/aistudio/core/serviceops/grpcops"
	mpb "github.com/vapusdata-oss/apis/protos/models/v1alpha1"
	aipb "github.com/vapusdata-oss/apis/protos/vapusai-studio/v1alpha1"
	pb "github.com/vapusdata-oss/apis/protos/vapusai-studio/v1alpha1"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func MarshalSpecs(spec protoreflect.ProtoMessage) string {
	// req := &mpb.MainRequestSpec{
	// 	ApiVersion: "v1alpha1",
	// 	Kind:       spec.ProtoReflect().Descriptor().Name(),
	// }
	byess, err := grpcops.ProtoYamlMarshal(spec)
	if err != nil {
		return ""
	}
	specStr := string(byess)
	return strings.TrimSpace(specStr)
}

var SpecMap = map[mpb.RequestObjects]string{
	mpb.RequestObjects_VAPUS_DOMAINS:       MarshalSpecs(OrganizationManagerRequest),
	mpb.RequestObjects_VAPUS_AIMODEL_NODES: MarshalSpecs(AinodeConfiguratorRequest),
	mpb.RequestObjects_VAPUS_ACCOUNT:       MarshalSpecs(AccountManagerRequest),
	mpb.RequestObjects_VAPUS_AIPROMPTS:     MarshalSpecs(AiPromptManagerRequest),
	mpb.RequestObjects_VAPUS_AIAGENTS:      MarshalSpecs(AiAgentManagerRequest),
	mpb.RequestObjects_VAPUS_PLUGINS:       MarshalSpecs(PluginManagerRequest),
	mpb.RequestObjects_VAPUS_AIGUARDRAILS:  MarshalSpecs(AiGuardrailManagerRequest),
}

var ResourceActionsMap = map[mpb.RequestObjects][]string{
	mpb.RequestObjects_VAPUS_AIMODEL_NODES: {
		aipb.AIModelNodeConfiguratorActions_CONFIGURE_AIMODEL_NODES.String(),
		aipb.AIModelNodeConfiguratorActions_PATCH_AIMODEL_NODES.String(),
		aipb.AIModelNodeConfiguratorActions_DELETE_AIMODEL_NODES.String(),
		aipb.AIModelNodeConfiguratorActions_SYNC_AI_MODELS.String(),
	},
	mpb.RequestObjects_VAPUS_DATAMARKETPLACE: {},
	mpb.RequestObjects_VAPUS_ACCOUNT: {
		pb.AccountAgentActions_UPDATE_PROFILE.String(),
		pb.AccountAgentActions_CONFIGURE_AISTUDIO_MODEL.String(),
	},
	mpb.RequestObjects_VAPUS_AIPROMPTS: {
		aipb.PromptAgentAction_CONFIGURE_PROMPT.String(),
		aipb.PromptAgentAction_PATCH_PROMPT.String(),
		aipb.PromptAgentAction_ARCHIVE_PROMPT.String(),
	},
	mpb.RequestObjects_VAPUS_AIAGENTS: {
		aipb.VapusAIAgentAction_CONFIGURE_AIAGENT.String(),
		aipb.VapusAIAgentAction_PATCH_AIAGENT.String(),
		aipb.VapusAIAgentAction_ARCHIVE_AIAGENT.String(),
	},
	mpb.RequestObjects_VAPUS_PLUGINS: {
		pb.PluginAgentAction_CONFIGURE_PLUGIN.String(),
		pb.PluginAgentAction_PATCH_PLUGIN.String(),
		pb.PluginAgentAction_TEST_PLUGIN.String(),
	},
	mpb.RequestObjects_VAPUS_AIGUARDRAILS: {
		aipb.VapusAIGuardrailsAction_CONFIGURE_GUARDRAIL.String(),
		aipb.VapusAIGuardrailsAction_PATCH_GUARDRAIL.String(),
		aipb.VapusAIGuardrailsAction_ARCHIVE_GUARDRAIL.String(),
	},
}

var DataSourceCreds *mpb.GenericCredentialObj = &mpb.GenericCredentialObj{}

var AccountManagerRequest *pb.AccountManagerRequest = &pb.AccountManagerRequest{
	Actions: pb.AccountAgentActions_UPDATE_PROFILE,
	Spec: &mpb.Account{
		AiAttributes: &mpb.AIAttributes{},
	},
}

var AinodeConfiguratorRequest *aipb.AIModelNodeConfiguratorRequest = &aipb.AIModelNodeConfiguratorRequest{
	Action: aipb.AIModelNodeConfiguratorActions_CONFIGURE_AIMODEL_NODES,
	Spec: []*mpb.AIModelNode{
		{
			Attributes: &mpb.AIModelNodeAttributes{
				NetworkParams: &mpb.AIModelNodeNetworkParams{
					Credentials: &mpb.GenericCredentialObj{},
				},
			},
		},
	},
}

var AiGuardrailManagerRequest *aipb.GuardrailsManagerRequest = &aipb.GuardrailsManagerRequest{
	Action: aipb.VapusAIGuardrailsAction_CONFIGURE_GUARDRAIL,
	Spec: &mpb.AIGuardrails{
		Contents:         &mpb.ContentGuardrailLevel{},
		Topics:           []*mpb.TopicGuardrails{{}},
		Words:            []*mpb.WordGuardRails{{}},
		SensitiveDataset: []*mpb.SensitiveDataGuardrails{{}},
		Base:             &mpb.Resourcebase{},
		GuardModel:       &mpb.GuardModels{},
	},
}

var AiAgentManagerRequest *aipb.AgentManagerRequest = &aipb.AgentManagerRequest{
	Action: aipb.VapusAIAgentAction_CONFIGURE_AIAGENT,
	Spec: &mpb.VapusAIAgent{
		AiModelMap: []*mpb.AIModelMap{{}},
		Steps: []*mpb.Steps{
			{}},
	},
}

var AiPromptManagerRequest *aipb.PromptManagerRequest = &aipb.PromptManagerRequest{
	Action: aipb.PromptAgentAction_CONFIGURE_PROMPT,
	Spec: []*mpb.AIModelPrompt{
		{
			Prompt: &mpb.Prompt{
				Sample: &mpb.Sample{},
				Tools:  []*mpb.ToolPrompts{{}},
			},
		},
	},
}

var PluginManagerRequest *pb.PluginManagerRequest = &pb.PluginManagerRequest{
	Action: pb.PluginAgentAction_CONFIGURE_PLUGIN,
	Spec: &mpb.Plugin{
		NetworkParams: &mpb.PluginNetworkParams{
			Credentials: DataSourceCreds,
		},
		DynamicParams: []*mpb.Mapper{{}},
	},
}

var OrganizationManagerRequest *pb.OrganizationManagerRequest = &pb.OrganizationManagerRequest{
	Actions: pb.OrganizationAgentActions_CONFIGURE_ORG,
	Spec: &mpb.Organization{
		SecretPasscode:       &mpb.CredentialSalt{},
		BackendSecretStorage: &mpb.BackendStorages{},
		ArtifactStorage:      &mpb.BackendStorages{},
	},
}

var EnumSpecs = map[string]map[string]int32{
	"AuthnMethod":                   mpb.AuthnMethod_value,
	"EncryptionAlgo":                mpb.EncryptionAlgo_value,
	"BEStoreAccessScope":            mpb.BEStoreAccessScope_value,
	"TlsType":                       mpb.TlsType_value,
	"RequestObjects":                mpb.RequestObjects_value,
	"ApiTokenType":                  mpb.ApiTokenType_value,
	"Frequency":                     mpb.Frequency_value,
	"BackendStorageTypes":           mpb.BackendStorageTypes_value,
	"BackendStorageOnboarding":      mpb.BackendStorageOnboarding_value,
	"StorageEngine":                 mpb.StorageEngine_value,
	"DataSourceType":                mpb.DataSourceType_value,
	"ArtifactTypes":                 mpb.ArtifactTypes_value,
	"LLMServiceProvider":            mpb.LLMServiceProvider_value,
	"LLMQueryType":                  mpb.LLMQueryType_value,
	"DataSensitivityClassification": mpb.DataSensitivityClassification_value,
	"ClassifiedTransformerActions":  mpb.ClassifiedTransformerActions_value,
	"ResourceScope":                 mpb.ResourceScope_value,
	"VersionBumpType":               mpb.VersionBumpType_value,
	"OrderBys":                      mpb.OrderBys_value,
	"ContentFormats":                mpb.ContentFormats_value,
	"AIModelNodeHosting":            mpb.AIModelNodeHosting_value,
	"EmbeddingType":                 mpb.EmbeddingType_value,
	"AIAgentTypes":                  mpb.VapusAiAgentTypes_value,
	"IntegrationPluginTypes":        mpb.IntegrationPluginTypes_value,
	"IntegrationPlugins":            mpb.IntegrationPlugins_value,
	"AgentStepEnum":                 mpb.AgentStepEnum_value,
	"EmailSettings":                 mpb.EmailSettings_value,
	"GuardRailLevels":               mpb.GuardRailLevels_value,
	"AgentStepValueType":            mpb.AgentStepValueType_value,
	"AIToolCallType":                mpb.AIToolCallType_value,
	"AIGuardrailScanMode":           mpb.AIGuardrailScanMode_value,
}
