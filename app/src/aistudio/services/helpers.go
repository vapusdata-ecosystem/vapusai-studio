package services

import (
	"context"
	"encoding/json"
	"log"
	"strings"

	"github.com/rs/zerolog"
	aicore "github.com/vapusdata-oss/aistudio/core/aistudio/core"
	aimodels "github.com/vapusdata-oss/aistudio/core/aistudio/providers"
	aitool "github.com/vapusdata-oss/aistudio/core/aistudio/tools"
	models "github.com/vapusdata-oss/aistudio/core/models"
	coreutils "github.com/vapusdata-oss/aistudio/core/utils"
	mpb "github.com/vapusdata-oss/apis/protos/models/v1alpha1"
)

func (x *VapusAINodeManagerAgent) setAiModelNodeCredentials(ctx context.Context, secretName string, modelNode *models.AIModelNode, logger zerolog.Logger) error {
	result, err := coreutils.StructToMap(modelNode.NetworkParams.Credentials)
	if err != nil {
		logger.Err(err).Msgf("error while converting struct to map for default platform jwt secrets.")
		return err
	}
	err = x.dmStore.SecretStore.WriteSecret(ctx, result, secretName)
	if err != nil {
		logger.Err(err).Msgf("error while swapping default platform JWT keys for given resource - %v", secretName)
		return err
	}
	modelNode.NetworkParams.SecretName = secretName
	modelNode.NetworkParams.Credentials = nil
	return nil
}

func crawlAIModels(ctx context.Context, modelNode *models.AIModelNode, logger zerolog.Logger) error {
	aiConfig, err := aimodels.NewAIModelNode(aimodels.WithAIModelNode(modelNode), aimodels.WithLogger(helperLogger))
	if err != nil {
		logger.Err(err).Msgf("error while creating AI model node for model - %v", modelNode.Name)
		return err
	}

	result, err := aiConfig.CrawlModels(ctx)
	if err != nil {
		logger.Err(err).Msgf("error while crawling models for model - %v", modelNode.Name)
		return err
	}
	log.Println("Crawled models for model - ", result)
	modelNode.GenerativeModels = make([]*models.AIModelBase, 0)
	modelNode.EmbeddingModels = make([]*models.AIModelBase, 0)
	for _, model := range result {
		if model.ModelType == mpb.AIModelType_EMBEDDING.String() {
			modelNode.EmbeddingModels = append(modelNode.EmbeddingModels, model)
		} else {
			modelNode.GenerativeModels = append(modelNode.GenerativeModels, model)
		}
	}
	return nil
}

func BuildAIPromptTemplate(obj *models.AIModelPrompt) {
	if obj == nil {
		return
	}
	templaterMap := []map[string]string{}
	template := obj.Prompt.UserMessage + "\n"
	if obj.Prompt.InputTag != "" {
		template = template +
			strings.Replace(aicore.StartTagTemplate, "TAG", obj.Prompt.InputTag, 1) +
			"[" + obj.Prompt.InputTag + "]" +
			strings.Replace(aicore.EndTagTemplate, "TAG", obj.Prompt.InputTag, 1) +
			"\n"
	} else {
		template = template + "[Input] - \n"
	}
	if obj.Prompt.ContextTag != "" {
		template = template +
			strings.Replace(aicore.StartTagTemplate, "TAG", obj.Prompt.ContextTag, 1) +
			"[" + obj.Prompt.ContextTag + "]" +
			strings.Replace(aicore.EndTagTemplate, "TAG", obj.Prompt.ContextTag, 1) +
			"\n"
	}
	if obj.Prompt.OutputTag != "" {
		template = template +
			strings.Replace(aicore.StartTagTemplate, "TAG", obj.Prompt.OutputTag, 1) +
			strings.Replace(aicore.EndTagTemplate, "TAG", obj.Prompt.OutputTag, 1) +
			"\n"
	}
	if obj.Prompt.Sample != nil {
		if obj.Prompt.Sample.InputText != "" {
			template = template + "[Sample Input] " + obj.Prompt.Sample.InputText + "\n"
		}
		if obj.Prompt.Sample.Response != "" {
			template = template + "[Sample Output] " + obj.Prompt.Sample.Response + "\n"
		}
		// template = template + "[Sample Input] " + obj.Prompt.Sample.InputText + "\n"
		// template = template + "[Sample Output] " + obj.Prompt.Sample.Response + "\n"
	}
	sysTemplate := obj.Prompt.SystemMessage
	if obj.Prompt.CustomJson != "" {
		sysTemplate = sysTemplate + "\n Expected JSON output format: " + obj.Prompt.CustomJson
	}
	obj.UserTemplate = template
	obj.Prompt.SystemMessage = sysTemplate
	templaterMap = append(templaterMap, map[string]string{
		"content": obj.UserTemplate,
		"role":    "user",
	})
	templaterMap = append(templaterMap, map[string]string{
		"content": obj.Prompt.SystemMessage,
		"role":    "system",
	})
	if len(obj.Prompt.Tools) > 0 {
		tools := []map[string]interface{}{}
		for _, tool := range obj.Prompt.Tools {
			tools = append(tools, map[string]interface{}{
				"type":   tool.Type,
				"schema": tool.ToolSchema,
			})
		}
		toolBytes, err := json.MarshalIndent(tools, "", "  ")
		if err != nil {
			return
		}
		templaterMap = append(templaterMap, map[string]string{
			"tools": string(toolBytes),
		})
	}
	bbytes, err := json.MarshalIndent(templaterMap, "", "  ")
	if err != nil {
		return
	}
	obj.Template = string(bbytes)
	// template = template + "[Instruction] - Please follow the above instructions to get proper content, input and expected output in desired format.\n Do not change the format of the content, input and expected output."
	return
}

