package datasvc

import (
	"reflect"
)

func validateType(t reflect.Type, expected reflect.Kind) (reflect.Type, error) {
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != expected {
		return nil, ErrInvalidDestinationType
	}
	return t, nil
}

// func GenerateQueryFromAI(ctx context.Context, vapusServiceClient *vapussvc.VapusServiceInternalClients, params *pkgs.AIQueryGeneratorparams) (string, error) {
// 	var qType string
// 	switch params.DatabaseType {
// 	default:
// 		qType = "SQL"
// 	}
// 	result, err := vapusServiceClient.Chat(ctx, &aipb.ChatRequest{
// 		Model:     params.AimodelParams.GenerativeModel,
// 		ModelNodeId: params.AimodelParams.GenerativeModelNode,
// 		Temperature: 0.2,
// 		Messages:    RenderInputPrompt(strings.Join(allFields, ",")),
// 	}, p.logger, vapussvc.ClientRetryStart)
// 		Actions: []aipb.AIModelNodeAction{aipb.AIModelNodeAction_GENERATE_CONTENT},
// 		Spec: &aipb.GeneratePromptParams{
// 			AiModel:     params.AimodelParams.GenerativeModel,
// 			ModelNodeId: params.AimodelParams.GenerativeModelNode,
// 			Temperature: 0.9,
// 			InputText:   params.TextQuery,
// 			// PromptTemplate: aicore.DataQueryPrompt.String(),
// 			PromptId:      "pr-e5447251-5cee-4102-9bc0-abcd767684d3",
// 			SystemMessage: "Please generate a" + qType + " query for the following data, query should be in proper" + qType + " format",
// 			Contexts: []*mpb.Mapper{
// 				{
// 					Key:   "Database Schema",
// 					Value: params.Schema,
// 				},
// 			},
// 		},
// 	}, params.Logger, vapussvc.ClientRetryStart)

// 	if err != nil {
// 		return "", err
// 	}
// 	if len(result.GetOutput()) > 0 {
// 		return result.GetOutput()[0].GetContent(), nil
// 	}
// 	return "", nil
// }
