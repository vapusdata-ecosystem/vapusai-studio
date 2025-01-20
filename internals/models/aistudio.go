package models

import (
	"context"
	"encoding/json"
	fmt "fmt"

	guuid "github.com/google/uuid"
	"github.com/uptrace/bun"
	"github.com/vapusdata-oss/aistudio/core/globals"
	utils "github.com/vapusdata-oss/aistudio/core/utils"
	mpb "github.com/vapusdata-oss/apis/protos/models/v1alpha1"
)

type AIModelNode struct {
	VapusBase             `bun:",embed" json:"base,omitempty" yaml:"base,omitempty" toml:"base,omitempty"`
	Name                  string                    `bun:"name,notnull,unique" json:"name,omitempty" yaml:"name"`
	GenerativeModels      []*AIModelBase            `bun:"generative_models,type:jsonb" json:"generativeModels,omitempty" yaml:"generativeModels"`
	EmbeddingModels       []*AIModelBase            `bun:"embedding_models,type:jsonb" json:"embeddingModels,omitempty" yaml:"embeddingModels"`
	DiscoverModels        bool                      `bun:"discover_models" json:"discoverModels,omitempty" yaml:"discoverModels"`
	NetworkParams         *AIModelNodeNetworkParams `bun:"network_params,type:jsonb" json:"networkParams,omitempty" yaml:"networkParams"`
	ApprovedOrganizations []string                  `bun:"approved_organizations,array" json:"approvedOrganizations,omitempty" yaml:"approvedOrganizations"`
	Hosting               string                    `bun:"hosting" json:"hosting,omitempty" yaml:"hosting"`
	ServiceProvider       string                    `bun:"service_provider" json:"serviceProvider,omitempty" yaml:"serviceProvider"`
	SecurityGuardrails    *SecurityGuardrails       `bun:"security_guardrails,type:jsonb" json:"securityGuardrails,omitempty" yaml:"securityGuardrails"`
}

func (m *AIModelNode) SetAccountId(accountId string) {
	if m != nil {
		m.OwnerAccount = accountId
	}
}

func (m *AIModelNode) ConvertToPb() *mpb.AIModelNode {
	if m != nil {
		obj := &mpb.AIModelNode{
			ModelNodeId:        m.VapusID,
			Name:               m.Name,
			NodeOwners:         m.Editors,
			Status:             m.Status,
			Org:                m.Organization,
			SecurityGuardrails: m.SecurityGuardrails.ConvertToPb(),
			ResourceBase:       m.ConvertToPbBase(),
			Attributes: &mpb.AIModelNodeAttributes{
				DiscoverModels:   m.DiscoverModels,
				GenerativeModels: make([]*mpb.AIModelBase, 0),
				EmbeddingModels:  make([]*mpb.AIModelBase, 0),
				NetworkParams:    m.NetworkParams.ConvertToPb(),
				Scope:            mpb.ResourceScope(mpb.ResourceScope_value[m.Scope]),
				Hosting:          mpb.AIModelNodeHosting(mpb.AIModelNodeHosting_value[m.Hosting]),
				ServiceProvider:  mpb.LLMServiceProvider(mpb.LLMServiceProvider_value[m.ServiceProvider]),
			},
		}
		for _, gm := range m.GenerativeModels {
			obj.Attributes.GenerativeModels = append(obj.Attributes.GenerativeModels, gm.ConvertToPb())
		}
		for _, em := range m.EmbeddingModels {
			obj.Attributes.EmbeddingModels = append(obj.Attributes.EmbeddingModels, em.ConvertToPb())
		}
		return obj
	}
	return nil
}

