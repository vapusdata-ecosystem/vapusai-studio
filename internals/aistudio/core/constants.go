package aicore

import "time"

const (
	USER      = "user"
	SYSTEM    = "system"
	ASSISTANT = "assistant"
	FUNCTION  = "function"
	TOOL      = "tool"
)

type StreamEvent string

const (
	StreamEventStart      StreamEvent = "start"
	StreamEventEnd        StreamEvent = "end"
	StreamEventData       StreamEvent = "data"
	StreamEventError      StreamEvent = "error"
	StreamGuardrailFailed StreamEvent = "guardrail_failed"
)

func (s StreamEvent) String() string {
	return string(s)
}

var retryStatusCodes = map[int]bool{
	429: true,
	500: true,
	502: true,
	503: true,
	504: true,
}

var defaultRetryWaitTime = 2 * time.Second

var StartTagTemplate = `{TAG}`
var EndTagTemplate = `{/TAG}`
