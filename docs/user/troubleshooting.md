# Troubleshooting Guide

This guide provides solutions to common issues you might encounter when using FLOSS.

## Installation Issues

### Command Not Found

If you get a "command not found" error when trying to run `floss`:

1. Check that FLOSS is installed:
   ```bash
   which floss
   ```
   
2. If FLOSS is installed but not in your PATH, add it to your PATH:
   ```bash
   export PATH=$PATH:/path/to/floss
   ```

3. On Windows, ensure the installation directory is in your PATH environment variable.

### Permission Denied

If you get a "permission denied" error:

1. Make sure the FLOSS binary is executable:
   ```bash
   chmod +x /path/to/floss
   ```

## Provider Configuration Issues

### API Key Errors

If you're getting API key errors:

1. Verify that your API key environment variables are set:
   ```bash
   echo $OPENAI_API_KEY
   ```
   
2. Check that your API key is valid and has not expired.

3. Ensure that your API key has the necessary permissions.

### Provider Not Found

If you get a "provider not found" error:

1. Check your configuration file for typos in provider names.

2. Ensure the provider is supported by your version of FLOSS.

3. Update providers if necessary:
   ```bash
   floss update-providers
   ```

### Authentication Failures

For providers that require authentication:

1. Run the authentication command:
   ```bash
   floss auth <provider>
   ```
   
2. Follow the prompts to complete authentication.

3. Check that your credentials are stored correctly.

## Session Issues

### Session Not Found

If you get a "session not found" error:

1. List available sessions to verify the session name:
   ```bash
   floss sessions list
   ```
   
2. Check for typos in the session name.

3. Verify that the session file exists in the sessions directory.

### Corrupted Sessions

If a session appears to be corrupted:

1. Try to load the session in debug mode:
   ```bash
   floss --debug --session=my-session
   ```
   
2. If the session cannot be recovered, you may need to delete it and start over:
   ```bash
   floss sessions delete my-session
   ```

## Performance Issues

### Slow Responses

If you're experiencing slow responses:

1. Check your internet connection.

2. Try switching to a different provider or model.

3. Verify that local services (like LSP servers) are not causing delays.

4. Enable debug mode to see where the delay is occurring:
   ```bash
   floss --debug
   ```

### High Memory Usage

If FLOSS is using too much memory:

1. Check if you have too many large context files loaded.

2. Reduce the number of context files in your configuration.

3. Close unused sessions.

4. Restart FLOSS to clear memory.

## LSP Integration Issues

### LSP Server Not Starting

If LSP servers are not starting:

1. Verify that the LSP server is installed and in your PATH:
   ```bash
   which gopls
   ```
   
2. Check your LSP configuration in your FLOSS configuration file.

3. Ensure that the LSP server supports the language you're working with.

### LSP Features Not Working

If LSP features are not working:

1. Check that the LSP server is running:
   ```bash
   ps aux | grep gopls
   ```
   
2. Verify that your project has the necessary files for the LSP server to work (e.g., go.mod for Go projects).

3. Check the LSP server logs for errors.

## MCP Integration Issues

### MCP Server Not Connecting

If MCP servers are not connecting:

1. Verify that the MCP server is installed and running.

2. Check your MCP configuration in your FLOSS configuration file.

3. Ensure that network connections to the MCP server are not blocked by firewalls.

### MCP Tools Not Available

If MCP tools are not available:

1. Check that the MCP server supports the tools you're trying to use.

2. Verify that the tools are properly configured in the MCP server.

3. Check the MCP server logs for errors.

## Network Issues

### Connection Timeouts

If you're experiencing connection timeouts:

1. Check your internet connection.

2. Verify that your API keys are correct and have not been rate-limited.

3. Try using a different network or proxy.

### SSL/TLS Errors

If you're getting SSL/TLS errors:

1. Ensure that your system has up-to-date SSL certificates.

2. Check if you're behind a proxy that might be interfering with SSL connections.

3. Try disabling SSL verification (not recommended for production):
   ```bash
   export NODE_TLS_REJECT_UNAUTHORIZED=0
   ```

## Configuration Issues

### Invalid Configuration

If you're getting configuration errors:

1. Validate your configuration file against the schema:
   ```bash
   floss schema validate
   ```
   
2. Check for syntax errors in your JSON.

3. Ensure all required fields are present.

### Configuration Not Loading

If your configuration is not being loaded:

1. Verify that your configuration file is in the correct location.

2. Check that the file has the correct name (`.floss.json` or `floss.json`).

3. Ensure that the file is readable:
   ```bash
   cat .floss.json
   ```

## Logging and Debugging

### Enabling Debug Logging

To enable debug logging:

```bash
floss --debug
```

This will provide more detailed information about what FLOSS is doing, which can help diagnose issues.

### Viewing Logs

To view FLOSS logs:

```bash
floss logs
```

To follow logs in real-time:

```bash
floss logs --follow
```

### Reporting Issues

If you encounter an issue that you cannot resolve:

1. Enable debug logging and reproduce the issue.

2. Check the logs for error messages.

3. Report the issue on the FLOSS GitHub repository with:
   - A description of the issue
   - Steps to reproduce
   - Debug logs
   - Your configuration file (with sensitive information removed)
   - Your operating system and FLOSS version

## Platform-Specific Issues

### Windows Issues

On Windows, you might encounter:

1. **Line ending issues**: Use LF line endings instead of CRLF in configuration files.

2. **Path issues**: Use forward slashes in paths, even on Windows.

3. **Terminal compatibility**: Use PowerShell or Windows Terminal for the best experience.

### macOS Issues

On macOS, you might encounter:

1. **Gatekeeper issues**: If you downloaded FLOSS from the internet, you might need to allow it in System Preferences > Security & Privacy.

2. **Homebrew issues**: Ensure Homebrew is up to date:
   ```bash
   brew update && brew upgrade
   ```

### Linux Issues

On Linux, you might encounter:

1. **Dependency issues**: Install required system dependencies:
   ```bash
   sudo apt-get install libwebkit2gtk-4.0-dev libgtk-3-dev libappindicator3-dev
   ```

2. **Permission issues**: Ensure FLOSS has the necessary permissions to access files and directories.

## Version Compatibility

### Updating FLOSS

To update FLOSS:

```bash
# If installed with Homebrew
brew upgrade floss

# If installed with NPM
npm update -g @nom-nom-hub/floss

# If installed with Go
go install github.com/nom-nom-hub/floss@latest
```

### Downgrading FLOSS

If you need to downgrade FLOSS:

1. Uninstall the current version.

2. Install the desired version:
   ```bash
   # For Homebrew
   brew install floss@version
   
   # For NPM
   npm install -g @nom-nom-hub/floss@version
   
   # For Go
   go install github.com/nom-nom-hub/floss@version
   ```

## Getting Help

If you're still having issues:

1. Check the [FAQ](faq.md) for answers to common questions.

2. Search the [GitHub issues](https://github.com/nom-nom-hub/floss/issues) to see if others have encountered the same issue.

3. Join the community chat or forum for real-time help.

4. Contact support if you have a commercial license.