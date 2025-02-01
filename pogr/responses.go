package pogr

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

// Response types
type initResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
	Payload struct {
		SessionID string `json:"session_id"`
	} `json:"payload,omitempty"`
}

type dataResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
	Payload struct {
		DataID string `json:"data_id"`
	} `json:"payload,omitempty"`
}

type genericResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
}

// handleInitResponse processes initialization responses
func (sdk *pogrSDK) handleInitResponse(req *Request) (string, error) {
	resp, err := sdk.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to execute request: %w", err)
	}

	var initResp initResponse
	if err := json.Unmarshal(resp.Body, &initResp); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	if !initResp.Success {
		return "", errors.New(initResp.Error)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	sdk.setSessionState(initResp.Payload.SessionID, true)
	return initResp.Payload.SessionID, nil
}

// handleDataResponse processes data submission responses
func (sdk *pogrSDK) handleDataResponse(req *Request) (string, error) {
	resp, err := sdk.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to execute request: %w", err)
	}

	var dataResp dataResponse
	if err := json.Unmarshal(resp.Body, &dataResp); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	if !dataResp.Success {
		return "", errors.New(dataResp.Error)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return dataResp.Payload.DataID, nil
}

// handleGenericResponse processes general responses
func (sdk *pogrSDK) handleGenericResponse(req *Request) error {
	resp, err := sdk.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}

	var genResp genericResponse
	if err := json.Unmarshal(resp.Body, &genResp); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	if !genResp.Success {
		return errors.New(genResp.Error)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}
