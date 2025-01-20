package models

import (
	fmt "fmt"

	guuid "github.com/google/uuid"
	globals "github.com/vapusdata-oss/aistudio/core/globals"
	utils "github.com/vapusdata-oss/aistudio/core/utils"
	mpb "github.com/vapusdata-oss/apis/protos/models/v1alpha1"
)

type AIModelPrompt struct {
	VapusBase       `bun:",embed" json:"base,omitempty" yaml:"base,omitempty" toml:"base,omitempty"`
	Name            string     `bun:"name,notnull,unique" json:"name,omitempty" yaml:"name"`
	PromptTypes     []string   `bun:"prompt_type,array" json:"promptType,omitempty" yaml:"promptType"`
	PreferredModels []string   `bun:"preferred_models,array" json:"preferredModels,omitempty" yaml:"preferredModels"`
	Editable        bool       `bun:"editable" json:"editable,omitempty" yaml:"editable" default:"true"`
	Prompt          *Prompt    `bun:"prompt,type:jsonb" json:"prompt,omitempty" yaml:"prompt"`
	IsJsonPrompt    bool       `bun:"is_json_prompt" json:"isJsonPrompt,omitempty" yaml:"isJsonPrompt"`
	Labels          []string   `bun:"labels,array" json:"labels,omitempty" yaml:"labels"`
	SpecDigest      *DigestVal `bun:"spec_digest,type:jsonb" json:"specDigest,omitempty" yaml:"specDigest"`
	ResponseFormat  *Mapper    `bun:"response_format,type:jsonb" json:"responseFormat,omitempty" yaml:"responseFormat"`
	UserTemplate    string     `bun:"user_template" json:"userTemplate,omitempty" yaml:"userTemplate"`
	Template        string     `bun:"template" json:"template,omitempty" yaml:"template"`
}

func (m *AIModelPrompt) SetAccountId(accountId string) {
	if m != nil {
		m.OwnerAccount = accountId
	}
}

func (d *AIModelPrompt) SetPromptId() {
	if d == nil {
		return
	}
	d.VapusID = fmt.Sprintf(globals.PROMPT_ID, guuid.New())
}

func (d *AIModelPrompt) PreSaveCreate(authzClaim map[string]string) {
	if d == nil {
		return
	}
	d.PreSaveVapusBase(authzClaim)
}

func (dn *AIModelPrompt) PreSaveUpdate(userId string) {
	if dn == nil {
		return
	}
	dn.UpdatedBy = userId
	dn.UpdatedAt = utils.GetEpochTime()
}

func (dn *AIModelPrompt) ConvertFromPb(pb *mpb.AIModelPrompt) *AIModelPrompt {
	if pb == nil {
		return nil
	}
	vv := &AIModelPrompt{
		Name:            pb.GetName(),
		PromptTypes:     pb.GetPromptTypes(),
		PreferredModels: pb.GetPreferredModels(),
		Editable:        pb.GetEditable(), // PromptScope:     mpb.ResourceScope(mpb.ResourceScope_value[pb.GetPromptScope()]),
		Prompt:          (&Prompt{}).ConvertFromPb(pb.GetPrompt()),
		IsJsonPrompt:    pb.GetIsJsonPrompt(),
		Labels:          pb.GetLabels(),
	}
	vv.Scope = pb.GetScope().String()
	return vv
}

func (dn *AIModelPrompt) ConvertToPb() *mpb.AIModelPrompt {
	if dn == nil {
		return nil
	}
	return &mpb.AIModelPrompt{
		Name:            dn.Name,
		PromptTypes:     dn.PromptTypes,
		PreferredModels: dn.PreferredModels,
		Editable:        dn.Editable,
		Org:             dn.Organization,
		Scope:           mpb.ResourceScope(mpb.ResourceScope_value[dn.Scope]),
		Prompt:          dn.Prompt.ConvertToPb(),
		IsJsonPrompt:    dn.IsJsonPrompt,
		PromptId:        dn.VapusID,
		Labels:          dn.Labels,
		UserTemplate:    dn.UserTemplate,
		Template:        dn.Template,
		ResourceBase:    dn.ConvertToPbBase(),
	}
}

type ToolPrompts struct {
	Definition   *FunctionCall `json:"definition,omitempty" yaml:"definition,omitempty"`
	AutoGenerate bool          `json:"autoGenerate,omitempty" yaml:"autoGenerate,omitempty"`
	SampleJSON   string        `json:"sampleJson,omitempty" yaml:"sampleJson,omitempty"`
	Type         string        `json:"type,omitempty" yaml:"type,omitempty"`
	ToolSchema   string        `json:"toolSchema,omitempty" yaml:"toolSchema,omitempty"`
}

