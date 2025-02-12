package pogr

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
)

// InitWithUserJWT initializes a session using a JWT token
func (sdk *pogrSDK) InitWithUserJWT(userJWT string) (string, error) {
	headers := map[string]string{
		"POGR_CLIENT":   sdk.config.ClientKey,
		"POGR_BUILD":    sdk.config.BuildKey,
		"Authorization": fmt.Sprintf("Bearer %s", userJWT),
		"Content-Type":  "application/json",
	}

	req := &Request{
		Method:  "POST",
		URL:     fmt.Sprintf("%s/init", sdk.config.BaseURL),
		Headers: headers,
	}

	return sdk.handleInitResponse(req)
}

// InitWithAssociationID initializes a session using an association ID
func (sdk *pogrSDK) InitWithAssociationID(associationID string) (string, error) {
	data := map[string]string{"association_id": associationID}
	payload, err := json.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("failed to marshal data: %w", err)
	}

	headers := map[string]string{
		"POGR_CLIENT":  sdk.config.ClientKey,
		"POGR_BUILD":   sdk.config.BuildKey,
		"Content-Type": "application/json",
	}

	req := &Request{
		Method:  "POST",
		URL:     fmt.Sprintf("%s/init", sdk.config.BaseURL),
		Headers: headers,
		Body:    payload,
	}

	return sdk.handleInitResponse(req)
}

// InitWithSteamTicket initializes a session using a Steam ticket
func (sdk *pogrSDK) InitWithSteamTicket(steamTicket string) (string, error) {
	headers := map[string]string{
		"POGR_CLIENT": sdk.config.ClientKey,
		"POGR_BUILD":  sdk.config.BuildKey,
	}

	req := &Request{
		Method:  "POST",
		URL:     fmt.Sprintf("%s/init?steam_ticket=%s", sdk.config.BaseURL, steamTicket),
		Headers: headers,
	}

	return sdk.handleInitResponse(req)
}

// SendData sends data with optional tags using available authentication method
func (sdk *pogrSDK) SendData(data interface{}, tags *Tags) (string, error) {
    payload := DataPayload{
        Data: data,
        Tags: tags,
    }

    jsonData, err := json.Marshal(payload)
    if err != nil {
        return "", fmt.Errorf("failed to marshal data: %w", err)
    }

    headers, err := sdk.getAuthHeaders()
    if err != nil {
        return "", err
    }
    headers["Content-Type"] = "application/json"

    ctx := context.Background()
    if sdk.config.Timeout > 0 {
        var cancel context.CancelFunc
        ctx, cancel = context.WithTimeout(ctx, sdk.config.Timeout)
        defer cancel()
    }

    req := &Request{
        Method:  "POST",
        URL:     fmt.Sprintf("%s/data", sdk.config.BaseURL),
        Headers: headers,
        Body:    jsonData,
        Context: ctx,
    }

    return sdk.handleDataResponse(req)
}

// EndSession safely ends the current session
func (sdk *pogrSDK) EndSession() error {
	sdk.mu.Lock()
	defer sdk.mu.Unlock()

	if !sdk.state.initialized || sdk.state.sessionID == "" {
		return ErrNoActiveSession
	}

	headers := map[string]string{
		"INTAKE_SESSION_ID": sdk.state.sessionID,
	}

	req := &Request{
		Method:  "POST",
		URL:     fmt.Sprintf("%s/end", sdk.config.BaseURL),
		Headers: headers,
	}

	if err := sdk.handleGenericResponse(req); err != nil {
		return err
	}

	sdk.state.sessionID = ""
	sdk.state.initialized = false
	return nil
}

// SendEvent sends an event with relevant details and optional user tags
func (sdk *pogrSDK) SendEvent(event string, subEvent string, eventType string, eventFlag string, eventKey string, eventData map[string]interface{}, tags *Tags) (string, error) {
	eventPayload := map[string]interface{}{
		"event":       event,
		"sub_event":   subEvent,
		"event_type":  eventType,
		"event_flag":  eventFlag,
		"event_key":   eventKey,
		"event_data":  eventData,
		"tags":        tags,
	}

	jsonData, err := json.Marshal(eventPayload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal event data: %w", err)
	}

	headers, err := sdk.getAuthHeaders()
	if err != nil {
		return "", err
	}
	headers["Content-Type"] = "application/json"

	req := &Request{
		Method:  "POST",
		URL:     fmt.Sprintf("%s/event", sdk.config.BaseURL),
		Headers: headers,
		Body:    jsonData,
	}

	return sdk.handleDataResponse(req)
}

