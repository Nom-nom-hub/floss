# Configuration Guide

This guide explains how to configure FLOSS for your specific needs.

## Configuration Files

FLOSS looks for configuration files in the following locations, in order of priority:

1. `.floss.json` (in the current directory)
2. `floss.json` (in the current directory)
3. `$HOME/.config/floss/floss.json` (Unix) or `%USERPROFILE%\AppData\Local\floss\floss.json` (Windows)

## Basic Configuration Structure

A basic FLOSS configuration file looks like this:

```json
{
  "$schema": "https://charm.land/floss.json",
  "options": {
    "debug": false
  }
}
```

## Configuration Options

### Models

Configure default models for different tasks:

```json
{
  "models": {
    "large": {
      "model": "gpt-4o",
      "provider": "openai"
    },
    "small": {
      "model": "gpt-4o-mini",
      "provider": "openai"
    }
  }
}
```

### Providers

Configure AI providers:

```json
{
  "providers": {
    "openai": {
      "name": "OpenAI",
      "base_url": "https://api.openai.com/v1/",
      "type": "openai",
      "api_key": "$OPENAI_API_KEY",
      "models": [
        {
          "id": "gpt-4o",
          "name": "GPT-4 Omni",
          "context_window": 128000,
          "default_max_tokens": 4096
        }
      ]
    }
  }
}
```

### Language Server Protocol (LSP)

Configure LSP servers for additional context:

```json
{
  "lsp": {
    "go": {
      "command": "gopls",
      "env": {
        "GOTOOLCHAIN": "go1.23.0"
      }
    }
  }
}
```

### Model Context Protocol (MCP)

Configure MCP servers for extended functionality:

```json
{
  "mcp": {
    "filesystem": {
      "type": "stdio",
      "command": "npx",
      "args": ["@modelcontextprotocol/server-filesystem", "--stdio"]
    }
  }
}
```

### Options

General application options:

```json
{
  "options": {
    "context_paths": [
      "README.md",
      "docs/"
    ],
    "debug": false,
    "disable_provider_auto_update": false,
    "attribution": {
      "co_authored_by": true,
      "generated_with": true
    }
  }
}
```

### Permissions

Control which tools can be used without prompting:

```json
{
  "permissions": {
    "allowed_tools": [
      "view",
      "ls",
      "grep"
    ]
  }
}
```

## Context Files

FLOSS uses context files to provide additional information to the AI. These files are specified in the `context_paths` option and can include documentation, coding standards, or any other relevant information.

## Validation

FLOSS provides a JSON schema for configuration validation. You can reference it in your configuration file:

```json
{
  "$schema": "https://charm.land/floss.json"
}
```

See [Configuration Schema Reference](../reference/configuration-schema.md) for detailed information about all configuration options.