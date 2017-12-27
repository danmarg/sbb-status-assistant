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
	Status struct {
		Code            int    `json:"code"`
		ErrorType       string `json:"errorType"`
		ErrorDetails    string `json:"errorDetails"`
		WebhookTimedOut bool   `json:"webhookTimedOut"`
	} `json:"status"`
	SessionID string `json:"sessionId"`
}

type DialogflowResponse struct {
	Speech      string `json:"speech"`
	DisplayText string `json:"displayText"`
	Source      string `json:"source"`
}
