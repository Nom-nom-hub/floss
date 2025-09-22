package cmd

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/nom-nom-hub/floss/internal/app"
	"github.com/nom-nom-hub/floss/internal/tui/util"
)

// Command represents a parsed command with its arguments
type Command struct {
	Name string
	Args []string
}

// Parser handles parsing and executing chat commands
type Parser struct {
	app *app.App
}

// NewParser creates a new command parser
func NewParser(app *app.App) *Parser {
	return &Parser{
		app: app,
	}
}

// ParseCommand parses a string into a command and arguments
func (p *Parser) ParseCommand(input string) *Command {
	trimmed := strings.TrimSpace(input)
	if !strings.HasPrefix(trimmed, "/") {
		return nil
	}

	// Remove the leading slash
	trimmed = strings.TrimPrefix(trimmed, "/")
	parts := strings.Fields(trimmed)
	
	if len(parts) == 0 {
		return nil
	}

	return &Command{
		Name: parts[0],
		Args: parts[1:],
	}
}

// IsCommand checks if the input is a command
func (p *Parser) IsCommand(input string) bool {
	return strings.HasPrefix(strings.TrimSpace(input), "/")
}

// ExecuteCommand executes a parsed command
func (p *Parser) ExecuteCommand(cmd *Command) tea.Cmd {
	switch strings.ToLower(cmd.Name) {
	case "list":
		// Check if the first argument is "agents"
		if len(cmd.Args) > 0 && strings.ToLower(cmd.Args[0]) == "agents" {
			return p.handleListAgents()
		}
		return util.ReportError(fmt.Errorf("Unknown command: %s. Type /help for available commands.", cmd.Name))
	case "help":
		return p.handleHelp()
	default:
		return util.ReportError(fmt.Errorf("Unknown command: %s. Type /help for available commands.", cmd.Name))
	}
}

// handleListAgents handles the "list agents" command
func (p *Parser) handleListAgents() tea.Cmd {
	if p.app.CoderAgent == nil {
		return util.ReportInfo("No agents configured")
	}

	// Get information about the coder agent
	agentInfo := "Available agents:\n"
	agentInfo += "  coder - The main coding assistant agent\n"
	
	// Add more detailed information if available
	if p.app.Config() != nil && p.app.Config().Agents != nil {
		for name, agentCfg := range p.app.Config().Agents {
			modelName := "unknown"
			if agentCfg.Model != "" {
				// Convert SelectedModelType to string
				modelName = string(agentCfg.Model)
			}
			agentInfo += fmt.Sprintf("    %s (model: %s)\n", name, modelName)
		}
	}

	return util.ReportInfo(agentInfo)
}

// handleHelp handles the "help" command
func (p *Parser) handleHelp() tea.Cmd {
	helpText := `Available commands:
  /list agents    - List available agents
  /help           - Show this help message`
	return util.ReportInfo(helpText)
}