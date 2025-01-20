package prompts

import "time"

const (
	USER      = "user"
	SYSTEM    = "system"
	ASSISTANT = "assistant"
	FUNCTION  = "function"
	TOOL      = "tool"
)

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

var DefaultMaxOPTokenLength int32 = 2000
