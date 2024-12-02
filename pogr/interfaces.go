package pogr

import (
	"context"
	"time"
)

// POGRService defines the interface for POGR SDK operations
type POGRService interface {
	// Session Management
	InitWithUserJWT(userJWT string) (string, error)
	InitWithAssociationID(associationID string) (string, error)
	InitWithSteamTicket(steamTicket string) (string, error)
	EndSession() error

	// Data Operations
	SendData(data interface{}, tags *Tags) (string, error)

	// Utility Methods
	IsInitialized() bool
	GetSessionID() string
	ValidateTag(key string) bool
	PrintConfig() string
}

// Config holds the configuration options for the SDK
type Config struct {
	ClientKey            string
	BuildKey             string
	AccessKey            string
	SecretKey            string
	BaseURL              string
	HTTPClient           HTTPClient
	Timeout              time.Duration
	EnableConnectionPool bool
	PoolConfig           *ConnectionPoolConfig
}

// ConnectionPoolConfig holds connection pool settings
type ConnectionPoolConfig struct {
	MaxIdleConns        int
	MaxIdleConnsPerHost int
	MaxConnsPerHost     int
	IdleConnTimeout     time.Duration
}

// Request represents an HTTP request with context support
type Request struct {
	Method  string
	URL     string
	Headers map[string]string
	Body    []byte
	Context context.Context
}

// Response represents an HTTP response
type Response struct {
	StatusCode int
	Body       []byte
	Headers    map[string]string
}

// HTTPClient interface for making HTTP requests
type HTTPClient interface {
	Do(req *Request) (*Response, error)
}

// Tags represents the available tag options for data
type Tags struct {
	DiscordID         string `json:"discord_id,omitempty"`
	SteamID           string `json:"steam_id,omitempty"`
	TwitchID          string `json:"twitch_id,omitempty"`
	AssociationID     string `json:"association_id,omitempty"`
	PogrGameSession   string `json:"pogr_game_session,omitempty"`
	XboxID            string `json:"xbox_id,omitempty"`
	BattlenetID       string `json:"battlenet_id,omitempty"`
	TwitterID         string `json:"twitter_id,omitempty"`
	LinkedinID        string `json:"linkedin_id,omitempty"`
	PogrPlayerID      string `json:"pogr_player_id,omitempty"`
	OverrideTimestamp string `json:"override_timestamp,omitempty"`
}

// DataPayload represents the structure for sending data with optional tags
type DataPayload struct {
	Data interface{} `json:"data"`
	Tags *Tags       `json:"tags,omitempty"`
}
