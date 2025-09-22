package provider

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/charmbracelet/catwalk/pkg/catwalk"
	"github.com/nom-nom-hub/floss/internal/config"
	"github.com/nom-nom-hub/floss/internal/llm/tools"
	"github.com/nom-nom-hub/floss/internal/log"
	"github.com/nom-nom-hub/floss/internal/message"
)

type qwenClient struct {
	providerOptions providerClientOptions
	client          *http.Client
}

type QwenClient ProviderClient

func newQwenClient(opts providerClientOptions) QwenClient {
	return &qwenClient{
		providerOptions: opts,
		client:          createQwenClient(opts),
	}
}

func createQwenClient(opts providerClientOptions) *http.Client {
	httpClient := &http.Client{}
	
	if config.Get().Options.Debug {
		httpClient = log.NewHTTPClient()
	}

	return httpClient
}

// Qwen API request/response structures
type QwenChatCompletionRequest struct {
	Model       string                   `json:"model"`
	Messages    []QwenMessage            `json:"messages"`
	Tools       []QwenTool               `json:"tools,omitempty"`
	MaxTokens   int                      `json:"max_tokens,omitempty"`
	Temperature float64                  `json:"temperature,omitempty"`
	Stream      bool                     `json:"stream,omitempty"`
}

type QwenMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type QwenTool struct {
	Type     string                 `json:"type"`
	Function QwenFunctionDefinition `json:"function"`
}

type QwenFunctionDefinition struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description,omitempty"`
	Parameters  map[string]interface{} `json:"parameters,omitempty"`
}

type QwenChatCompletionResponse struct {
	ID      string            `json:"id"`
	Object  string            `json:"object"`
	Created int64             `json:"created"`
	Model   string            `json:"model"`
	Choices []QwenChoice      `json:"choices"`
	Usage   QwenUsage         `json:"usage"`
}

type QwenChoice struct {
	Index        int              `json:"index"`
	Message      QwenMessage      `json:"message"`
	FinishReason string           `json:"finish_reason"`
}

type QwenUsage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type QwenStreamChunk struct {
	ID      string        `json:"id"`
	Object  string        `json:"object"`
	Created int64         `json:"created"`
	Model   string        `json:"model"`
	Choices []QwenChoice  `json:"choices"`
}

type tokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	ResourceURL  string `json:"resource_url,omitempty"`
}

func (q *qwenClient) convertMessages(messages []message.Message) []QwenMessage {
	var qwenMessages []QwenMessage

	// Add system message first
	systemMessage := q.providerOptions.systemMessage
	if q.providerOptions.systemPromptPrefix != "" {
		systemMessage = q.providerOptions.systemPromptPrefix + "\n" + systemMessage
	}

	if systemMessage != "" {
		qwenMessages = append(qwenMessages, QwenMessage{
			Role:    "system",
			Content: systemMessage,
		})
	}

	for _, msg := range messages {
		switch msg.Role {
		case message.User:
			qwenMessages = append(qwenMessages, QwenMessage{
				Role:    "user",
				Content: msg.Content().String(),
			})
		case message.Assistant:
			qwenMessages = append(qwenMessages, QwenMessage{
				Role:    "assistant",
				Content: msg.Content().String(),
			})
		case message.Tool:
			for _, result := range msg.ToolResults() {
				qwenMessages = append(qwenMessages, QwenMessage{
					Role:    "tool",
					Content: result.Content,
				})
			}
		}
	}

	return qwenMessages
}

func (q *qwenClient) convertTools(tools []tools.BaseTool) []QwenTool {
	var qwenTools []QwenTool

	for _, tool := range tools {
		info := tool.Info()
		qwenTools = append(qwenTools, QwenTool{
			Type: "function",
			Function: QwenFunctionDefinition{
				Name:        info.Name,
				Description: info.Description,
				Parameters: map[string]interface{}{
					"type":       "object",
					"properties": info.Parameters,
					"required":   info.Required,
				},
			},
		})
	}

	return qwenTools
}

func (q *qwenClient) finishReason(reason string) message.FinishReason {
	switch reason {
	case "stop":
		return message.FinishReasonEndTurn
	case "length":
		return message.FinishReasonMaxTokens
	case "tool_calls":
		return message.FinishReasonToolUse
	default:
		return message.FinishReasonUnknown
	}
}

func (q *qwenClient) preparedParams(messages []QwenMessage, tools []QwenTool) QwenChatCompletionRequest {
	model := q.providerOptions.model(q.providerOptions.modelType)
	cfg := config.Get()

	modelConfig := cfg.Models[config.SelectedModelTypeLarge]
	if q.providerOptions.modelType == config.SelectedModelTypeSmall {
		modelConfig = cfg.Models[config.SelectedModelTypeSmall]
	}

	maxTokens := model.DefaultMaxTokens
	if modelConfig.MaxTokens > 0 {
		maxTokens = modelConfig.MaxTokens
	}

	// Override max tokens if set in provider options
	if q.providerOptions.maxTokens > 0 {
		maxTokens = q.providerOptions.maxTokens
	}

	params := QwenChatCompletionRequest{
		Model:     model.ID,
		Messages:  messages,
		Tools:     tools,
		MaxTokens: int(maxTokens),
	}

	return params
}

