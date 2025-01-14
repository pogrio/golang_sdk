package main

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/pogrio/golang_sdk/pogr"
)

var (
	clientID        string
	buildID         string
	accessKey       string
	secretKey       string
	intakeBaseURL   string
	pogrJWT         string
	twitchID        string
	associationID   string
	steamAuthTicket string
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found: %v", err)
	}

	clientID = os.Getenv("POGR_CLIENT_ID")
	buildID = os.Getenv("POGR_BUILD_ID")
	accessKey = os.Getenv("POGR_ACCESS_KEY")
	secretKey = os.Getenv("POGR_SECRET_KEY")
	intakeBaseURL = os.Getenv("INTAKE_BASE_URL")
	pogrJWT = os.Getenv("POGR_JWT")
	twitchID = os.Getenv("TWITCH_ID")
	associationID = os.Getenv("POGR_ASSOCIATION_ID")
	steamAuthTicket = os.Getenv("STEAM_AUTHENTICATION_TICKET")
}

func main() {
	config := pogr.Config{
		ClientKey:            clientID,
		BuildKey:             buildID,
		AccessKey:            accessKey,
		SecretKey:            secretKey,
		BaseURL:              intakeBaseURL,
		Timeout:              30 * time.Second,
		EnableConnectionPool: true,
	}

	sdk := pogr.NewPOGRSDK(config)

	log.Printf("SDK Configuration:\n%s", sdk.PrintConfig())

	// Authentication examples
	runJWTExample(sdk)
	runAssociationIDExample(sdk)
	runSteamTicketExample(sdk)

	// Data example
	sendTestData(sdk, "Session-based auth")

	// Example usage of additional endpoints
	runEventExample(sdk)
	runLogExample(sdk)
	runMetricsExample(sdk)
	runMonitorExample(sdk)
}

func runJWTExample(sdk pogr.POGRService) {
	sessionID, err := sdk.InitWithUserJWT(pogrJWT)
	if err != nil {
		log.Printf("JWT initialization failed: %v", err)
		return
	}
	log.Printf("JWT Session initialized: %s", sessionID)
}

func runAssociationIDExample(sdk pogr.POGRService) {
	sessionID, err := sdk.InitWithAssociationID(associationID)
	if err != nil {
		log.Printf("Association ID initialization failed: %v", err)
		return
	}
	log.Printf("Association ID Session initialized: %s", sessionID)
}

func runSteamTicketExample(sdk pogr.POGRService) {
	sessionID, err := sdk.InitWithSteamTicket(steamAuthTicket)
	if err != nil {
		log.Printf("Steam Ticket initialization failed: %v", err)
		return
	}
	log.Printf("Steam Ticket Session initialized: %s", sessionID)
}

func sendTestData(sdk pogr.POGRService, authMethod string) {
	data := map[string]interface{}{
		"auth_method": authMethod,
		"timestamp":   time.Now().Unix(),
		"test_data":   "Hello POGR!",
	}

	dataID, err := sdk.SendData(data, nil)
	if err != nil {
		log.Printf("Failed to send data with %s: %v", authMethod, err)
		return
	}
	log.Printf("Data sent successfully: %s", dataID)
}

func runEventExample(sdk pogr.POGRService) {
	eventData := map[string]interface{}{
		"player_id":        "12345",
		"achievement_name": "Master Explorer",
	}

	tags := &pogr.Tags{
		DiscordID: "9480bc67-88e6-42ee-bfb6-0c70137d1fad",
	}

	eventID, err := sdk.SendEvent("player_login", "level_up", "achievement", "completed", "level_5_unlocked", eventData, tags)
	if err != nil {
		log.Printf("Failed to send event: %v", err)
		return
	}
	log.Printf("Event sent successfully: %s", eventID)
}

func runLogExample(sdk pogr.POGRService) {
	logData := map[string]interface{}{
		"user_id":    "user_56789",
		"timestamp":  "2025-01-01T12:00:00Z",
		"ip_address": "192.168.1.1",
	}

	tags := &pogr.Tags{
		TwitchID:          "f88e0a15-26fa-492a-b83b-9861a44522df",
		OverrideTimestamp: time.Now().Add(-1 * time.Hour).Format(time.RFC3339),
	}

	logID, err := sdk.SendLog("authentication", "live", "info", "user-login", "User logged in successfully", logData, tags)
	if err != nil {
		log.Printf("Failed to send log: %v", err)
		return
	}
	log.Printf("Log sent successfully: %s", logID)
}

func runMetricsExample(sdk pogr.POGRService) {
	metrics := map[string]interface{}{
		"players_online":         250,
		"average_latency_ms":     35.7,
		"server_load_percentage": 80.5,
	}

	tags := &pogr.Tags{
		SteamID: "steam_001234",
	}

	metricsID, err := sdk.SendMetrics("game_server", "production", metrics, tags)
	if err != nil {
		log.Printf("Failed to send metrics: %v", err)
		return
	}
	log.Printf("Metrics sent successfully: %s", metricsID)
}

func runMonitorExample(sdk pogr.POGRService) {
	dlls := []string{"graphics.dll", "physics.dll", "audio.dll"}

	settings := map[string]interface{}{
		"resolution": "1920x1080",
		"fps_limit":  60,
	}

	monitorID, err := sdk.SendMonitorData(10.5, 4096000, dlls, settings)
	if err != nil {
		log.Printf("Failed to send monitor data: %v", err)
		return
	}
	log.Printf("Monitor data sent successfully: %s", monitorID)
}
