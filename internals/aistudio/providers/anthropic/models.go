package anthropic

const (
	defaultModel           = "claude-3-5-sonnet-20241022"
	eventStreamContentType = "text/event-stream"
	jsonContentType        = "application/json"
	defaultEndpoint        = "https://api.anthropic.com"
	baseAPIPath            = "/v1"
	defaultAPIVersion      = "2023-06-01"
	EOS                    = "\x00"
	versionHeaderKey       = "anthropic-version"
	apiKeyHeader           = "x-api-key"
	messagePath            = "/messages"
)

type contentType string

const (
	messageTypeText  contentType = "text"
	messageTypeImage contentType = "image"
)

type metadata struct {
	UserID string `json:"user_id"`
}

type message struct {
	Role    string    `json:"role"`
	Content []content `json:"content"`
}

type aerror struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

type content struct {
	Type   contentType    `json:"type"`
	Text   *string        `json:"text,omitempty"`
	Source *contentSource `json:"source,omitempty"`
	Name   *string        `json:"name,omitempty"`
	Input  *string        `json:"input,omitempty"`
}

type contentSource struct {
	Type      string `json:"type"`
	MediaType string `json:"media_type"`
	Data      string `json:"data"`
}

type Tool struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	InputSchema string `json:"input_schema,omitempty"`
}

type GenerativeRequest struct {
	Model         string    `json:"model"`
	Messages      []message `json:"messages"`
	System        string    `json:"system"`
	MaxTokens     int       `json:"max_tokens"`
	Metadata      metadata  `json:"metadata"`
	StopSequences []string  `json:"stop_sequences"`
	Stream        bool      `json:"stream,omitempty"`
	Temperature   float64   `json:"temperature,omitempty"`
	TopP          float64   `json:"top_p,omitempty"`
	TopK          int       `json:"top_k,omitempty"`
	Tools         []Tool    `json:"tools,omitempty"`
}

type usage struct {
	InputTokens             int `json:"input_tokens"`
	OutputTokens            int `json:"output_tokens"`
	CacheCreationInputToken int `json:"cache_creation_input_tokens"`
	CacheReadInputTokens    int `json:"cache_read_input_tokens"`
}

type GenerativeResponse struct {
	HTTPStatusCode    int       `json:"-"`
	acceptContentType string    `json:"-"`
	ID                string    `json:"id"`
	Type              string    `json:"type"`
	Error             aerror    `json:"error"`
	Role              string    `json:"role"`
	Content           []content `json:"content"`
	Model             string    `json:"model"`
	StopReason        *string   `json:"stop_reason"`
	StopSequence      *string   `json:"stop_sequence"`
	Usage             usage     `json:"usage"`
	RawBody           []byte    `json:"-"`
}

type streamResponse struct {
	Event string           `json:"event"`
	Data  *streamEventData `json:"data"`
}

type streamEventData struct {
	Type         string             `json:"type"`
	Index        *int               `json:"index"`
	Delta        *contentDelta      `json:"delta"`
	ContentBlock *contentDelta      `json:"content_block"`
	Message      *eventMessageStart `json:"message"`
}

type eventMessageStart struct {
	ID           string       `json:"id"`
	Type         string       `json:"type"`
	Role         string       `json:"role"`
	Content      []any        `json:"content"`
	Model        string       `json:"model"`
	StopReason   any          `json:"stop_reason"`
	StopSequence any          `json:"stop_sequence"`
	Usage        *streamUsage `json:"usage"`
}

type streamUsage struct {
	InputTokens  int `json:"input_tokens"`
	OutputTokens int `json:"output_tokens"`
}

type contentDelta struct {
	Type         string `json:"type"`
	Text         string `json:"text"`
	StopSequence any    `json:"stop_sequence"`
	StopReason   string `json:"stop_reason"`
}