func (q *qwenClient) refreshToken() error {
	// Load refresh token from credentials file
	credentialsPath := config.Get().Options.DataDirectory + "/credentials.json"
	
	// Read existing credentials
	credentials := make(map[string]map[string]string)
	if data, err := os.ReadFile(credentialsPath); err == nil {
		json.Unmarshal(data, &credentials)
	}

	qwenCreds, exists := credentials["qwen"]
	if !exists {
		return fmt.Errorf("no Qwen credentials found")
	}

	refreshToken, exists := qwenCreds["refresh_token"]
	if !exists {
		return fmt.Errorf("no refresh token found")
	}

	// Prepare form data for token refresh
	data := url.Values{}
	data.Set("client_id", "f0304373b74a44d2b584a3fb70ca9e56")
	data.Set("refresh_token", refreshToken)
	data.Set("grant_type", "refresh_token")

	// Make request
	resp, err := http.PostForm("https://chat.qwen.ai/api/v1/oauth2/token", data)
	if err != nil {
		return fmt.Errorf("failed to refresh token: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("token refresh failed with status %d: %s", resp.StatusCode, string(body))
	}

	// Parse response
	var tokenResp tokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return fmt.Errorf("failed to parse token response: %w", err)
	}

	// Update the API key in the provider options
	q.providerOptions.apiKey = tokenResp.AccessToken

	// Save the new tokens
	credentials["qwen"]["access_token"] = tokenResp.AccessToken
	credentials["qwen"]["refresh_token"] = tokenResp.RefreshToken
	
	// Save the resource URL if provided
	if tokenResp.ResourceURL != "" {
		credentials["qwen"]["resource_url"] = tokenResp.ResourceURL
	}

	// Write credentials back to file
	credData, err := json.MarshalIndent(credentials, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal credentials: %w", err)
	}

	if err := os.WriteFile(credentialsPath, credData, 0600); err != nil {
		return fmt.Errorf("failed to write credentials file: %w", err)
	}

	// Also update the API key in the config
	if err := config.Get().SetProviderAPIKey("qwen", "Bearer "+tokenResp.AccessToken); err != nil {
		return fmt.Errorf("failed to save API key to config: %w", err)
	}

	return nil
}

func (q *qwenClient) send(ctx context.Context, messages []message.Message, tools []tools.BaseTool) (response *ProviderResponse, err error) {
	params := q.preparedParams(q.convertMessages(messages), q.convertTools(tools))
	
	// Set up the request
	baseURL := q.providerOptions.baseURL
	if baseURL == "" {
		// Try to get the resource URL from credentials
		credentialsPath := config.Get().Options.DataDirectory + "/credentials.json"
		credentials := make(map[string]map[string]string)
		if data, err := os.ReadFile(credentialsPath); err == nil {
			json.Unmarshal(data, &credentials)
		}
		
		if qwenCreds, exists := credentials["qwen"]; exists {
			if resourceURL, exists := qwenCreds["resource_url"]; exists && resourceURL != "" {
				baseURL = resourceURL
			}
		}
		
		// Fallback to the default DashScope endpoint
		if baseURL == "" {
			baseURL = "https://dashscope.aliyuncs.com/compatible-mode/v1"
		}
	}
	
	url := baseURL + "/chat/completions"
	jsonData, err := json.Marshal(params)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, strings.NewReader(string(jsonData)))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+q.providerOptions.apiKey)
	
	// Add required DashScope headers
	req.Header.Set("X-DashScope-CacheControl", "enable")
	req.Header.Set("X-DashScope-UserAgent", "Floss/1.0 (Windows; x86_64)")
	req.Header.Set("X-DashScope-AuthType", "qwen-oauth")
	
	// Add extra headers
	for key, value := range q.providerOptions.extraHeaders {
		req.Header.Set(key, value)
	}

	// Make the request
	resp, err := q.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Check if we need to refresh token
	if resp.StatusCode == http.StatusUnauthorized {
		// Try to refresh token
		if refreshErr := q.refreshToken(); refreshErr != nil {
			return nil, fmt.Errorf("token refresh failed: %w", refreshErr)
		}

		// Retry the request with new token
		req.Header.Set("Authorization", "Bearer "+q.providerOptions.apiKey)
		resp, err = q.client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("failed to send request after token refresh: %w", err)
		}
		defer resp.Body.Close()
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var qwenResponse QwenChatCompletionResponse
	if err := json.Unmarshal(body, &qwenResponse); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if len(qwenResponse.Choices) == 0 {
		return nil, fmt.Errorf("received empty response from Qwen API")
	}

	content := qwenResponse.Choices[0].Message.Content
	finishReason := q.finishReason(qwenResponse.Choices[0].FinishReason)

	return &ProviderResponse{
		Content:      content,
		ToolCalls:    []message.ToolCall{}, // Qwen tool calling would need to be implemented
		Usage:        q.usage(qwenResponse),
		FinishReason: finishReason,
	}, nil
}