func (m *AIModelNode) ConvertFromPb(pb *mpb.AIModelNode) *AIModelNode {
	if pb == nil {
		return nil
	}
	obj := &AIModelNode{
		Name:               pb.GetName(),
		DiscoverModels:     pb.GetAttributes().GetDiscoverModels(),
		GenerativeModels:   make([]*AIModelBase, 0),
		EmbeddingModels:    make([]*AIModelBase, 0),
		NetworkParams:      (&AIModelNodeNetworkParams{}).ConvertFromPb(pb.GetAttributes().GetNetworkParams()),
		Hosting:            pb.GetAttributes().GetHosting().String(),
		ServiceProvider:    pb.GetAttributes().GetServiceProvider().String(),
		SecurityGuardrails: (&SecurityGuardrails{}).ConvertFromPb(pb.GetSecurityGuardrails()),
	}
	obj.Scope = pb.GetAttributes().GetScope().String()
	obj.Editors = pb.GetNodeOwners()
	for _, gm := range pb.GetAttributes().GetGenerativeModels() {
		obj.GenerativeModels = append(obj.GenerativeModels, (&AIModelBase{}).ConvertFromPb(gm))
	}
	for _, em := range pb.GetAttributes().GetEmbeddingModels() {
		obj.EmbeddingModels = append(obj.EmbeddingModels, (&AIModelBase{}).ConvertFromPb(em))
	}
	return obj
}

func (n *AIModelNode) GetModelNodeId() string {
	if n != nil {
		return n.VapusID
	}
	return ""
}

func (n *AIModelNode) GetName() string {
	if n != nil {
		return n.Name
	}
	return ""
}

func (n *AIModelNode) GetNodeOwners() []string {
	if n != nil {
		return n.Editors
	}
	return nil
}

func (n *AIModelNode) GetStatus() string {
	if n != nil {
		return n.Status
	}
	return ""
}

func (n *AIModelNode) GetGenerativeModels() []*AIModelBase {
	if n != nil {
		return n.GenerativeModels
	}
	return nil
}

func (n *AIModelNode) GetEmbeddingModels() []*AIModelBase {
	if n != nil {
		return n.EmbeddingModels
	}
	return nil
}

func (n *AIModelNode) GetDiscoverModels() bool {
	if n != nil {
		return n.DiscoverModels
	}
	return false
}

func (n *AIModelNode) GetNetworkParams() *AIModelNodeNetworkParams {
	if n != nil {
		return n.NetworkParams
	}
	return nil
}

func (n *AIModelNode) GetScope() string {
	if n != nil {
		return n.Scope
	}
	return ""
}

func (n *AIModelNode) GetApprovedOrganizations() []string {
	if n != nil {
		return n.ApprovedOrganizations
	}
	return nil
}

func (n *AIModelNode) GetHosting() string {
	if n != nil {
		return n.Hosting
	}
	return ""
}

func (n *AIModelNode) GetServiceProvider() string {
	if n != nil {
		return n.ServiceProvider
	}
	return ""
}

type AIModelNodeNetworkParams struct {
	Url                 string                  `json:"url,omitempty" yaml:"url"`
	ApiVersion          string                  `json:"apiVersion,omitempty" yaml:"apiVersion"`
	LocalPath           string                  `json:"localPath,omitempty" yaml:"localPath"`
	Credentials         *GenericCredentialModel `json:"credentials,omitempty" yaml:"credentials"`
	SecretName          string                  `json:"secretName,omitempty" yaml:"secretName"`
	IsAlreadyInSecretBs bool                    `json:"isAlreadyInSecretBS,omitempty" yaml:"isAlreadyInSecretBS"`
}

func (n *AIModelNodeNetworkParams) GetUrl() string {
	if n != nil {
		return n.Url
	}
	return ""
}

func (n *AIModelNodeNetworkParams) GetApiVersion() string {
	if n != nil {
		return n.ApiVersion
	}
	return ""
}

func (n *AIModelNodeNetworkParams) GetLocalPath() string {
	if n != nil {
		return n.LocalPath
	}
	return ""
}

func (n *AIModelNodeNetworkParams) GetCredentials() *GenericCredentialModel {
	if n != nil {
		return n.Credentials
	}
	return nil
}

func (n *AIModelNodeNetworkParams) GetSecretName() string {
	if n != nil {
		return n.SecretName
	}
	return ""
}

func (n *AIModelNodeNetworkParams) GetIsAlreadyInSecretBs() bool {
	if n != nil {
		return n.IsAlreadyInSecretBs
	}
	return false
}

