package producer

import (
	"context"
	"encoding/json"
	"encore.dev/pubsub"
	"fmt"
	"github.com/slack-go/slack"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type OiEvent struct {
	UserID      string
	Text        string
	ResponseURL string
}

var OiEvents = pubsub.NewTopic[*OiEvent]("oi-events", pubsub.TopicConfig{
	DeliveryGuarantee: pubsub.AtLeastOnce,
})

var (
	// encore.dev does some clever stuff to inject any secrets (stored in the cloud) at runtime
	secrets struct {
		oiSlackSigningSecret string
	}
)

func handleVerifySecret(req *http.Request) ([]byte, error) {
	if req.Header.Get("X-Test-Mode") == "true" {
		body, err := io.ReadAll(req.Body)
		if err != nil {
			return nil, err
		}

		return body, nil
	}

	sv, err := slack.NewSecretsVerifier(req.Header, secrets.oiSlackSigningSecret)
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}

	_, err = sv.Write(body)
	if err != nil {
		return nil, err
	}

	err = sv.Ensure()
	if err != nil {
		return nil, err
	}

	return body, nil
}

func buildResponseData(outputText string) ([]byte, error) {
	return json.Marshal(map[string]string{
		"response_type": "in_channel",
		"text":          outputText,
	})
}

func handleResponseToRequest(w http.ResponseWriter, outputText string, status int) {
	data, err := buildResponseData(outputText)
	if err != nil {
		log.Printf("!!! error: buildResponseData() failed because %v",
			err.Error(),
		)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_, err = w.Write(data)
	return
}

// Oi is our handler; TODO: remove "raw" while in "encore run" and see how hard it is to fault find
//
//encore:api public raw path=/oi
func Oi(w http.ResponseWriter, req *http.Request) {
	panic("oh no")

	log.Printf("%#+v", req)

	log.Printf(">>> Oi()")
	defer log.Printf("<<< Oi()")

	log.Printf("--- verifying secret...")
	body, err := handleVerifySecret(req)
	if err != nil {
		log.Printf("!!! error: handleVerifySecret() failed because %v",
			err.Error(),
		)
		return
	}
	log.Printf("--- verified.")

	log.Printf("--- parsing query parameters...")
	q, err := url.ParseQuery(string(body))
	if err != nil {
		log.Printf("!!! error: url.ParseQuery() failed because %v",
			err.Error(),
		)
		return
	}
	log.Printf("--- parsed; q=%#+v", q)

	// token=[removed]
	_ = q.Get("token")

	// team_id=[removed]
	_ = q.Get("team_id")

	// team_domain=ftphub
	_ = q.Get("team_domain")

	// channel_id=[removed]
	_ = q.Get("channel_id")

	// channel_name=directmessage
	_ = q.Get("channel_name")

	// user_id=[removed]
	userID := q.Get("user_id")

	// user_name=initialed85
	_ = q.Get("user_name")

	// command=%2Foi
	_ = q.Get("command")

	// text=test
	text := strings.TrimSpace(q.Get("text"))

	// api_app_id=[removed]
	_ = q.Get("api_app_id")

	// is_enterprise_install=false
	_ = q.Get("is_enterprise_install")

	// response_url=https%3A%2F%2Fhooks.slack.com%2Fcommands%2FT2UR9F7MM%2F5178155597170%2FNi7g2rzaDmkMp5OS9sXQ6kmr
	responseUrl := strings.TrimSpace(q.Get("response_url"))

	// trigger_id=5163612806935.96859517735.c4cf9453c51011f5cd1f7ac4b55d3d65
	_ = q.Get("trigger_id")

	if responseUrl == "" {
		log.Printf("!!! error: responseUrl is empty; cannot continue")
		return
	}

	quotedText := ""
	for _, line := range strings.Split(text, "\n") {
		quotedText += fmt.Sprintf("> %v\n", line)
	}

	log.Printf("--- sending immediate 200...")
	handleResponseToRequest(
		w,
		fmt.Sprintf("<@%v> asked:\n\n%v", userID, quotedText),
		http.StatusOK,
	)
	log.Printf("--- sent.")

	log.Printf("--- text=%#+v", text)
	log.Printf("--- responseUrl=%#+v", responseUrl)

	log.Printf("--- publishing event...")
	_, err = OiEvents.Publish(
		context.Background(),
		&OiEvent{
			UserID:      userID,
			Text:        text,
			ResponseURL: responseUrl,
		},
	)
	if err != nil {
		log.Printf("!!! error: failed to publish event because %v", err.Error())
		return
	}
	log.Printf("--- published.")
}