func (d *ToolPrompts) ConvertFromPb(pbTool *mpb.ToolPrompts) *ToolPrompts {
	if pbTool == nil {
		return nil
	}
	return &ToolPrompts{
		ToolSchema:   pbTool.GetToolSchema(),
		AutoGenerate: pbTool.GetAutoGenerate(),
		SampleJSON:   pbTool.GetSampleJson(),
		Definition:   GetFunctionCallFromString(pbTool.GetToolSchema()),
	}
}

func (d *ToolPrompts) ConvertToPb() *mpb.ToolPrompts {
	if d == nil {
		return nil
	}
	return &mpb.ToolPrompts{
		ToolSchema:   d.ToolSchema,
		AutoGenerate: d.AutoGenerate,
		SampleJson:   d.SampleJSON,
	}
}

type Prompt struct {
	SystemMessage string         `json:"systemMessage,omitempty" yaml:"systemMessage,omitempty"`
	UserMessage   string         `json:"userMessage,omitempty" yaml:"userMessage,omitempty"`
	Tools         []*ToolPrompts `json:"tools,omitempty" yaml:"tools,omitempty"`
	InputTag      string         `json:"inputTag,omitempty" yaml:"inputTag,omitempty"`
	OutputTag     string         `json:"outputTag,omitempty" yaml:"outputTag,omitempty"`
	ContextTag    string         `json:"contextTag,omitempty" yaml:"contextTag,omitempty"`
	Sample        *Sample        `json:"sample,omitempty" yaml:"sample,omitempty"`
	CustomJson    string         `json:"customJson,omitempty" yaml:"customJson,omitempty"`
}

func (d *Prompt) ConvertFromPb(pbPrompt *mpb.Prompt) *Prompt {
	if pbPrompt == nil {
		return nil
	}
	return &Prompt{
		SystemMessage: pbPrompt.GetSystemMessage(),
		UserMessage:   pbPrompt.GetUserMessage(),
		Tools: func() []*ToolPrompts {
			var tools []*ToolPrompts
			for _, tool := range pbPrompt.GetTools() {
				tools = append(tools, (&ToolPrompts{}).ConvertFromPb(tool))
			}
			return tools
		}(),
		InputTag:   pbPrompt.GetInputTag(),
		OutputTag:  pbPrompt.GetOutputTag(),
		ContextTag: pbPrompt.GetContextTag(),
		Sample:     (&Sample{}).ConvertFromPb(pbPrompt.GetSample()),
		CustomJson: pbPrompt.GetCustomJson(),
	}
}

func (d *Prompt) ConvertToPb() *mpb.Prompt {
	if d == nil {
		return nil
	}
	return &mpb.Prompt{
		SystemMessage: d.SystemMessage,
		UserMessage:   d.UserMessage,
		Tools: func() []*mpb.ToolPrompts {
			var tools []*mpb.ToolPrompts
			for _, tool := range d.Tools {
				tools = append(tools, tool.ConvertToPb())
			}
			return tools
		}(),
		InputTag:   d.InputTag,
		OutputTag:  d.OutputTag,
		ContextTag: d.ContextTag,
		Sample:     d.Sample.ConvertToPb(),
		CustomJson: d.CustomJson,
	}
}

type Sample struct {
	InputText string `json:"inputText,omitempty" yaml:"inputText,omitempty"`
	Response  string `json:"response,omitempty" yaml:"response,omitempty"`
}

func (d *Sample) ConvertFromPb(pbSample *mpb.Sample) *Sample {
	if pbSample == nil {
		return nil
	}
	return &Sample{
		InputText: pbSample.GetInputText(),
		Response:  pbSample.GetResponse(),
	}
}

func (d *Sample) ConvertToPb() *mpb.Sample {
	if d == nil {
		return nil
	}
	return &mpb.Sample{
		InputText: d.InputText,
		Response:  d.Response,
	}
}

type PromptTag struct {
	Start string `json:"start,omitempty" yaml:"start,omitempty"`
	End   string `json:"end,omitempty" yaml:"end,omitempty"`
}

func (d *PromptTag) ConvertFromPb(pbTag *mpb.PromptTag) *PromptTag {
	if pbTag == nil {
		return nil
	}
	return &PromptTag{
		Start: pbTag.GetStart(),
		End:   pbTag.GetEnd(),
	}
}

func (d *PromptTag) ConvertToPb() *mpb.PromptTag {
	if d == nil {
		return nil
	}
	return &mpb.PromptTag{
		Start: d.Start,
		End:   d.End,
	}
}