func (m *AIModelNodeNetworkParams) ConvertToPb() *mpb.AIModelNodeNetworkParams {
	if m != nil {
		return &mpb.AIModelNodeNetworkParams{
			Url:                 m.Url,
			ApiVersion:          m.ApiVersion,
			LocalPath:           m.LocalPath,
			Credentials:         m.Credentials.ConvertToPb(),
			SecretName:          m.SecretName,
			IsAlreadyInSecretBs: m.IsAlreadyInSecretBs,
		}
	}
	return nil
}

func (m *AIModelNodeNetworkParams) ConvertFromPb(pb *mpb.AIModelNodeNetworkParams) *AIModelNodeNetworkParams {
	if pb == nil {
		return nil
	}
	return &AIModelNodeNetworkParams{
		Url:                 pb.GetUrl(),
		ApiVersion:          pb.GetApiVersion(),
		LocalPath:           pb.GetLocalPath(),
		Credentials:         (&GenericCredentialModel{}).ConvertFromPb(pb.GetCredentials()),
		SecretName:          pb.GetSecretName(),
		IsAlreadyInSecretBs: pb.GetIsAlreadyInSecretBs(),
	}
}

type AIModelBase struct {
	ModelName        string   `json:"modelName,omitempty" yaml:"modelName,omitempty"`
	ModelId          string   `json:"modelId,omitempty" yaml:"modelId,omitempty"`
	ModelType        string   `json:"modelType,omitempty" yaml:"modelType,omitempty"`
	OwnedBy          string   `json:"ownedBy,omitempty" yaml:"ownedBy,omitempty"`
	InputTokenLimit  int32    `json:"inputTokenLimit,omitempty" yaml:"inputTokenLimit,omitempty"`
	OutputTokenLimit int32    `json:"outputTokenLimit,omitempty" yaml:"outputTokenLimit,omitempty"`
	SupprtedOps      []string `json:"supportedOps,omitempty" yaml:"supportedOps,omitempty"`
	Version          string   `json:"version,omitempty" yaml:"version,omitempty"`
}

func (m *AIModelBase) ConvertToPb() *mpb.AIModelBase {
	if m != nil {
		return &mpb.AIModelBase{
			ModelName: m.ModelName,
			ModelId:   m.ModelId,
			ModelType: mpb.AIModelType(mpb.AIModelType_value[m.ModelType]),
			OwnedBy:   m.OwnedBy,
		}
	}
	return nil
}

func (m *AIModelBase) ConvertFromPb(pb *mpb.AIModelBase) *AIModelBase {
	if pb == nil {
		return nil
	}
	return &AIModelBase{
		ModelName: pb.GetModelName(),
		ModelId:   pb.GetModelId(),
		ModelType: pb.GetModelType().String(),
		OwnedBy:   pb.GetOwnedBy(),
	}
}

func (dn *AIModelNode) SetAINodeId() {
	if dn == nil {
		return
	}
	if dn.VapusID == "" {
		dn.VapusID = fmt.Sprintf(globals.VAPUS_AIMODEL_NODE_ID, guuid.New())
	}
}

func (dn *AIModelNode) PreSaveCreate(authzClaim map[string]string) {
	if dn == nil {
		return
	}
	dn.PreSaveVapusBase(authzClaim)
}

func (dn *AIModelNode) PreSaveUpdate(userId string) {
	if dn == nil {
		return
	}
	dn.UpdatedBy = userId
	dn.UpdatedAt = utils.GetEpochTime()
}

func (dn *AIModelNode) PreSaveDelete(authzClaim map[string]string) {
	if dn == nil {
		return
	}
	dn.PreDeleteVapusBase(authzClaim)
}

func (dn *AIModelNode) GetCredentials(param string) *GenericCredentialModel {
	if dn == nil {
		return nil
	}
	return dn.NetworkParams.Credentials
}

func (dn *AIModelNode) Delete(userId string) {
	if dn == nil {
		return
	}
	dn.DeletedBy = userId
	dn.DeletedAt = utils.GetEpochTime()
}

type VectorEmbeddings struct {
	Vectors32 []float32 `json:"vectors32,omitempty" yaml:"vectors32"`
	Vectors64 []float64 `json:"vectors64,omitempty" yaml:"vectors64"`
}

