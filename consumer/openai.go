package consumer

import (
	"bytes"
	"context"
	"encoding/json"
	"encore.app/producer"
	"encore.dev/pubsub"
	"fmt"
	"github.com/sashabaranov/go-openai"
	"log"
	"net/http"
	"strings"
	"time"
)

var (
	// encore.dev does some clever stuff to inject any secrets (stored in the cloud) at runtime
	secrets struct {
		oiOpenAIAPIKey string
	}
	httpClient *http.Client
	_          = pubsub.NewSubscription(
		producer.OiEvents, "oi-events",
		pubsub.SubscriptionConfig[*producer.OiEvent]{
			Handler: HandleOiEvent,
		},
	)
)

func init() {
	httpClient = &http.Client{
		Timeout: time.Second * 60,
	}
}

func handleOpenAIInteraction(content string) (string, error) {
	client := openai.NewClient(secrets.oiOpenAIAPIKey)

	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: content,
				},
			},
		},
	)
	if err != nil {
		return "", err
	}

	return resp.Choices[0].Message.Content, nil
}

func buildResponseData(text string) ([]byte, error) {
	return json.Marshal(map[string]string{
		"response_type": "in_channel",
		"text":          text,
	})
}

func handleResponseToResponseURL(responseUrl string, text string) error {
	data, err := buildResponseData(text)
	if err != nil {
		return fmt.Errorf("buildResponseData() failed because %v",
			err.Error(),
		)
	}

	req, err := http.NewRequest(http.MethodPost, responseUrl, bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("http.NewRequest() failed because %v",
			err.Error(),
		)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("httpClient.Do() failed because %v",
			err.Error(),
		)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("httpClient.Do() failed because %v",
			resp.Status,
		)
	}

	return nil
}

func HandleOiEvent(ctx context.Context, oiEvent *producer.OiEvent) error {
	log.Printf(">>> HandleOiEvent()")
	defer log.Printf("<<< HandleOiEvent()")

	var err error

	if oiEvent.Text == "" {
		log.Printf("--- sending empty text request to responseUrl")
		err = handleResponseToResponseURL(
			oiEvent.ResponseURL,
			fmt.Sprintf("Oi <@%v>! What mate?", oiEvent.UserID),
		)
		if err != nil {
			log.Printf("!!! error: handleResponseToResponseURL() failed because %v",
				err.Error(),
			)
			return err
		}
		log.Printf("--- sent.")
	}

	log.Printf("--- interacting with OpenAI...")
	content, err := handleOpenAIInteraction(oiEvent.Text)
	if err != nil {
		log.Printf("!!! error: handleOpenAIInteraction() failed because %v",
			err.Error(),
		)
		_ = handleResponseToResponseURL(
			oiEvent.ResponseURL,
			fmt.Sprintf("Oi <@%v>! Sorry mate: %v", oiEvent.UserID, err.Error()),
		)
		return err
	}
	log.Printf("--- interacted; content=%#+v", content)

	log.Printf("--- sending content request to oiEvent.ResponseURL")
	err = handleResponseToResponseURL(
		oiEvent.ResponseURL,
		fmt.Sprintf("<@%v> %v", oiEvent.UserID, strings.TrimSpace(content)),
	)
	if err != nil {
		log.Printf("!!! error: handleResponseToResponseURL() failed because %v",
			err.Error(),
		)
		return err
	}
	log.Printf("--- sent.")

	return nil
}
