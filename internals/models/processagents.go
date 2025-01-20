package models

import (
	"context"
	fmt "fmt"

	guuid "github.com/google/uuid"
	"github.com/rs/zerolog"
)

// Non DB model
type VapusInterfaceAgentBase struct {
	Ctx       context.Context                       `json:"ctx" yaml:"ctx"`
	CtxClaim  map[string]string                     `json:"ctx_claim" yaml:"ctx_claim"`
	Action    string                                `json:"action" yaml:"action"`
	Logger    zerolog.Logger                        `json:"logger" yaml:"logger"`
	AgentId   string                                `json:"agent_id" yaml:"agent_id"`
	InitAt    int64                                 `json:"init_at" yaml:"init_at"`
	FinishAt  int64                                 `json:"finish_at" yaml:"finish_at"`
	Status    string                                `json:"status" yaml:"status"`
	Siblings  []map[string]*VapusInterfaceAgentBase `json:"siblings" yaml:"siblings"`
	AgentType string                                `json:"agent_type" yaml:"agent_type"`
	MetaData  map[string]interface{}                `json:"meta_data" yaml:"meta_data"`
}

func (x *VapusInterfaceAgentBase) FinalLog() {
	x.Logger.Info().Msgf("%v - %v action %v started at %v and finished at %v with status %v", x.AgentType, x.AgentId, x.Action, x.InitAt, x.FinishAt, x.Status)
	x.SetAgentLog(x.Logger.Info(), "finalLog", fmt.Sprintf("%v - %v action %v started at %v and finished at %v with status %v", x.AgentType, x.AgentId, x.Action, x.InitAt, x.FinishAt, x.Status))
}

func (x *VapusInterfaceAgentBase) GetAgentLogs() map[string]interface{} {
	return x.MetaData
}

func (x *VapusInterfaceAgentBase) SetAgentLog(loggerEvent *zerolog.Event, key, val string) {
	loggerEvent.Msgf("%s -- %v", key, val)
	if x.MetaData == nil {
		x.MetaData = make(map[string]interface{})
	}
	x.MetaData[key] = val
}

func (x *VapusInterfaceAgentBase) GetAction() string {
	return x.Action
}

func (x *VapusInterfaceAgentBase) SetAgentId() {
	if x.AgentId == "" {
		x.AgentId = guuid.New().String()
	}
}

// func (x *VapusInterfaceAgentBase) SetNewGoRCtx(withAuth bool) error {
// 	if withAuth {
// 		nCtx, err := pbtools.SwapNewContextWithAuthToken(x.Ctx)
// 		if err != nil {
// 			x.Logger.Error().Msgf("error swapping context with auth token: %v", err)
// 			return err
// 		}
// 		x.Ctx = nCtx
// 	} else {
// 		x.Ctx = context.Background()
// 	}
// 	return nil
// }
