# Qwen Configuration Example

This directory contains an example configuration file for using Qwen with Floss.

## Usage

1. Authenticate with Qwen:
   ```bash
   floss auth qwen
   ```

2. Copy the `qwen_config.json` file to your Floss configuration directory:
   - On Linux/macOS: `~/.config/floss/floss.json`
   - On Windows: `%LOCALAPPDATA%\floss\floss.json`

3. Or merge the Qwen provider configuration with your existing configuration file.

4. Start Floss:
   ```bash
   floss
   ```

The configuration will automatically use your Qwen credentials stored by the authentication process.