type AIModelStudioLog struct {
	VapusBase      `bun:",embed" json:"base,omitempty" yaml:"base,omitempty" toml:"base,omitempty"`
	Organization   string `bun:"organization" json:"organization,omitempty" yaml:"organization"`
	Input          string `bun:"input,type:text" json:"Input,omitempty" yaml:"Input"`
	ParsedInput    string `bun:"parsed_input,type:text" json:"parsedInput,omitempty" yaml:"parsedInput"`
	Search         string `bun:"search,type:tsvector" json:"requestDataSearch,omitempty" yaml:"requestDataSearch"`
	ResponseStatus string `bun:"response_status" json:"responseStatus,omitempty" yaml:"responseStatus"`
	ResponseLength int64  `bun:"response_length" json:"responseLength,omitempty" yaml:"responseLength"`
	ThreadId       string `bun:"thread_id" json:"threadId,omitempty" yaml:"threadId"`
	Error          string `bun:"error" json:"error,omitempty" yaml:"error"`
	Mode           string `bun:"mode" json:"mode,omitempty" yaml:"mode,omitempty"`
	Output         string `bun:"output,type:text" json:"output,omitempty" yaml:"output"`
	// LogVectors     pgvector.Vector `bun:"log_vectors,type:vector(1536)" json:"logVectors,omitempty" yaml:"logVectors,omitempty"`
}

var _ bun.AfterCreateTableHook = (*AIModelStudioLog)(nil)

func (*AIModelStudioLog) AfterCreateTable(ctx context.Context, query *bun.CreateTableQuery) error {
	var err error
	_, err = query.DB().NewCreateIndex().IfNotExists().
		Model((*AIModelStudioLog)(nil)).TableExpr("ai_model_studio_log").
		Index("search_idx").
		ColumnExpr("search").
		Using("gin").
		Exec(ctx)
	return err
}

func (dm *AIModelStudioLog) PreSaveCreate(authzClaim map[string]string) {
	if dm == nil {
		return
	}
	dm.PreSaveVapusBase(authzClaim)
}

func (dm *AIModelStudioLog) PreSaveUpdate(userId string) {
	if dm == nil {
		return
	}
	dm.UpdatedBy = userId
	dm.UpdatedAt = utils.GetEpochTime()
}

type FunctionCall struct {
	Name           string             `json:"name" yaml:"name"`
	Parameters     *FunctionParameter `json:"parameters" yaml:"parameters"`
	Description    string             `json:"description" yaml:"description"`
	RequiredFields []string           `json:"requiredFields" yaml:"requiredFields"`
}

func GetFunctionCallFromString(val string) *FunctionCall {
	toolSchema := &FunctionCall{}
	err := json.Unmarshal([]byte(val), toolSchema)
	if err != nil {
		return nil
	}
	return toolSchema
}

func (f *FunctionCall) GetStringParamSchema() string {
	if f != nil {
		bbytes, err := json.MarshalIndent(f.Parameters, "", "  ")
		if err != nil {
			return ""
		}
		return string(bbytes)
	}
	return ""
}

type FunctionParameter struct {
	Type        string                          `json:"type" yaml:"type"`
	Properties  map[string]*ParameterProperties `json:"properties" yaml:"properties"`
	Description string                          `json:"description" yaml:"description"`
}
type ParameterProperties struct {
	Type        string                          `json:"type" yaml:"type"`
	Description string                          `json:"description" yaml:"description"`
	Properties  map[string]*ParameterProperties `json:"properties,omitempty" yaml:"properties,omitempty"`
}

type SecurityGuardrails struct {
	Guardrails []string `json:"guardrails,omitempty" yaml:"guardrails,omitempty"`
}

func (s *SecurityGuardrails) ConvertFromPb(pb *mpb.SecurityGuardrails) *SecurityGuardrails {
	if pb == nil {
		return nil
	}
	return &SecurityGuardrails{
		Guardrails: pb.GetGuardrails(),
	}
}

func (s *SecurityGuardrails) ConvertToPb() *mpb.SecurityGuardrails {
	if s != nil {
		return &mpb.SecurityGuardrails{
			Guardrails: s.Guardrails,
		}
	}
	return nil
}
