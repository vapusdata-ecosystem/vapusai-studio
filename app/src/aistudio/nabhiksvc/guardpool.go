package nabhiksvc

import (
	"context"

	dmstores "github.com/vapusdata-oss/aistudio/aistudio/datastoreops"
	gdrl "github.com/vapusdata-oss/aistudio/core/aistudio/guardrails"
	"github.com/vapusdata-oss/aistudio/core/models"
)

var GuardrailPoolManager *GuardrailPool

type GuardrailPool struct {
	ModelGuardRails    map[string][]string
	AccountGuardRails  []string
	GuardrailClientMap map[string]*gdrl.GuardRailClient
	ModelGuardrailMap  map[string][]string
}

func InitGuardrailPool() (*GuardrailPool, error) {
	GuardrailPoolManager = &GuardrailPool{
		ModelGuardRails:    make(map[string][]string),
		AccountGuardRails:  make([]string, 0),
		GuardrailClientMap: make(map[string]*gdrl.GuardRailClient),
	}
	guardrails, err := dmstores.DMStoreManager.ListAIGuardrails(context.Background(), "", nil)
	if err != nil {
		return nil, err
	}
	modelNodes, err := dmstores.DMStoreManager.ListAIModelNodes(context.Background(), "status = 'ACTIVE' AND deleted_at IS NULL ORDER BY created_at DESC", nil)
	if err != nil {
		return nil, err
	}
	for _, guardrail := range guardrails {
		nodePool := make([]*gdrl.GuardModelNodePool, 0)
		attr := dmstores.Account.AIAttributes
		if attr.GuardrailModelNode == "" || attr.GuardrailModel == "" {
			continue
		}
		nodeConn := AIModelNodeConnectionPoolManager.GetConnectionById(attr.GuardrailModelNode)

		nodePool = append(nodePool, &gdrl.GuardModelNodePool{
			Connection: nodeConn,
			IsAccount:  true,
			Model:      attr.GuardrailModel,
		})
		gd := gdrl.New(
			gdrl.WithSpec(guardrail),
			gdrl.WithModelPool(nodePool),
		)
		GuardrailPoolManager.GuardrailClientMap[guardrail.VapusID] = gd
	}
	for _, modelNode := range modelNodes {
		if modelNode.SecurityGuardrails != nil {
			for _, guardrail := range modelNode.SecurityGuardrails.Guardrails {
				if gd, ok := GuardrailPoolManager.GuardrailClientMap[guardrail]; ok {
					GuardrailPoolManager.AddGuardrails(modelNode.VapusID, gd.Guardrail.VapusID)
				}
			}
		}
	}
	return GuardrailPoolManager, nil
}

func (gp *GuardrailPool) AddGuardrails(modelsNode string, guardrailId string) {
	if modelsNode == "" {
		gp.AccountGuardRails = append(gp.AccountGuardRails, guardrailId)
		return
	} else {
		if _, ok := gp.ModelGuardRails[modelsNode]; !ok {
			gp.ModelGuardRails[modelsNode] = make([]string, 0)
		}
		gp.ModelGuardRails[modelsNode] = append(gp.ModelGuardRails[modelsNode], guardrailId)
	}
}

func (gp *GuardrailPool) GetGuardrail(guardrailId string) *gdrl.GuardRailClient {
	if guardrail, ok := gp.GuardrailClientMap[guardrailId]; ok {
		return guardrail
	}
	return nil
}

func (gp *GuardrailPool) UpdateGuardrailPool(guardrail *models.AIGuardrails) {
	nodePool := make([]*gdrl.GuardModelNodePool, 0)
	nodeConn := AIModelNodeConnectionPoolManager.GetConnectionById(guardrail.GuardModel.ModelNodeID)
	if nodeConn == nil {
		return
	}
	nodePool = append(nodePool, &gdrl.GuardModelNodePool{
		Connection: nodeConn,
		IsAccount:  true,
		Model:      guardrail.GuardModel.ModelID,
	})
	gd := gdrl.New(
		gdrl.WithSpec(guardrail),
		gdrl.WithModelPool(nodePool),
	)
	GuardrailPoolManager.GuardrailClientMap[guardrail.VapusID] = gd
	return
}

func (gp *GuardrailPool) RemoveGuardrail(guardrailId string) {
	for i, guardrail := range gp.AccountGuardRails {
		if guardrail == guardrailId {
			gp.AccountGuardRails = append(gp.AccountGuardRails[:i], gp.AccountGuardRails[i+1:]...)
		}
	}
	for _, guardrails := range gp.ModelGuardRails {
		for i, guardrail := range guardrails {
			if guardrail == guardrailId {
				guardrails = append(guardrails[:i], guardrails[i+1:]...)
			}
		}
	}
}
