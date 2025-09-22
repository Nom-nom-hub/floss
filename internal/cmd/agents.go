package cmd

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/nom-nom-hub/floss/internal/agent"
	"github.com/spf13/cobra"
)

// agentsCmd represents the agents command
var agentsCmd = &cobra.Command{
	Use:   "agents",
	Short: "Manage virtual company agents",
	Long:  `Manage the virtual company agents that simulate a full development team.`,
}

// agentsListCmd represents the agents list command
var agentsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all available agent roles",
	Long:  `List all available agent roles in the virtual company.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Initialize config like the main command does
		appInstance, err := SetupApp(cmd)
		if err != nil {
			return err
		}
		cfg := appInstance.Config()

		// Load agent system configuration
		agentConfig, err := agent.LoadAgentSystemConfig(cfg.Options.DataDirectory)
		if err != nil {
			return fmt.Errorf("failed to load agent system config: %w", err)
		}

		// Print company structure
		fmt.Printf("Company: %s\n", agentConfig.CompanyStructure.Name)
		fmt.Printf("Description: %s\n\n", agentConfig.CompanyStructure.Description)

		// Print departments and roles
		for i, dept := range agentConfig.CompanyStructure.Departments {
			fmt.Printf("Department %d: %s\n", i, dept.Name)
			fmt.Printf("  Description: %s\n", dept.Description)
			if dept.HeadRole != "" {
				def, exists := agentConfig.AgentDefinitions[dept.HeadRole]
				if exists {
					fmt.Printf("  Head: %s (%s)\n", def.Name, def.Description)
				} else {
					fmt.Printf("  Head: %s (Definition NOT FOUND)\n", dept.HeadRole)
				}
			}

			for j, team := range dept.Teams {
				fmt.Printf("  Team %d: %s\n", j, team.Name)
				fmt.Printf("    Description: %s\n", team.Description)
				if team.LeadRole != "" {
					def, exists := agentConfig.AgentDefinitions[team.LeadRole]
					if exists {
						fmt.Printf("    Lead: %s (%s)\n", def.Name, def.Description)
					} else {
						fmt.Printf("    Lead: %s (Definition NOT FOUND)\n", team.LeadRole)
					}
				}

				for k, role := range team.MemberRoles {
					def, exists := agentConfig.AgentDefinitions[role]
					if exists {
						fmt.Printf("    Member %d: %s (%s)\n", k, def.Name, def.Description)
					} else {
						fmt.Printf("    Member %d: %s (Definition NOT FOUND)\n", k, role)
					}
				}
			}
			fmt.Println()
		}

		return nil
	},
}

// agentsInfoCmd represents the agents info command
var agentsInfoCmd = &cobra.Command{
	Use:   "info [role]",
	Short: "Get detailed information about a specific agent role",
	Long:  `Get detailed information about a specific agent role including capabilities and tools.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		// Initialize config like the main command does
		appInstance, err := SetupApp(cmd)
		if err != nil {
			return err
		}
		cfg := appInstance.Config()

		role := agent.AgentRole(args[0])

		// Load agent system configuration
		agentConfig, err := agent.LoadAgentSystemConfig(cfg.Options.DataDirectory)
		if err != nil {
			return fmt.Errorf("failed to load agent system config: %w", err)
		}

		// Get agent definition
		def, exists := agentConfig.AgentDefinitions[role]
		if !exists {
			return fmt.Errorf("agent role %s not found", role)
		}

		// Print agent information
		fmt.Printf("Agent Role: %s\n", def.Name)
		fmt.Printf("ID: %s\n", def.ID)
		fmt.Printf("Description: %s\n", def.Description)
		fmt.Printf("Role: %s\n", def.Role)
		fmt.Printf("Model: %s\n", def.Model)
		
		fmt.Printf("Allowed Tools: %s\n", strings.Join(def.AllowedTools, ", "))
		fmt.Printf("Context Paths: %s\n", strings.Join(def.ContextPaths, ", "))
		
		fmt.Printf("Capabilities:\n")
		for _, cap := range def.Capabilities {
			fmt.Printf("  - %s\n", cap)
		}
		
		fmt.Printf("Prompt Template:\n%s\n", def.PromptTemplate)

		return nil
	},
}

