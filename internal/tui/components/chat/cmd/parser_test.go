package cmd

import (
	"testing"

	"github.com/nom-nom-hub/floss/internal/app"
)

func TestCommandParsing(t *testing.T) {
	// Create a mock app
	mockApp := &app.App{}

	parser := NewParser(mockApp)

	// Test parsing a command
	cmd := parser.ParseCommand("/list agents")
	if cmd == nil {
		t.Error("Expected command to be parsed, got nil")
	}
	if cmd.Name != "list" {
		t.Errorf("Expected command name 'list', got '%s'", cmd.Name)
	}
	if len(cmd.Args) != 1 || cmd.Args[0] != "agents" {
		t.Errorf("Expected args ['agents'], got %v", cmd.Args)
	}

	// Test parsing a non-command
	cmd = parser.ParseCommand("hello world")
	if cmd != nil {
		t.Error("Expected nil for non-command, got command")
	}

	// Test checking if text is a command
	if !parser.IsCommand("/help") {
		t.Error("Expected '/help' to be recognized as a command")
	}
	if parser.IsCommand("hello") {
		t.Error("Expected 'hello' to not be recognized as a command")
	}
}