// SendLog submits a log entry for monitoring and auditing purposes
func (sdk *pogrSDK) SendLog(service string, environment string, severity string, logType string, logMessage string, data map[string]interface{}, tags *Tags) (string, error) {
	logPayload := map[string]interface{}{
		"service":     service,
		"environment": environment,
		"severity":    severity,
		"type":        logType,
		"log":         logMessage,
		"data":        data,
		"tags":        tags,
	}

	jsonData, err := json.Marshal(logPayload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal log data: %w", err)
	}

	headers, err := sdk.getAuthHeaders()
	if err != nil {
		return "", err
	}
	headers["Content-Type"] = "application/json"

	req := &Request{
		Method:  "POST",
		URL:     fmt.Sprintf("%s/logs", sdk.config.BaseURL),
		Headers: headers,
		Body:    jsonData,
	}

	return sdk.handleDataResponse(req)
}

// SendMetrics sends real-time metrics for monitoring purposes
func (sdk *pogrSDK) SendMetrics(service string, environment string, metrics map[string]interface{}, tags *Tags) (string, error) {
	metricsPayload := map[string]interface{}{
		"service":     service,
		"environment": environment,
		"metrics":     metrics,
		"tags":        tags,
	}

	jsonData, err := json.Marshal(metricsPayload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal metrics data: %w", err)
	}

	headers, err := sdk.getAuthHeaders()
	if err != nil {
		return "", err
	}
	headers["Content-Type"] = "application/json"

	req := &Request{
		Method:  "POST",
		URL:     fmt.Sprintf("%s/metrics", sdk.config.BaseURL),
		Headers: headers,
		Body:    jsonData,
	}

	return sdk.handleDataResponse(req)
}

// SendMonitorData sends system resource usage data
func (sdk *pogrSDK) SendMonitorData(cpuUsage float64, memoryUsage int, dllsLoaded []string, settings map[string]interface{}) (string, error) {
	monitorPayload := map[string]interface{}{
		"cpu_usage":    cpuUsage,
		"memory_usage": memoryUsage,
		"dlls_loaded":  dllsLoaded,
		"settings":     settings,
	}

	jsonData, err := json.Marshal(monitorPayload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal monitor data: %w", err)
	}

	headers, err := sdk.getAuthHeaders()
	if err != nil {
		return "", err
	}
	headers["Content-Type"] = "application/json"

	req := &Request{
		Method:  "POST",
		URL:     fmt.Sprintf("%s/monitor", sdk.config.BaseURL),
		Headers: headers,
		Body:    jsonData,
	}

	return sdk.handleDataResponse(req)
}



// IsInitialized returns the initialization status
func (sdk *pogrSDK) IsInitialized() bool {
	sdk.mu.RLock()
	defer sdk.mu.RUnlock()
	return sdk.state.initialized
}

// GetSessionID returns the current session ID
func (sdk *pogrSDK) GetSessionID() string {
	sdk.mu.RLock()
	defer sdk.mu.RUnlock()
	return sdk.state.sessionID
}

// ValidateTag checks if a tag key is valid
func (sdk *pogrSDK) ValidateTag(key string) bool {
	validTags := map[string]bool{
		"steam_id":           true,
		"twitch_id":          true,
		"association_id":     true,
		"pogr_game_session":  true,
		"xbox_id":            true,
		"battlenet_id":       true,
		"twitter_id":         true,
		"linkedin_id":        true,
		"pogr_player_id":     true,
		"discord_id":         true,
		"override_timestamp": true,
	}
	return validTags[key]
}

func (sdk *pogrSDK) getAuthHeaders() (map[string]string, error) {
	headers := make(map[string]string)

	// Check auth methods in priority order
	if sessionID := sdk.getSessionID(); sessionID != "" {
		headers["INTAKE_SESSION_ID"] = sessionID
		return headers, nil
	}

	if sdk.hasAccessKeyAuth() {
		headers["ACCESS_KEY"] = sdk.config.AccessKey
		headers["SECRET_KEY"] = sdk.config.SecretKey
		return headers, nil
	}

	if sdk.hasClientKeyAuth() {
		headers["POGR_CLIENT"] = sdk.config.ClientKey
		headers["POGR_BUILD"] = sdk.config.BuildKey
		return headers, nil
	}

	return nil, errors.New("no valid authentication method available")
}

func (sdk *pogrSDK) getSessionID() string {
	sdk.mu.RLock()
	defer sdk.mu.RUnlock()
	return sdk.state.sessionID
}

func (sdk *pogrSDK) hasAccessKeyAuth() bool {
	return sdk.config.AccessKey != "" && sdk.config.SecretKey != ""
}

func (sdk *pogrSDK) hasClientKeyAuth() bool {
	return sdk.config.ClientKey != "" && sdk.config.BuildKey != ""
}