// agentsWorkflowsCmd represents the agents workflows command
var agentsWorkflowsCmd = &cobra.Command{
	Use:   "workflows",
	Short: "List all available workflows",
	Long:  `List all available workflows that agents can execute.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Initialize config like the main command does
		appInstance, err := SetupApp(cmd)
		if err != nil {
			return err
		}
		cfg := appInstance.Config()

		// Load agent system configuration
		agentConfig, err := agent.LoadAgentSystemConfig(cfg.Options.DataDirectory)
		if err != nil {
			return fmt.Errorf("failed to load agent system config: %w", err)
		}

		// Print workflows
		fmt.Printf("Available Workflows:\n\n")
		for _, workflow := range agentConfig.Workflows {
			fmt.Printf("Workflow: %s\n", workflow.Name)
			fmt.Printf("  ID: %s\n", workflow.ID)
			fmt.Printf("  Description: %s\n", workflow.Description)
			fmt.Printf("  Steps: %d\n\n", len(workflow.Steps))
		}

		return nil
	},
}

// agentsInitCmd represents the agents init command
var agentsInitCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize the agent system with default configuration",
	Long:  `Initialize the agent system with default configuration files.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Initialize config like the main command does
		appInstance, err := SetupApp(cmd)
		if err != nil {
			return err
		}
		cfg := appInstance.Config()

		// Create default agent system configuration
		agentConfig := agent.DefaultAgentSystemConfig()

		// Save the configuration
		if err := agent.SaveAgentSystemConfig(agentConfig, cfg.Options.DataDirectory); err != nil {
			return fmt.Errorf("failed to save agent system config: %w", err)
		}

		fmt.Printf("Agent system initialized with default configuration in %s\n", cfg.Options.DataDirectory)
		return nil
	},
}

// agentsConfigCmd represents the agents config command
var agentsConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "Show current agent system configuration",
	Long:  `Show the current agent system configuration.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Initialize config like the main command does
		appInstance, err := SetupApp(cmd)
		if err != nil {
			return err
		}
		cfg := appInstance.Config()

		// Load agent system configuration
		agentConfig, err := agent.LoadAgentSystemConfig(cfg.Options.DataDirectory)
		if err != nil {
			return fmt.Errorf("failed to load agent system config: %w", err)
		}

		// Convert to JSON for display
		data, err := json.MarshalIndent(agentConfig, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal agent system config: %w", err)
		}

		fmt.Printf("%s\n", data)
		return nil
	},
}

// agentsEnableCmd represents the agents enable command
var agentsEnableCmd = &cobra.Command{
	Use:   "enable",
	Short: "Enable the agent system",
	Long:  `Enable the agent system to allow virtual company simulation.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Initialize config like the main command does
		appInstance, err := SetupApp(cmd)
		if err != nil {
			return err
		}
		cfg := appInstance.Config()

		// Load agent system configuration
		agentConfig, err := agent.LoadAgentSystemConfig(cfg.Options.DataDirectory)
		if err != nil {
			// If config doesn't exist, create default
			agentConfig = agent.DefaultAgentSystemConfig()
		}

		// Enable the agent system
		agentConfig.Enabled = true

		// Save the configuration
		if err := agent.SaveAgentSystemConfig(agentConfig, cfg.Options.DataDirectory); err != nil {
			return fmt.Errorf("failed to save agent system config: %w", err)
		}

		fmt.Println("Agent system enabled")
		return nil
	},
}

// agentsDisableCmd represents the agents disable command
var agentsDisableCmd = &cobra.Command{
	Use:   "disable",
	Short: "Disable the agent system",
	Long:  `Disable the agent system to stop virtual company simulation.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Initialize config like the main command does
		appInstance, err := SetupApp(cmd)
		if err != nil {
			return err
		}
		cfg := appInstance.Config()

		// Load agent system configuration
		agentConfig, err := agent.LoadAgentSystemConfig(cfg.Options.DataDirectory)
		if err != nil {
			// If config doesn't exist, create default
			agentConfig = agent.DefaultAgentSystemConfig()
		}

		// Disable the agent system
		agentConfig.Enabled = false

		// Save the configuration
		if err := agent.SaveAgentSystemConfig(agentConfig, cfg.Options.DataDirectory); err != nil {
			return fmt.Errorf("failed to save agent system config: %w", err)
		}

		fmt.Println("Agent system disabled")
		return nil
	},
}

func init() {
	// Add subcommands to agents command
	agentsCmd.AddCommand(agentsListCmd)
	agentsCmd.AddCommand(agentsInfoCmd)
	agentsCmd.AddCommand(agentsWorkflowsCmd)
	agentsCmd.AddCommand(agentsInitCmd)
	agentsCmd.AddCommand(agentsConfigCmd)
	agentsCmd.AddCommand(agentsEnableCmd)
	agentsCmd.AddCommand(agentsDisableCmd)

	// Add agents command to root command
	rootCmd.AddCommand(agentsCmd)
}