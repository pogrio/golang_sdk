// Package main demonstrates the usage of the POGR SDK
package main

import (
	"log"
	"os"
	"time"

	"github.com/pogrio/golang_sdk/pogr"

	"github.com/joho/godotenv"
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

	// Session based intake
	config1 := pogr.Config{
		ClientKey:            clientID,
		BuildKey:             buildID,
		BaseURL:              intakeBaseURL,
		Timeout:              30 * time.Second,
		EnableConnectionPool: true,
	}

	sdkWithSession := pogr.NewPOGRSDK(config1)
	log.Printf("SDK Configuration:\n%s", sdkWithSession.PrintConfig())

	// Run examples with different authentication methods
	runJWTExample(sdkWithSession)
	runAssociationIDExample(sdkWithSession)
	runSteamTicketExample(sdkWithSession)

	// AccessKey based intake
	config2 := pogr.Config{
		AccessKey:            accessKey,
		SecretKey:            secretKey,
		BaseURL:              intakeBaseURL,
		Timeout:              30 * time.Second,
		EnableConnectionPool: true,
	}

	sdkWithAccessKey := pogr.NewPOGRSDK(config2)
	log.Printf("SDK Configuration:\n%s", sdkWithAccessKey.PrintConfig())

	sendTestDataWithoutSession(sdkWithAccessKey, "accessKey")

	// ClientID without session
	sdkWithClientID := pogr.NewPOGRSDK(config1)
	log.Printf("SDK Configuration:\n%s", sdkWithClientID.PrintConfig())

	sendTestDataWithoutSession(sdkWithAccessKey, "clientID")

}

// runJWTExample demonstrates JWT authentication
func runJWTExample(sdk pogr.POGRService) {
	sessionID, err := sdk.InitWithUserJWT(pogrJWT)
	if err != nil {
		log.Printf("JWT initialization failed: %v", err)
		return
	}
	log.Printf("JWT Session initialized: %s", sessionID)

	// Send test data
	sendTestData(sdk, "JWT")
}

// runAssociationIDExample demonstrates Association ID authentication
func runAssociationIDExample(sdk pogr.POGRService) {
	sessionID, err := sdk.InitWithAssociationID(associationID)
	if err != nil {
		log.Printf("Association ID initialization failed: %v", err)
		return
	}
	log.Printf("Association ID Session initialized: %s", sessionID)

	// Send test data
	sendTestData(sdk, "AssociationID")
}

// runSteamTicketExample demonstrates Steam Ticket authentication
func runSteamTicketExample(sdk pogr.POGRService) {
	sessionID, err := sdk.InitWithSteamTicket(steamAuthTicket)
	if err != nil {
		log.Printf("Steam Ticket initialization failed: %v", err)
		return
	}
	log.Printf("Steam Ticket Session initialized: %s", sessionID)

	// Send test data
	sendTestData(sdk, "SteamTicket")
}

// sendTestData sends sample data with the specified authentication method
func sendTestData(sdk pogr.POGRService, authMethod string) {

	data := map[string]interface{}{
		"auth_method": authMethod,
		"timestamp":   time.Now().Unix(),
		"test_data":   "Hello POGR!",
	}

	dataID1, err := sdk.SendData(data, nil)
	if err != nil {
		log.Printf("Failed to send data with %s auth: %v", authMethod, err)
		return
	}
	log.Printf("Data sent with %s auth: %s", authMethod, dataID1)

	tags := &pogr.Tags{
		OverrideTimestamp: time.Now().Add(-30 * time.Minute).Format(time.RFC3339),
	}

	dataID2, err := sdk.SendData(data, tags)
	if err != nil {
		log.Printf("Failed to send data with %s auth: %v", authMethod, err)
		return
	}
	log.Printf("Data sent with %s auth: %s", authMethod, dataID2)

	complexData := map[string]interface{}{
		"data":           "string",
		"any":            1,
		"type":           true,
		"doesn_t_matter": []int{1, 2, 3, 4, 5, 6},
		"any_thing": map[string]interface{}{
			"we":     "string",
			"accept": "string",
			"it":     true,
		},
	}

	dataID3, err := sdk.SendData(complexData, tags)
	if err != nil {
		log.Printf("Failed to send data with %s auth: %v", authMethod, err)
		return
	}
	log.Printf("Data sent with %s auth: %s", authMethod, dataID3)

	if err := sdk.EndSession(); err != nil {
		log.Printf("Failed to end %s session: %v", authMethod, err)
	}

}

// sendTestDataWithoutSession sends data without session
func sendTestDataWithoutSession(sdk pogr.POGRService, intakeType string) {
	data := map[string]interface{}{
		"auth_method": intakeType,
		"timestamp":   time.Now().Unix(),
		"test_data":   "Hello POGR!",
	}

	dataID1, err := sdk.SendData(data, nil)
	if err != nil {
		log.Printf("Failed to send data with %s: %v", intakeType, err)
		return
	}
	log.Printf("Data sent with %s auth: %s", intakeType, dataID1)

	tags := &pogr.Tags{
		TwitchID:          twitchID,
		OverrideTimestamp: time.Now().Add(-30 * time.Minute).Format(time.RFC3339),
	}

	dataID2, err := sdk.SendData(data, tags)
	if err != nil {
		log.Printf("Failed to send data with %s auth: %v", intakeType, err)
		return
	}
	log.Printf("Data sent with %s auth: %s", intakeType, dataID2)

	complexData := map[string]interface{}{
		"data":           "string",
		"any":            1,
		"type":           true,
		"doesn_t_matter": []int{1, 2, 3, 4, 5, 6},
		"any_thing": map[string]interface{}{
			"we":     "string",
			"accept": "string",
			"it":     true,
		},
	}

	dataID3, err := sdk.SendData(complexData, tags)
	if err != nil {
		log.Printf("Failed to send data with %s auth: %v", intakeType, err)
		return
	}
	log.Printf("Data sent with %s auth: %s", intakeType, dataID3)
}

// demonstrateUtilities shows usage of utility methods
func demonstrateUtilities(sdk pogr.POGRService) {
	log.Printf("SDK Initialized: %v", sdk.IsInitialized())
	log.Printf("Session ID: %s", sdk.GetSessionID())

	tags := []string{"steam_id", "twitch_id", "discord_id", "invalid_tag"}
	for _, tag := range tags {
		log.Printf("Tag '%s' valid: %v", tag, sdk.ValidateTag(tag))
	}
}