func (q *qwenClient) stream(ctx context.Context, messages []message.Message, tools []tools.BaseTool) <-chan ProviderEvent {
	params := q.preparedParams(q.convertMessages(messages), q.convertTools(tools))
	params.Stream = true

	eventChan := make(chan ProviderEvent)

	go func() {
		defer close(eventChan)

		// Set up the request
		baseURL := q.providerOptions.baseURL
		if baseURL == "" {
			// Try to get the resource URL from credentials
			credentialsPath := config.Get().Options.DataDirectory + "/credentials.json"
			credentials := make(map[string]map[string]string)
			if data, err := os.ReadFile(credentialsPath); err == nil {
				json.Unmarshal(data, &credentials)
			}
			
			if qwenCreds, exists := credentials["qwen"]; exists {
				if resourceURL, exists := qwenCreds["resource_url"]; exists && resourceURL != "" {
					baseURL = resourceURL
				}
			}
			
			// Fallback to the default DashScope endpoint
			if baseURL == "" {
				baseURL = "https://dashscope.aliyuncs.com/compatible-mode/v1"
			}
		}

		url := baseURL + "/chat/completions"
		jsonData, err := json.Marshal(params)
		if err != nil {
			eventChan <- ProviderEvent{
				Type:  EventError,
				Error: fmt.Errorf("failed to marshal request: %w", err),
			}
			return
		}

		req, err := http.NewRequestWithContext(ctx, "POST", url, strings.NewReader(string(jsonData)))
		if err != nil {
			eventChan <- ProviderEvent{
				Type:  EventError,
				Error: fmt.Errorf("failed to create request: %w", err),
			}
			return
		}

		// Set headers
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+q.providerOptions.apiKey)
		
		// Add required DashScope headers
		req.Header.Set("X-DashScope-CacheControl", "enable")
		req.Header.Set("X-DashScope-UserAgent", "Floss/1.0 (Windows; x86_64)")
		req.Header.Set("X-DashScope-AuthType", "qwen-oauth")

		// Add extra headers
		for key, value := range q.providerOptions.extraHeaders {
			req.Header.Set(key, value)
		}

		// Make the request
		resp, err := q.client.Do(req)
		if err != nil {
			eventChan <- ProviderEvent{
				Type:  EventError,
				Error: fmt.Errorf("failed to send request: %w", err),
			}
			return
		}

		// Check if we need to refresh token
		if resp.StatusCode == http.StatusUnauthorized {
			// Try to refresh token
			if refreshErr := q.refreshToken(); refreshErr != nil {
				eventChan <- ProviderEvent{
					Type:  EventError,
					Error: fmt.Errorf("token refresh failed: %w", refreshErr),
				}
				resp.Body.Close()
				return
			}

			// Retry the request with new token
			req.Header.Set("Authorization", "Bearer "+q.providerOptions.apiKey)
			resp.Body.Close()
			resp, err = q.client.Do(req)
			if err != nil {
				eventChan <- ProviderEvent{
					Type:  EventError,
					Error: fmt.Errorf("failed to send request after token refresh: %w", err),
				}
				return
			}
		}

		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			eventChan <- ProviderEvent{
				Type:  EventError,
				Error: fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body)),
			}
			return
		}

		// Process the stream
		reader := bufio.NewReader(resp.Body)
		currentContent := ""
		
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				if err == io.EOF {
					break
				}
				eventChan <- ProviderEvent{
					Type:  EventError,
					Error: fmt.Errorf("failed to read stream: %w", err),
				}
				return
			}

			// Skip empty lines
			if line == "" || line == "\n" {
				continue
			}

			// Handle SSE format
			if strings.HasPrefix(line, "data: ") {
				data := strings.TrimPrefix(line, "data: ")
				data = strings.TrimSpace(data)

				// Check for stream end
				if data == "[DONE]" {
					break
				}

				// Parse the chunk
				var chunk QwenStreamChunk
				if err := json.Unmarshal([]byte(data), &chunk); err != nil {
					slog.Debug("Failed to parse stream chunk", "error", err, "data", data)
					continue
				}

				if len(chunk.Choices) > 0 {
					content := chunk.Choices[0].Message.Content
					if content != "" {
						eventChan <- ProviderEvent{
							Type:    EventContentDelta,
							Content: content,
						}
						currentContent += content
					}
				}
			}
		}

		// Send completion event
		eventChan <- ProviderEvent{
			Type: EventComplete,
			Response: &ProviderResponse{
				Content:      currentContent,
				ToolCalls:    []message.ToolCall{},
				Usage:        TokenUsage{}, // Would need to implement usage tracking
				FinishReason: message.FinishReasonEndTurn,
			},
		}
	}()

	return eventChan
}

func (q *qwenClient) usage(response QwenChatCompletionResponse) TokenUsage {
	return TokenUsage{
		InputTokens:  int64(response.Usage.PromptTokens),
		OutputTokens: int64(response.Usage.CompletionTokens),
	}
}

func (q *qwenClient) Model() catwalk.Model {
	return q.providerOptions.model(q.providerOptions.modelType)
}