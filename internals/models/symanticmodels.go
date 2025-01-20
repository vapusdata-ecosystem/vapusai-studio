package models

import (
	"github.com/pgvector/pgvector-go"
)

type VapusSearchResponse struct {
	Results []*VapusSearchItem `json:"results"`
}

type VapusSearchItem struct {
	Description string `json:"description"`
	*BaseIdentifier
}

type QueryAttributes struct {
	Error         string `json:"error"`
	EmptyResponse bool   `json:"empty_response"`
}

type RagTextQueryLog struct {
	VapusBase       `bun:",embed" json:"base,omitempty" yaml:"base,omitempty" toml:"base,omitempty"`
	TextQuery       string `bun:"text_query" json:"text_query,omitempty" yaml:"text_query,omitempty"`
	TextQuerySearch string `bun:"text_query_search,type:tsvector" json:"text_query_search,omitempty" yaml:"text_query_search,omitempty"`
	DataQuery       string `bun:"data_query" json:"data_query,omitempty" yaml:"data_query,omitempty"`
	DataProduct     string `bun:"data_product" json:"data_product,omitempty" yaml:"data_product,omitempty"`
	GenericSqlQuery string `bun:"generic_sql_query" json:"generic_sql_query,omitempty" yaml:"generic_sql_query,omitempty"`
	// Attributes          *QueryAttributes `bun:"attributes,type:jsonb" json:"attributes,omitempty" yaml:"attributes,omitempty"`
	TextQueryEmbeddings pgvector.Vector `bun:"text_query_embeddings,type:vector(1536)"`
	IsValid             bool            `bun:"is_valid" json:"is_valid,omitempty" yaml:"is_valid,omitempty"`
	ErrorText           string          `bun:"error_text" json:"error_text,omitempty" yaml:"error_text,omitempty"`
	ResponseLength      int64           `bun:"response_length" json:"response_length,omitempty" yaml:"response_length,omitempty"`
}

func (dm *RagTextQueryLog) PreSaveCreate(authzClaim map[string]string) {
	if dm == nil {
		return
	}
	dm.PreSaveVapusBase(authzClaim)
}
