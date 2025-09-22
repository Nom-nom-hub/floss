# Architecture Overview

This document provides a high-level overview of the FLOSS architecture.

## System Components

### Core Engine

The core engine handles:
- Session management
- Context processing
- Tool execution
- Message routing
- Provider integration

### Providers

Providers are the AI model interfaces. FLOSS supports multiple provider types:
- OpenAI-compatible APIs
- Anthropic-compatible APIs
- Google Gemini
- Qwen
- Amazon Bedrock
- Vertex AI
- Custom providers

### Language Server Protocol (LSP) Integration

FLOSS integrates with LSP servers to provide additional context about your code:
- Code completion
- Error detection
- Refactoring suggestions
- Symbol information

### Model Context Protocol (MCP) Integration

MCP servers extend FLOSS functionality:
- Filesystem access
- HTTP services
- Custom tools and data sources

### Agent System

The Virtual Company Agent System simulates a complete software development company:
- Specialized AI agents for different roles
- Communication protocols
- Collaboration workflows
- Multi-agent problem solving

## Data Flow

1. User input is processed by the core engine
2. Context is gathered from:
   - Current session
   - LSP servers
   - MCP servers
   - File system
3. Request is sent to the configured AI provider
4. Response is processed and any tool calls are executed
5. Results are presented to the user

## Session Management

FLOSS maintains sessions to preserve context:
- Message history
- Tool call results
- File modifications
- Agent interactions

## Extensibility Points

FLOSS is designed to be extensible:
- Custom providers
- Additional LSP servers
- New MCP integrations
- Custom tools
- Agent system customization

## Security Architecture

FLOSS follows security best practices:
- Permission system for tool execution
- Environment variable isolation
- Secure credential storage
- Network access controls