func BuildPromptSchema(obj *models.AIModelPrompt) {
	for _, tool := range obj.Prompt.Tools {
		if tool.AutoGenerate {
			log.Println("Auto generating tool schema for tool - ", tool.Definition)
		}
	}
}

func BuildAgentFuntioncallRenderer(obj *models.VapusAIAgent) string {
	if obj == nil {
		return ""
	}
	required := []string{}
	tool := &models.FunctionCall{
		Name: coreutils.SlugifyBase(obj.Name),
		Parameters: &models.FunctionParameter{
			Type:        aitool.FuncParamType,
			Properties:  map[string]*models.ParameterProperties{},
			Description: obj.Description,
		},
	}
	for _, step := range obj.Steps {
		tool.Parameters.Properties[step.Id] = &models.ParameterProperties{
			Type:        strings.ToLower(step.ValueType),
			Description: step.Prompt,
		}
		if step.Required {
			required = append(required, step.Id)
		}
	}
	tool.RequiredFields = required
	// schema, err := coreutils.GenericMarshaler(tool, mpb.ContentFormats_JSON.String())
	// if err != nil {
	// 	return ""
	// }
	schema, err := json.MarshalIndent(tool, "", "  ")
	if err != nil {
		return ""
	}
	return string(schema)
}

func BuildGuardrailSchema(obj *models.AIGuardrails) string {
	if obj == nil {
		return ""
	}
	tool := &models.FunctionCall{
		Name:        coreutils.SlugifyBase(obj.Name),
		Description: obj.Description,
		Parameters: &models.FunctionParameter{
			Type:        aitool.FuncParamType,
			Properties:  map[string]*models.ParameterProperties{},
			Description: "Guardrails for AI models, based on the different topics and its conent, the content will be filtered.",
		},
	}

	topicParams := &models.ParameterProperties{
		Type: strings.ToLower(mpb.AgentStepValueType_OBJECT.String()),
		Description: `Guardrails for different topics,based on the different topics, the content will be filtered.
		 If the content is not related to the topic, return false for below properties. If they are related to the topic, return true.`,
		Properties: map[string]*models.ParameterProperties{},
	}
	contentParams := &models.ParameterProperties{
		Type: strings.ToLower(mpb.AgentStepValueType_OBJECT.String()),
		Description: `Guardrails for different contents,based on the different contents, the content will be filtered. Related levels & types of content is defined in enums. 
		If user input is not related to the content, return NONE for below properties. If they are related to the content, return the severity.
		Valid values for severity are: NONE, LOW, MEDIUM, HIGH that should be returned based on the content.`,
		Properties: map[string]*models.ParameterProperties{},
	}
	for _, topic := range obj.Topics {
		topicParams.Properties[coreutils.SlugifyBase(topic.Topic)] = &models.ParameterProperties{
			Type:        strings.ToLower(mpb.AgentStepValueType_BOOLEAN.String()),
			Description: topic.Description,
		}
	}
	tool.Parameters.Properties["topic_guardrails"] = topicParams
	contentParams.Properties[coreutils.SlugifyBase("Hate Speech")] = &models.ParameterProperties{
		Type:        strings.ToLower(mpb.AgentStepValueType_STRING.String()),
		Description: "Hate speech is a communication that carries no meaning other than the expression of hatred for some group, especially in circumstances in which the communication is likely to provoke violence.",
	}
	contentParams.Properties[coreutils.SlugifyBase("Insults")] = &models.ParameterProperties{
		Type:        strings.ToLower(mpb.AgentStepValueType_STRING.String()),
		Description: `Insults are words or actions that are intended to be rude or offensive, often because they are directed at a particular person or group.`,
	}
	contentParams.Properties[coreutils.SlugifyBase("Threats")] = &models.ParameterProperties{
		Type:        strings.ToLower(mpb.AgentStepValueType_STRING.String()),
		Description: `Threats are statements of an intention to inflict pain, injury, damage, or other hostile action on someone in retribution for something done or not done.`,
	}
	contentParams.Properties[coreutils.SlugifyBase("Sexual")] = &models.ParameterProperties{
		Type:        strings.ToLower(mpb.AgentStepValueType_STRING.String()),
		Description: `Sexual content is any material depicting, describing, or alluding to sexual behavior or anatomy that is intended to arouse sexual feelings.`,
	}
	contentParams.Properties[coreutils.SlugifyBase("Misconduct")] = &models.ParameterProperties{
		Type:        strings.ToLower(mpb.AgentStepValueType_STRING.String()),
		Description: `Misconduct is behavior that is illegal or dishonest, or that is considered morally wrong by most people.`,
	}
	tool.Parameters.Properties["content_guardrails"] = contentParams
	schema, err := json.MarshalIndent(tool, "", "  ")
	if err != nil {
		return ""
	}
	return string(schema)
}
