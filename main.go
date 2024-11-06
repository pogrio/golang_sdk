package pogr

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
)

type pogrSDK struct {
	clientKey   string
	buildKey    string
	sessionID   string
	httpClient  *http.Client
	baseURL     string
}

type initResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
	Payload struct {
		SessionID string `json:"session_id"`
	} `json:"payload,omitempty"`
}

func NewPOGRSDK(clientKey, buildKey, baseURL string) *pogrSDK {

	if baseURL == "" {
		baseURL = os.Getenv("POGR_BASE_URL")
		if baseURL == "" {
			baseURL = "https://api.pogr.io/v1/intake" 
		}
	}

	return &pogrSDK{
		clientKey:  clientKey,
		buildKey:   buildKey,
		httpClient: &http.Client{},
		baseURL:    baseURL,
	}
}

func (sdk *pogrSDK) InitWithUserJWT(userJWT string) error {
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/init", sdk.baseURL), nil)
	if err != nil {
		return err
	}
	req.Header.Set("POGR_CLIENT", sdk.clientKey)
	req.Header.Set("POGR_BUILD", sdk.buildKey)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", userJWT))

	resp, err := sdk.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var initResp initResponse
	if err := json.NewDecoder(resp.Body).Decode(&initResp); err != nil {
		return err
	}

	if !initResp.Success {
		return errors.New(initResp.Error)
	}

	sdk.sessionID = initResp.Payload.SessionID
	return nil
}

func (sdk *pogrSDK) InitWithAssociationID(associationID string) error {
	data := map[string]string{"association_id": associationID}
	payload, err := json.Marshal(data)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/init", sdk.baseURL), bytes.NewBuffer(payload))
	if err != nil {
		return err
	}
	req.Header.Set("POGR_CLIENT", sdk.clientKey)
	req.Header.Set("POGR_BUILD", sdk.buildKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := sdk.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var initResp initResponse
	if err := json.NewDecoder(resp.Body).Decode(&initResp); err != nil {
		return err
	}

	if !initResp.Success {
		return errors.New(initResp.Error)
	}

	sdk.sessionID = initResp.Payload.SessionID
	return nil
}

type genericResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
}

func (sdk *pogrSDK) EndSession() error {
	if sdk.sessionID == "" {
		return errors.New("no active session to end")
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/end", sdk.baseURL), nil)
	if err != nil {
		return err
	}
	req.Header.Set("INTAKE_SESSION_ID", sdk.sessionID)

	resp, err := sdk.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var response genericResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return err
	}

	if !response.Success {
		return errors.New(response.Error)
	}

	sdk.sessionID = ""
	return nil
}

func (sdk *pogrSDK) TriggerEvent(eventData map[string]interface{}) error {

	return nil
}

func (sdk *pogrSDK) SendData(data map[string]interface{}) error {

	return nil
}

func (sdk *pogrSDK) SendLogs(logs []string) error {

	return nil
}

func (sdk *pogrSDK) SendMetrics(metrics map[string]interface{}) error {

	return nil
}

func (sdk *pogrSDK) MonitorData() error {

	return nil
}