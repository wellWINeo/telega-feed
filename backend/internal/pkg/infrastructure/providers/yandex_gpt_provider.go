package providers

import (
	"TelegaFeed/internal/pkg/core/entities"
	"TelegaFeed/pkg/myhttp"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type YandexGPTProvider struct {
	http              myhttp.HttpClient
	catalogIdentifier string
}

func (y *YandexGPTProvider) GenerateSummary(ctx context.Context, article *entities.Article) (string, error) {
	payload := yandexGptRequestModel{
		ModelUri: fmt.Sprintf("gpt://%s/%s", y.catalogIdentifier, defaultModel),
		CompletionOptions: completionOptions{
			Stream:      false,
			Temperature: 0.5,
			MaxTokens:   "2000",
			ReasoningOptions: reasoningOptions{
				Mode: ReasoningModeDisabled,
			},
		},
		Messages: []message{
			{
				Role: SystemRole,
				Text: "TODO",
			},
			{
				Role: UserRole,
				Text: fmt.Sprintf("Title: %s\n\n%s", article.Title, article.Url),
			},
		},
	}

	payloadJson, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	request, err := http.NewRequest(
		http.MethodPost,
		"https://llm.api.cloud.yandex.net/foundationModels/v1/completion",
		bytes.NewBuffer(payloadJson),
	)

	if err != nil {
		return "", err
	}

	request.Header.Set("Content-Type", "application/payloadJson")
	request.Header.Set("Authorization", "Bearer "+"IAM_TOKEN") // TODO

	resp, err := y.http.Do(request)
	if err != nil {
		return "", err
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	var gptResponse yandexGptResponseModel

	if err := json.NewDecoder(resp.Body).Decode(&gptResponse); err != nil {
		return "", err
	}

	if len(gptResponse.Result.Alternatives) == 0 {
		return "", errors.New("no alternatives")
	}

	return gptResponse.Result.Alternatives[0].Message.Text, nil
}

func (y *YandexGPTProvider) GenerateDigest(ctx context.Context, articles []*entities.Article) (string, error) {
	//TODO implement me
	panic("implement me")
}
