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

type DialogflowResponse_Data_Google_SystemIntent struct {
	Intent string `json:"intent,omitempty"`
	Data   struct {
		Type        string   `json:"@type,omitempty"`
		OptContext  string   `json:"opt_context,omitempty"`
		Permissions []string `json:"permissions,omitempty"`
	} `json:"data,omitempty"`
}

type DialogflowResponse_Data_Google struct {
	ExpectUserResponse bool                                         `json:"expectUserResponse,omitempty"`
	IsSsml             bool                                         `json:"isSsmp,omitempty"`
	NoInputPrompts     []interface{}                                `json:"noInputPrompts,omitempty"`
	SystemIntent       *DialogflowResponse_Data_Google_SystemIntent `json:"systemIntent,omitempty"`
}

type DialogflowResponse_Data struct {
	Google *DialogflowResponse_Data_Google `json:"google,omitempty"`
}

type DialogflowResponse struct {
	Speech        string                   `json:"speech,omitempty"`
	DisplaySpeech string                   `json:"displayText,omitempty"`
	Data          *DialogflowResponse_Data `json:"data,omitempty"`
	ContextOut    []struct {
		Name       string `json:"name,omitempty"`
		Lifespan   int    `json:"lifespan,omitempty"`
		Parameters struct {
		} `json:"parameters,omitempty"`
	} `json:"contextOut,omitempty"`
}
