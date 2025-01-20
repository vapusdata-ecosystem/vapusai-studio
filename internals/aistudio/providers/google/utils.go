package googlegenai

import (
	"context"
	"slices"
	"strings"

	"github.com/vapusdata-oss/aistudio/core/models"
	mpb "github.com/vapusdata-oss/apis/protos/models/v1alpha1"
	"google.golang.org/api/iterator"
)

func (o *GoogleGenAI) CrawlModels(ctx context.Context) (result []*models.AIModelBase, err error) {
	modelsIter := o.client.ListModels(ctx)

	for {
		model, err := modelsIter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			continue
		}
		obj := &models.AIModelBase{ // Use the imported type
			ModelId:          model.Name,
			OwnedBy:          "google",
			ModelName:        strings.ReplaceAll(model.Name, "models/", ""),
			Version:          model.Version,
			SupprtedOps:      model.SupportedGenerationMethods,
			InputTokenLimit:  model.InputTokenLimit,
			OutputTokenLimit: model.OutputTokenLimit,
		}
		if slices.Contains(model.SupportedGenerationMethods, "embedContent") {
			obj.ModelType = mpb.AIModelType_EMBEDDING.String()
		} else {
			obj.ModelType = mpb.AIModelType_LLM.String()
		}
		result = append(result, obj)
	}
	return result, nil
}
