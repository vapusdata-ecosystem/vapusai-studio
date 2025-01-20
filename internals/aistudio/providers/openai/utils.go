package openaistd

import (
	"context"
	"strings"

	"github.com/vapusdata-oss/aistudio/core/models"
	mpb "github.com/vapusdata-oss/apis/protos/models/v1alpha1"
)

func (o *OpenAI) CrawlModels(ctx context.Context) (result []*models.AIModelBase, err error) {
	openaiModels, err := o.client.ListModels(ctx)
	if err != nil {
		o.log.Err(err).Msg("error while getting models from openai")
		return nil, err
	}
	for _, model := range openaiModels.Models {
		if strings.Contains(model.ID, "embed") {
			result = append(result, &models.AIModelBase{ // Use the imported type
				ModelId:   model.ID,
				OwnedBy:   model.OwnedBy,
				ModelType: mpb.AIModelType_EMBEDDING.String(),
				ModelName: model.ID,
			})
		} else {
			result = append(result, &models.AIModelBase{ // Use the imported type
				ModelId:   model.ID,
				OwnedBy:   model.OwnedBy,
				ModelType: mpb.AIModelType_LLM.String(),
				ModelName: model.ID,
			})
		}
	}
	return result, nil
}
