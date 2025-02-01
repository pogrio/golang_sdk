package pogr

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

// Common errors
var (
	ErrNoActiveSession = errors.New("no active session")
	ErrInvalidData     = errors.New("invalid data provided")
	ErrUnauthorized    = errors.New("unauthorized request")
)

// pogrSDK implements the POGRService interface with thread-safety
type pogrSDK struct {
	config     Config
	httpClient HTTPClient
	mu         sync.RWMutex // Protects session state
	state      sessionState
}

// sessionState encapsulates mutable session data
type sessionState struct {
	sessionID   string
	initialized bool
}

// NewPOGRSDK creates a new thread-safe instance of the POGR SDK
func NewPOGRSDK(config Config) POGRService {
	if config.BaseURL == "" {
		config.BaseURL = "https://api.pogr.io/v1/intake"
	}

	if config.HTTPClient == nil {
		config.HTTPClient = NewDefaultHTTPClient(config)
	}

	return &pogrSDK{
		config:     config,
		httpClient: config.HTTPClient,
	}
}

// NewDefaultHTTPClient creates a default HTTP client with optional connection pooling
func NewDefaultHTTPClient(config Config) HTTPClient {
	var transport *http.Transport

	if config.EnableConnectionPool {
		poolConfig := config.PoolConfig
		if poolConfig == nil {
			poolConfig = DefaultPoolConfig()
		}
		transport = &http.Transport{
			MaxIdleConns:        poolConfig.MaxIdleConns,
			MaxIdleConnsPerHost: poolConfig.MaxIdleConnsPerHost,
			MaxConnsPerHost:     poolConfig.MaxConnsPerHost,
			IdleConnTimeout:     poolConfig.IdleConnTimeout,
		}
	} else {
		transport = &http.Transport{}
	}

	return &defaultHTTPClient{
		client: &http.Client{
			Transport: transport,
			Timeout:   config.Timeout,
		},
	}
}

// DefaultPoolConfig returns default connection pool settings
func DefaultPoolConfig() *ConnectionPoolConfig {
	return &ConnectionPoolConfig{
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 100,
		MaxConnsPerHost:     100,
		IdleConnTimeout:     90 * time.Second,
	}
}

// setSessionState safely updates session state
func (sdk *pogrSDK) setSessionState(sessionID string, initialized bool) {
	sdk.mu.Lock()
	defer sdk.mu.Unlock()
	sdk.state.sessionID = sessionID
	sdk.state.initialized = initialized
}

// PrintConfig returns a string representation of the current configuration
func (sdk *pogrSDK) PrintConfig() string {
	return fmt.Sprintf(`
POGR SDK Configuration:
BaseURL: %s
ClientKey: %s
BuildKey: %s
AccessKey: %s
SecretKey: %s
Connection Pool Enabled: %v
Timeout: %v`,
		sdk.config.BaseURL,
		sdk.config.ClientKey,
		sdk.config.BuildKey,
		sdk.config.AccessKey,
		sdk.config.SecretKey,
		sdk.config.EnableConnectionPool,
		sdk.config.Timeout)
}

// defaultHTTPClient implements the HTTPClient interface
type defaultHTTPClient struct {
	client *http.Client
}

func (c *defaultHTTPClient) Do(req *Request) (*Response, error) {
	var httpReq *http.Request
	var err error

	if req.Context != nil {
		httpReq, err = http.NewRequestWithContext(req.Context, req.Method, req.URL, bytes.NewBuffer(req.Body))
	} else {
		httpReq, err = http.NewRequest(req.Method, req.URL, bytes.NewBuffer(req.Body))
	}

	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	for key, value := range req.Headers {
		httpReq.Header.Set(key, value)
	}

	resp, err := c.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	headers := make(map[string]string)
	for key, values := range resp.Header {
		if len(values) > 0 {
			headers[key] = values[0]
		}
	}

	return &Response{
		StatusCode: resp.StatusCode,
		Body:       body,
		Headers:    headers,
	}, nil
}
