package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/nom-nom-hub/floss/internal/config"
	"github.com/spf13/cobra"
)

const (
	qwenDeviceCodeURL = "https://chat.qwen.ai/api/v1/oauth2/device/code"
	qwenTokenURL      = "https://chat.qwen.ai/api/v1/oauth2/token"
	qwenClientID      = "f0304373b74a44d2b584a3fb70ca9e56"
	qwenScope         = "openid profile email model.completion"
	qwenGrantType     = "urn:ietf:params:oauth:grant-type:device_code"
)

type deviceCodeResponse struct {
	DeviceCode              string `json:"device_code"`
	UserCode                string `json:"user_code"`
	VerificationURI         string `json:"verification_uri"`
	VerificationURIComplete string `json:"verification_uri_complete"`
	ExpiresIn               int    `json:"expires_in"`
	Interval                int    `json:"interval"`
}

type tokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
}

func init() {
	rootCmd.AddCommand(authCmd)
}

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Authentication commands",
	Long:  `Authentication commands for various providers`,
}

var qwenAuthCmd = &cobra.Command{
	Use:   "qwen",
	Short: "Authenticate with Qwen",
	Long:  `Authenticate with Qwen using OAuth2 device code flow`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return authenticateWithQwen(cmd.Context())
	},
}

func init() {
	authCmd.AddCommand(qwenAuthCmd)
}

func authenticateWithQwen(ctx context.Context) error {
	fmt.Println("Starting Qwen authentication process...")
	
	// Step 1: Request device code from Qwen
	fmt.Println("Requesting device code from Qwen...")
	deviceCodeResp, codeVerifier, err := requestDeviceCode()
	if err != nil {
		return fmt.Errorf("failed to request device code: %w", err)
	}
	fmt.Println("Device code received successfully!")

	// Step 2: Display login instructions to user
	fmt.Printf("Please visit %s to authenticate\n", deviceCodeResp.VerificationURIComplete)
	fmt.Printf("Enter code: %s\n", deviceCodeResp.UserCode)
	fmt.Println("Waiting for authentication...")
	
	// Force flush output to ensure it's displayed
	// This might help with terminal display issues
	fmt.Print("")

	// Step 3: Poll token endpoint until authorized
	fmt.Println("Starting token polling...")
	accessToken, refreshToken, err := pollForToken(deviceCodeResp, codeVerifier)
	if err != nil {
		return fmt.Errorf("failed to get token: %w", err)
	}
	fmt.Println("Token received successfully!")

	// Step 4: Save tokens to credentials store
	fmt.Println("Saving tokens to credentials store...")
	if err := saveTokens(accessToken, refreshToken); err != nil {
		return fmt.Errorf("failed to save tokens: %w", err)
	}
	fmt.Println("Tokens saved successfully!")

	fmt.Println("Successfully authenticated with Qwen!")
	return nil
}

func requestDeviceCode() (*deviceCodeResponse, string, error) {
	// Generate code verifier and challenge
	codeVerifier, err := generateCodeVerifier()
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate code verifier: %w", err)
	}
	codeChallenge := generateCodeChallenge(codeVerifier)

	// Prepare form data
	data := url.Values{}
	data.Set("client_id", qwenClientID)
	data.Set("scope", qwenScope)
	data.Set("code_challenge", codeChallenge)
	data.Set("code_challenge_method", "S256") // PKCE method

	// Make request
	resp, err := http.PostForm(qwenDeviceCodeURL, data)
	if err != nil {
		return nil, "", fmt.Errorf("failed to request device code: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, "", fmt.Errorf("device code request failed with status %d: %s", resp.StatusCode, string(body))
	}

	// Parse response
	var deviceCodeResp deviceCodeResponse
	if err := json.Unmarshal(body, &deviceCodeResp); err != nil {
		return nil, "", fmt.Errorf("failed to parse device code response: %w", err)
	}

	return &deviceCodeResp, codeVerifier, nil
}

