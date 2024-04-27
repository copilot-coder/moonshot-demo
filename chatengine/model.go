package chatengine

import "errors"

var (
	ErrNoResponse = errors.New("no response")
	ErrRateLimit  = errors.New("rate limit")
)

type Config struct {
	ApiKey string
}

type ChatReq struct {
	Model     string    `json:"model"`
	Messages  []Message `json:"messages"`
	MaxTokens int       `json:"max_tokens"`
	Stream    bool      `json:"stream"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatResponse struct {
	Id      string       `json:"id"`
	Choices []ChoiceItem `json:"choices"`
	Err     error        `json:"-"`
}

type ChoiceItem struct {
	Message Message `json:"message"`
	Delta   struct {
		Content string `json:"content"`
	} `json:"delta"`
	FinishReason string `json:"finish_reason"`
	Usage        struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}
