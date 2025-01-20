package tools

import (
	"github.com/rs/zerolog"
	svcops "github.com/vapusdata-oss/aistudio/core/serviceops"
	mpb "github.com/vapusdata-oss/apis/protos/models/v1alpha1"
)

type ToolPropertiesType string

func (t ToolPropertiesType) String() string {
	return string(t)
}

const (
	FuncParamType                    = "object"
	StringProp    ToolPropertiesType = "string"
)

type Options struct {
	Input          string
	Sysmess        string
	Tools          []*mpb.ToolCall
	ActionAnalyzer bool
	ParamAnalyzer  bool
	FetchData      bool
	RetryCount     int `default:"0"`
}

type ToolCaller struct {
	AIModel     string
	AIModelNode string
	Client      *svcops.VapusSvcInternalClients
	Logger      zerolog.Logger
}
