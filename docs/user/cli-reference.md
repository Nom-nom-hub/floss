# FLOSS CLI Reference

This document provides a comprehensive reference for all FLOSS command-line interface commands.

## Global Options

| Option | Description |
|--------|-------------|
| `--help` | Display help for the command |
| `--version` | Display the version of FLOSS |
| `--config` | Specify a configuration file path |
| `--debug` | Enable debug logging |

## Main Commands

### `floss`

Start an interactive FLOSS session.

```bash
floss
```

Options:
- `--provider <provider>`: Specify which provider to use
- `--model <model>`: Specify which model to use
- `--session <session>`: Load a specific session
- `--prompt <prompt>`: Provide an initial prompt

### `floss auth`

Authenticate with providers that require it.

```bash
floss auth <provider>
```

Currently supports:
- `qwen`: Authenticate with Qwen using OAuth2 device code flow

### `floss logs`

View FLOSS logs.

```bash
floss logs
```

Options:
- `--tail <number>`: Show last N lines (default: 1000)
- `--follow`: Follow logs in real time

### `floss agents`

Manage the Virtual Company Agent System.

```bash
floss agents <subcommand>
```

Subcommands:
- `list`: List available agent roles
- `info <role>`: Get information about a specific agent role
- `workflows`: List available workflows
- `init`: Initialize agent system with default configuration
- `config`: Show current agent system configuration
- `enable`: Enable the agent system
- `disable`: Disable the agent system

### `floss update-providers`

Update provider information.

```bash
floss update-providers [source]
```

Sources:
- Remote Catwalk repository (default)
- Custom URL
- Local file
- Embedded (reset to build-time providers)

### `floss run`

Run a command non-interactively.

```bash
floss run <command>
```

### `floss schema`

Generate JSON schema for configuration.

```bash
floss schema
```

## Environment Variables

FLOSS can be configured using environment variables. See [Environment Variables Reference](../reference/environment-variables.md) for details.