func pollForToken(deviceCodeResp *deviceCodeResponse, codeVerifier string) (string, string, error) {
	interval := deviceCodeResp.Interval
	if interval == 0 {
		interval = 5 // Default to 5 seconds if not specified
	}

	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	defer ticker.Stop()

	timeout := time.After(time.Duration(deviceCodeResp.ExpiresIn) * time.Second)

	for {
		select {
		case <-timeout:
			return "", "", fmt.Errorf("authentication timed out")
		case <-ticker.C:
			accessToken, refreshToken, err := tryGetToken(deviceCodeResp.DeviceCode, codeVerifier)
			if err != nil {
				// Check if it's an expected error (authorization pending)
				if strings.Contains(err.Error(), "authorization_pending") {
					continue // Keep polling
				}
				return "", "", err // Unexpected error
			}
			return accessToken, refreshToken, nil
		}
	}
}

func tryGetToken(deviceCode, codeVerifier string) (string, string, error) {
	// Prepare form data
	data := url.Values{}
	data.Set("client_id", qwenClientID)
	data.Set("device_code", deviceCode)
	data.Set("grant_type", qwenGrantType)
	data.Set("code_verifier", codeVerifier)

	// Make request
	resp, err := http.PostForm(qwenTokenURL, data)
	if err != nil {
		return "", "", fmt.Errorf("failed to request token: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", fmt.Errorf("failed to read response: %w", err)
	}

	// Check if this is an error response
	var errorResp struct {
		Error string `json:"error"`
	}
	if err := json.Unmarshal(body, &errorResp); err == nil && errorResp.Error != "" {
		return "", "", fmt.Errorf("%s", errorResp.Error)
	}

	// Parse successful response
	var tokenResp tokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return "", "", fmt.Errorf("failed to parse token response: %w", err)
	}

	return tokenResp.AccessToken, tokenResp.RefreshToken, nil
}

func saveTokens(accessToken, refreshToken string) error {
	// Get config directory
	cfg := config.Get()
	if cfg == nil {
		// If config is nil, try to initialize it
		var err error
		cfg, err = config.Load("", "", false)
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}
		if cfg == nil {
			return fmt.Errorf("config is still nil after loading")
		}
	}
	
	dataDir := cfg.Options.DataDirectory
	if dataDir == "" {
		// Use default data directory if not set
		dataDir = filepath.Join(os.Getenv("LOCALAPPDATA"), "floss")
		if dataDir == "" {
			// Fallback for non-Windows systems
			homeDir, err := os.UserHomeDir()
			if err != nil {
				return fmt.Errorf("failed to get home directory: %w", err)
			}
			dataDir = filepath.Join(homeDir, ".floss")
		}
	}

	// Create credentials file path
	credentialsPath := filepath.Join(dataDir, "credentials.json")
	
	// Ensure the directory exists
	if err := os.MkdirAll(dataDir, 0700); err != nil {
		return fmt.Errorf("failed to create data directory: %w", err)
	}

	// Read existing credentials
	credentials := make(map[string]map[string]string)
	if data, err := os.ReadFile(credentialsPath); err == nil {
		if err := json.Unmarshal(data, &credentials); err != nil {
			// If we can't parse existing credentials, start fresh
			credentials = make(map[string]map[string]string)
		}
	}

	// Update Qwen credentials
	if credentials["qwen"] == nil {
		credentials["qwen"] = make(map[string]string)
	}
	credentials["qwen"]["access_token"] = accessToken
	credentials["qwen"]["refresh_token"] = refreshToken

	// Write credentials back to file
	data, err := json.MarshalIndent(credentials, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal credentials: %w", err)
	}

	if err := os.WriteFile(credentialsPath, data, 0600); err != nil {
		return fmt.Errorf("failed to write credentials file: %w", err)
	}

	// Also save the API key to the config
	// Note: We might not be able to save to the main config file in this context
	// but the credentials file is sufficient for authentication
	
	return nil
}