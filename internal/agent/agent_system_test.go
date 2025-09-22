package agent

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDefaultAgentSystemConfig(t *testing.T) {
	cfg := DefaultAgentSystemConfig()
	require.NotNil(t, cfg.CompanyStructure)
	require.NotEmpty(t, cfg.AgentDefinitions)
	require.NotEmpty(t, cfg.CommunicationProtocols)
	require.NotEmpty(t, cfg.CollaborationProtocols)
	require.NotEmpty(t, cfg.Workflows)
	require.True(t, cfg.Enabled)
}

func TestLoadAgentSystemConfig(t *testing.T) {
	t.Parallel()
	tdir := t.TempDir()

	// Case 1: config file does not exist, should return default config
	cfg, err := LoadAgentSystemConfig(tdir)
	require.NoError(t, err)
	require.NotNil(t, cfg.CompanyStructure)
	require.NotEmpty(t, cfg.AgentDefinitions)
	require.True(t, cfg.Enabled)

	// Case 2: config file exists and is valid
	validConfigPath := filepath.Join(tdir, "agents.json")
	err = os.WriteFile(validConfigPath, []byte(`{
		"company_structure": {
			"name": "Test Company",
			"description": "A test company",
			"departments": []
		},
		"agent_definitions": {},
		"communication_protocols": {},
		"collaboration_protocols": {},
		"workflows": {},
		"enabled": true
	}`), 0o644)
	require.NoError(t, err)

	cfg, err = LoadAgentSystemConfig(tdir)
	require.NoError(t, err)
	require.Equal(t, "Test Company", cfg.CompanyStructure.Name)
	require.True(t, cfg.Enabled)

	// Case 3: config file exists but is invalid JSON
	invalidConfigPath := filepath.Join(tdir, "agents.json")
	err = os.WriteFile(invalidConfigPath, []byte(`invalid json`), 0o644)
	require.NoError(t, err)

	_, err = LoadAgentSystemConfig(tdir)
	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to parse agent system config")
}
