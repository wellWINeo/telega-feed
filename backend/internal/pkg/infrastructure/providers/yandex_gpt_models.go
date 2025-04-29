package providers

const (
	ReasoningModeDisabled = "DISABLED"

	SystemRole = "system"
	UserRole   = "user"
)

const defaultModel = "yandexgpt-lite"

// request
type yandexGptRequestModel struct {
	ModelUri          string            `json:"modelUri"`
	CompletionOptions completionOptions `json:"completionOptions"`
	Messages          []message         `json:"messages"`
}

type completionOptions struct {
	Stream           bool             `json:"stream"`
	Temperature      float64          `json:"temperature"`
	MaxTokens        string           `json:"maxTokens"`
	ReasoningOptions reasoningOptions `json:"reasoningOptions"`
}

type reasoningOptions struct {
	Mode string `json:"mode"`
}

type message struct {
	Role string `json:"role"`
	Text string `json:"text"`
}

// response
type yandexGptResponseModel struct {
	Result struct {
		Alternatives []struct {
			Message message `json:"message"`
			Status  string  `json:"status"`
		} `json:"alternatives"`
		Usage struct {
			InputTextTokens  string `json:"inputTextTokens"`
			CompletionTokens string `json:"completionTokens"`
			TotalTokens      string `json:"totalTokens"`
		} `json:"usage"`
		ModelVersion string `json:"modelVersion"`
	} `json:"result"`
}
