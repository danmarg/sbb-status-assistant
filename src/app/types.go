package app

import (
	"encoding/json"
	"time"
)

type DialogflowRequest struct {
	ID        string    `json:"id"`
	Timestamp time.Time `json:"timestamp"`
	Lang      string    `json:"lang"`
	Result    struct {
		Source           string `json:"source"`
		ResolvedQuery    string `json:"resolvedQuery"`
		Action           string `json:"action"`
		ActionIncomplete bool   `json:"actionIncomplete"`
		Parameters       struct {
			Source      string      `json:"source"`
			Destination string      `json:"destination"`
			Transport   []string    `json:"transport"`
			Route       []string    `json:"route"`
			Limit       json.Number `json:"limit"`
			DateTime    string      `json:"date-time"`
			Query       string      `json:"query"`
		} `json:"parameters"`
		Contexts []interface{} `json:"contexts"`
		Metadata struct {
			IntentID                  string `json:"intentId"`
			WebhookUsed               string `json:"webhookUsed"`
			WebhookForSlotFillingUsed string `json:"webhookForSlotFillingUsed"`
			WebhookResponseTime       int    `json:"webhookResponseTime"`
			IntentName                string `json:"intentName"`
		} `json:"metadata"`
		Fulfillment struct {
			Speech   string `json:"speech"`
			Messages []struct {
				Type   int    `json:"type"`
				Speech string `json:"speech"`
			} `json:"messages"`
		} `json:"fulfillment"`
		Score float32 `json:"score"`
	} `json:"result"`
	OriginalRequest struct {
		Data struct {
			User struct {
				UserID  string `json:"user_id"`
				Profile struct {
					DisplayName string `json:"display_name"`
					GivenName   string `json:"given_name"`
					FamilyName  string `json:"family_name"`
				} `json:"profile"`
				AccessToken string `json:"access_token"`
			} `json:"user"`
			Device struct {
				Location struct {
					Coordinates struct {
						Latitude  float64 `json:"latitude"`
						Longitude float64 `json:"longitude"`
					} `json:"coordinates"`
					FormattedAddress string `json:"formatted_address"`
					ZipCode          string `json:"zip_code"`
					City             string `json:"city"`
				} `json:"location"`
			} `json:"device"`
		} `json:"data"`
	} `json:"originalRequest"`
	Status struct {
		Code            int    `json:"code"`
		ErrorType       string `json:"errorType"`
		ErrorDetails    string `json:"errorDetails"`
		WebhookTimedOut bool   `json:"webhookTimedOut"`
	} `json:"status"`
	SessionID string `json:"sessionId"`
}

type DialogflowResponse struct {
	Speech        string `json:"speech"`
	DisplaySpeech string `json:"displayText"`
	Data          struct {
		Google struct {
			ExpectUserResponse bool          `json:"expect_user_response"`
			IsSsml             bool          `json:"is_ssml"`
			NoInputPrompts     []interface{} `json:"no_input_prompts"`
			SystemIntent       struct {
				Intent string `json:"intent"`
				Spec   struct {
					PermissionValueSpec struct {
						OptContext  string   `json:"opt_context"`
						Permissions []string `json:"permissions"`
					} `json:"permission_value_spec"`
				} `json:"spec"`
			} `json:"system_intent"`
		} `json:"google"`
	} `json:"data"`
	ContextOut []struct {
		Name       string `json:"name"`
		Lifespan   int    `json:"lifespan"`
		Parameters struct {
		} `json:"parameters"`
	} `json:"contextOut"`
}
