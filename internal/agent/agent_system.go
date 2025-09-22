package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/nom-nom-hub/floss/internal/csync"
	"github.com/nom-nom-hub/floss/internal/config"
	"github.com/nom-nom-hub/floss/internal/history"
	"github.com/nom-nom-hub/floss/internal/llm/agent"
	"github.com/nom-nom-hub/floss/internal/lsp"
	"github.com/nom-nom-hub/floss/internal/message"
	"github.com/nom-nom-hub/floss/internal/permission"
	"github.com/nom-nom-hub/floss/internal/session"
)

// AgentSystemConfig represents the configuration for the agent system
type AgentSystemConfig struct {
	CompanyStructure    CompanyStructure                     `json:"company_structure"`
	AgentDefinitions    map[AgentRole]AgentDefinition       `json:"agent_definitions"`
	CommunicationProtocols map[string]AgentCommunication    `json:"communication_protocols"`
	CollaborationProtocols map[string]CollaborationProtocol `json:"collaboration_protocols"`
	Workflows           map[string]AgentWorkflow            `json:"workflows"`
	Enabled             bool                                `json:"enabled"`
}

// DefaultAgentSystemConfig returns the default agent system configuration
func DefaultAgentSystemConfig() AgentSystemConfig {
	return AgentSystemConfig{
		CompanyStructure:    DefaultCompanyStructure(),
		AgentDefinitions:    DefaultAgentDefinitions(),
		CommunicationProtocols: DefaultCommunicationProtocols(),
		CollaborationProtocols: DefaultCollaborationProtocols(),
		Workflows:           DefaultWorkflows(),
		Enabled:             true,
	}
}

// LoadAgentSystemConfig loads the agent system configuration from a file
func LoadAgentSystemConfig(configDir string) (AgentSystemConfig, error) {
	configPath := filepath.Join(configDir, "agents.json")
	
	// Check if the file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// Return default configuration if file doesn't exist
		return DefaultAgentSystemConfig(), nil
	}

	// Read the file
	data, err := os.ReadFile(configPath)
	if err != nil {
		return AgentSystemConfig{}, fmt.Errorf("failed to read agent system config file: %w", err)
	}

	// Parse the JSON
	var config AgentSystemConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return AgentSystemConfig{}, fmt.Errorf("failed to parse agent system config: %w", err)
	}

	return config, nil
}

// SaveAgentSystemConfig saves the agent system configuration to a file
func SaveAgentSystemConfig(config AgentSystemConfig, configDir string) error {
	// Create the directory if it doesn't exist
	if err := os.MkdirAll(configDir, 0o755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Convert to JSON
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal agent system config: %w", err)
	}

	// Write to file
	configPath := filepath.Join(configDir, "agents.json")
	if err := os.WriteFile(configPath, data, 0o644); err != nil {
		return fmt.Errorf("failed to write agent system config file: %w", err)
	}

	return nil
}

// CreateAgentServices creates agent services for each defined agent role
func CreateAgentServices(
	cfg *config.Config,
	permissions permission.Service,
	sessions session.Service,
	messages message.Service,
	history history.Service,
	lspClients *csync.Map[string, *lsp.Client],
) (*csync.Map[AgentRole, agent.Service], error) {
	agentServices := csync.NewMap[AgentRole, agent.Service]()

	// Load agent system configuration
	agentConfig, err := LoadAgentSystemConfig(cfg.Options.DataDirectory)
	if err != nil {
		return nil, fmt.Errorf("failed to load agent system config: %w", err)
	}

	// If agent system is not enabled, return empty map
	if !agentConfig.Enabled {
		return agentServices, nil
	}

	// Create agent services for each defined agent
	for role, def := range agentConfig.AgentDefinitions {
		// Create agent configuration
		agentCfg := config.Agent{
			ID:           def.ID,
			Name:         def.Name,
			Description:  def.Description,
			Model:        def.Model,
			AllowedTools: def.AllowedTools,
			ContextPaths: def.ContextPaths,
		}

		// Create the agent service
		ctx := context.Background()
		agentService, err := agent.NewAgent(ctx, agentCfg, permissions, sessions, messages, history, lspClients)
		if err != nil {
			return nil, fmt.Errorf("failed to create agent service for role %s: %w", role, err)
		}

		// Add to the map
		agentServices.Set(role, agentService)
	}

	return agentServices, nil
}

// GetAgentDefinition returns the definition for a specific agent role
func GetAgentDefinition(cfg *config.Config, role AgentRole) (AgentDefinition, bool) {
	// Load agent system configuration
	agentConfig, err := LoadAgentSystemConfig(cfg.Options.DataDirectory)
	if err != nil {
		return AgentDefinition{}, false
	}

	def, exists := agentConfig.AgentDefinitions[role]
	return def, exists
}

// GetWorkflow returns a specific workflow by ID
func GetWorkflow(cfg *config.Config, workflowID string) (AgentWorkflow, bool) {
	// Load agent system configuration
	agentConfig, err := LoadAgentSystemConfig(cfg.Options.DataDirectory)
	if err != nil {
		return AgentWorkflow{}, false
	}

	workflow, exists := agentConfig.Workflows[workflowID]
	return workflow, exists
}

// ListAgentRoles returns a list of all defined agent roles
func ListAgentRoles(cfg *config.Config) []AgentRole {
	// Load agent system configuration
	agentConfig, err := LoadAgentSystemConfig(cfg.Options.DataDirectory)
	if err != nil {
		return []AgentRole{}
	}

	var roles []AgentRole
	for role := range agentConfig.AgentDefinitions {
		roles = append(roles, role)
	}

	return roles
}

// ListWorkflows returns a list of all defined workflows
func ListWorkflows(cfg *config.Config) []AgentWorkflow {
	// Load agent system configuration
	agentConfig, err := LoadAgentSystemConfig(cfg.Options.DataDirectory)
	if err != nil {
		return []AgentWorkflow{}
	}

	var workflows []AgentWorkflow
	for _, workflow := range agentConfig.Workflows {
		workflows = append(workflows, workflow)
	}

	return workflows
}