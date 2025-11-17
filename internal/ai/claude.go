package ai

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/anthropic-ai/anthropic-sdk-go"
	"github.com/anthropic-ai/anthropic-sdk-go/option"
)

// ClaudeClient represents the Claude AI client
type ClaudeClient struct {
	client *anthropic.Client
}

// NewClaudeClient creates a new Claude API client
func NewClaudeClient(apiKey string) *ClaudeClient {
	client := anthropic.NewClient(
		option.WithAPIKey(apiKey),
	)

	return &ClaudeClient{
		client: client,
	}
}

// QueryFleet sends a natural language query about the fleet
func (c *ClaudeClient) QueryFleet(userQuery string, fleetContext FleetContext) (*QueryResponse, error) {
	systemPrompt := `You are a fleet dispatch AI assistant for Radius Recycling. Your role is to:

1. Answer questions about vehicle, driver, and trailer locations
2. Help find specific assets or vehicles
3. Provide route information and distance calculations
4. Interpret natural language queries about fleet operations

CRITICAL RULES:
- Always use Driver ID numbers (never names) for privacy
- Provide clear, concise responses
- Include relevant location details when available
- Format coordinates as "latitude, longitude"
- Convert distances to miles and durations to hours/minutes

DATA YOU HAVE ACCESS TO:
- Real-time vehicle locations and status
- Driver locations and IDs
- Asset (trailer/container) locations
- Facility locations and geofences

OUTPUT FORMAT:
Provide clear, structured responses with:
- Direct answer to the query
- Relevant details (locations, distances, times)
- Any important warnings or notes`

	// Build context message
	contextJSON, err := json.MarshalIndent(fleetContext, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal fleet context: %w", err)
	}

	userMessage := fmt.Sprintf(`Query: %s

Current Fleet Data:
%s

Please analyze this query and provide a clear, helpful response based on the current fleet data.`, userQuery, string(contextJSON))

	// Call Claude API
	message, err := c.client.Messages.New(context.Background(), anthropic.MessageNewParams{
		Model:     anthropic.F(anthropic.ModelClaude_3_5_Sonnet_20241022),
		MaxTokens: anthropic.Int(2000),
		System: anthropic.F([]anthropic.TextBlockParam{
			anthropic.NewTextBlock(systemPrompt),
		}),
		Messages: anthropic.F([]anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock(userMessage)),
		}),
	})

	if err != nil {
		return nil, fmt.Errorf("Claude API call failed: %w", err)
	}

	// Extract response text
	var responseText string
	if len(message.Content) > 0 {
		if textBlock, ok := message.Content[0].AsUnion().(anthropic.TextBlock); ok {
			responseText = textBlock.Text
		}
	}

	return &QueryResponse{
		Response:   responseText,
		Model:      string(message.Model),
		StopReason: string(message.StopReason),
		InputTokens:  int(message.Usage.InputTokens),
		OutputTokens: int(message.Usage.OutputTokens),
	}, nil
}

// FleetContext represents the current state of the fleet
type FleetContext struct {
	QueryTime       string                 `json:"query_time"`
	VehicleCount    int                    `json:"vehicle_count"`
	DriverCount     int                    `json:"driver_count"`
	AssetCount      int                    `json:"asset_count"`
	Vehicles        []VehicleSummary       `json:"vehicles,omitempty"`
	Drivers         []DriverSummary        `json:"drivers,omitempty"`
	Assets          []AssetSummary         `json:"assets,omitempty"`
	Facilities      []FacilitySummary      `json:"facilities,omitempty"`
}

// VehicleSummary represents a vehicle for AI context
type VehicleSummary struct {
	VehicleNumber string   `json:"vehicle_number"`
	Status        string   `json:"status"`
	Location      string   `json:"location"`
	Latitude      *float64 `json:"latitude,omitempty"`
	Longitude     *float64 `json:"longitude,omitempty"`
	SpeedMPH      *int     `json:"speed_mph,omitempty"`
	FuelGallons   *float64 `json:"fuel_gallons,omitempty"`
	DriverID      *int     `json:"driver_id,omitempty"`
}

// DriverSummary represents a driver for AI context
type DriverSummary struct {
	DriverID         int      `json:"driver_id"`
	Status           string   `json:"status"`
	Location         string   `json:"location"`
	Latitude         *float64 `json:"latitude,omitempty"`
	Longitude        *float64 `json:"longitude,omitempty"`
	HOSRemainingMins *int     `json:"hos_remaining_minutes,omitempty"`
	VehicleNumber    string   `json:"vehicle_number,omitempty"`
}

// AssetSummary represents an asset for AI context
type AssetSummary struct {
	AssetName   string   `json:"asset_name"`
	AssetType   string   `json:"asset_type"`
	Status      string   `json:"status"`
	Location    string   `json:"location"`
	Latitude    *float64 `json:"latitude,omitempty"`
	Longitude   *float64 `json:"longitude,omitempty"`
	AttachedTo  string   `json:"attached_to,omitempty"`
}

// FacilitySummary represents a facility for AI context
type FacilitySummary struct {
	Name      string  `json:"name"`
	Type      string  `json:"type"`
	Address   string  `json:"address"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// QueryResponse represents Claude's response
type QueryResponse struct {
	Response     string `json:"response"`
	Model        string `json:"model"`
	StopReason   string `json:"stop_reason"`
	InputTokens  int    `json:"input_tokens"`
	OutputTokens int    `json:"output_tokens"`
}
