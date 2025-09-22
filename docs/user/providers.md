# Providers Guide

This guide explains how to configure and use different AI providers with FLOSS.

## Supported Providers

FLOSS supports a wide range of AI providers out of the box:

| Provider | Environment Variable | Notes |
|----------|---------------------|-------|
| Anthropic | `ANTHROPIC_API_KEY` | |
| OpenAI | `OPENAI_API_KEY` | |
| OpenRouter | `OPENROUTER_API_KEY` | |
| Cerebras | `CEREBRAS_API_KEY` | |
| Google Gemini | `GEMINI_API_KEY` | |
| Google Cloud VertexAI | `VERTEXAI_PROJECT`, `VERTEXAI_LOCATION` | Requires `gcloud auth application-default login` |
| Groq | `GROQ_API_KEY` | |
| Qwen | `QWEN_API_KEY` | Supports OAuth2 device code flow |
| AWS Bedrock | `AWS_ACCESS_KEY_ID`, `AWS_SECRET_ACCESS_KEY`, `AWS_REGION` | Requires AWS configuration |
| Azure OpenAI | `AZURE_OPENAI_API_ENDPOINT`, `AZURE_OPENAI_API_KEY`, `AZURE_OPENAI_API_VERSION` | |

## Provider Configuration

### Qwen Provider

FLOSS supports Qwen as a provider. You can authenticate with Qwen using the OAuth2 device code flow:

```bash
floss auth qwen
```

This will guide you through the authentication process and store your credentials securely.

Once authenticated, you can configure Qwen as a provider by adding it to your configuration file:

```json
{
  "providers": {
    "qwen": {
      "name": "Qwen",
      "type": "openai",
      "base_url": "https://api.qwen.com/v1/",
      "api_key": "$QWEN_API_KEY",
      "models": [
        {
          "id": "qwen-max",
          "name": "Qwen Max",
          "context_window": 32768,
          "default_max_tokens": 2048
        }
      ]
    }
  }
}
```

### Disabling Provider Auto-Updates

By default, FLOSS automatically updates providers and models from Catwalk, the open source Floss provider database. This means that when new providers and models are added to Catwalk, they will automatically appear in Floss without requiring an update.

For air-gapped environments, this might not be what you want, and this feature can be disabled.

To disable automatic provider updates, set `disable_provider_auto_update` in your config:

```json
{
  "options": {
    "disable_provider_auto_update": true
  }
}
```

Or set the `FLOSS_DISABLE_PROVIDER_AUTO_UPDATE` environment variable:

```bash
export FLOSS_DISABLE_PROVIDER_AUTO_UPDATE=1
```

### Manually Updating Providers

Manually updating providers is possible with the `floss update-providers` command.

```bash
# Update providers remotely from Catwalk.
floss update-providers

# Update providers from a custom Catwalk base URL.
floss update-providers https://example.com/

# Update providers from a local file.
floss update-providers /path/to/local-providers.json

# Reset providers to the embedded version, embedded at floss at build time.
floss update-providers embedded

floss update-providers --help
```

## Custom Providers

FLOSS supports custom provider configurations for both OpenAI-compatible and Anthropic-compatible APIs.

### OpenAI-Compatible APIs

Here's an example configuration for Deepseek, which uses an OpenAI-compatible API. Don't forget to set `DEEPSEEK_API_KEY` in your environment.

```json
{
  "$schema": "https://nom-nom-hub.land/floss.json",
  "providers": {
    "deepseek": {
      "type": "openai",
      "base_url": "https://api.deepseek.com/v1",
      "api_key": "$DEEPSEEK_API_KEY",
      "models": [
        {
          "id": "deepseek-chat",
          "name": "Deepseek V3",
          "cost_per_1m_in": 0.27,
          "cost_per_1m_out": 1.1,
          "cost_per_1m_in_cached": 0.07,
          "cost_per_1m_out_cached": 1.1,
          "context_window": 64000,
          "default_max_tokens": 5000
        }
      ]
    }
  }
}
```

### Anthropic-Compatible APIs

Custom Anthropic-compatible providers follow this format:

```json
{
  "$schema": "https://nom-nom-hub.land/floss.json",
  "providers": {
    "custom-anthropic": {
      "type": "anthropic",
      "base_url": "https://api.anthropic.com/v1",
      "api_key": "$ANTHROPIC_API_KEY",
      "extra_headers": {
        "anthropic-version": "2023-06-01"
      },
      "models": [
        {
          "id": "claude-sonnet-4-20250514",
          "name": "Claude Sonnet 4",
          "cost_per_1m_in": 3,
          "cost_per_1m_out": 15,
          "cost_per_1m_in_cached": 3.75,
          "cost_per_1m_out_cached": 0.3,
          "context_window": 200000,
          "default_max_tokens": 50000,
          "can_reason": true,
          "supports_attachments": true
        }
      ]
    }
  }
}
```

## Local Models

Local models can also be configured via OpenAI-compatible API. Here are two common examples:

### Ollama

```json
{
  "providers": {
    "ollama": {
      "name": "Ollama",
      "base_url": "http://localhost:11434/v1/",
      "type": "openai",
      "models": [
        {
          "name": "Qwen 3 30B",
          "id": "qwen3:30b",
          "context_window": 256000,
          "default_max_tokens": 20000
        }
      ]
    }
  }
}
```

### LM Studio

```json
{
  "providers": {
    "lmstudio": {
      "name": "LM Studio",
      "base_url": "http://localhost:1234/v1/",
      "type": "openai",
      "models": [
        {
          "name": "Qwen 3 30B",
          "id": "qwen/qwen3-30b-a3b-2507",
          "context_window": 256000,
          "default_max_tokens": 20000
        }
      ]
    }
  }
}
```

## Provider Selection

You can specify which provider to use in several ways:

1. **Command Line**: Use the `--provider` flag
   ```bash
   floss --provider=openai
   ```

2. **Configuration**: Set a default provider in your configuration file
   ```json
   {
     "providers": {
       "default": "openai"
     }
   }
   ```

3. **Environment Variable**: Set the `FLOSS_DEFAULT_PROVIDER` environment variable
   ```bash
   export FLOSS_DEFAULT_PROVIDER=anthropic
   ```

## Switching Providers

FLOSS allows you to switch providers mid-session while preserving context. This is useful when you want to leverage the strengths of different models for different tasks.

To switch providers during a session, use the `/provider` command followed by the provider name:

```
/provider anthropic
```

## Provider-Specific Features

Different providers may support different features:

- **Reasoning**: Some providers support advanced reasoning capabilities
- **Attachments**: Some providers can process file attachments
- **Caching**: Some providers support cached completions for cost savings
- **Streaming**: Most providers support streaming responses for real-time output

Check your provider's documentation for specific feature support.