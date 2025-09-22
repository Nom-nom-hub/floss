# FAQ

Frequently Asked Questions about FLOSS.

## What is FLOSS?

FLOSS is a terminal-based AI coding assistant that helps software developers with their daily tasks. It provides multi-model support with various AI providers, session-based context management, Language Server Protocol (LSP) integration, and Model Context Protocol (MCP) extensibility. FLOSS works directly in your terminal, making it easy to integrate into your existing development workflow.

## How is FLOSS different from other AI coding assistants?

FLOSS has several unique features that set it apart from other AI coding assistants:

1. **Terminal-based**: FLOSS runs directly in your terminal, making it lightweight and easy to integrate into your existing workflow.

2. **Multi-model support**: FLOSS supports a wide range of AI providers, allowing you to choose the best model for each task.

3. **Session-based context**: FLOSS maintains context across interactions, allowing for more coherent and helpful conversations.

4. **LSP integration**: FLOSS integrates with Language Server Protocol servers to provide additional context about your code.

5. **MCP extensibility**: FLOSS supports the Model Context Protocol, allowing you to extend its functionality with custom tools and services.

6. **Cross-platform**: FLOSS works on macOS, Linux, Windows, and other platforms.

## How much does FLOSS cost?

FLOSS itself is free and open-source software. However, you will need to pay for access to the AI providers you use with FLOSS. Each provider has its own pricing model, which you can find on their respective websites.

## Do I need to be online to use FLOSS?

FLOSS requires an internet connection to communicate with AI providers. However, some features may work offline, such as viewing previous sessions and basic text editing.

## Which AI providers does FLOSS support?

FLOSS supports a wide range of AI providers, including:

- Anthropic (Claude)
- OpenAI (GPT)
- Google (Gemini)
- Qwen
- Groq
- Mistral
- Together AI
- And many more

You can also configure FLOSS to work with local models through Ollama, LM Studio, and other local AI services.

## How do I configure FLOSS to use a specific AI provider?

You can configure FLOSS to use a specific AI provider by setting environment variables or by configuring providers in your FLOSS configuration file. See the [Providers Guide](providers.md) for detailed instructions.

## Can I use FLOSS with local AI models?

Yes, FLOSS can be configured to work with local AI models through services like Ollama and LM Studio. See the [Providers Guide](providers.md) for instructions on configuring local models.

## How do I update FLOSS?

The method for updating FLOSS depends on how you installed it:

- **Homebrew**: `brew upgrade floss`
- **NPM**: `npm update -g @nom-nom-hub/floss`
- **Go**: `go install github.com/nom-nom-hub/floss@latest`
- **Manual installation**: Download the latest release from GitHub

## How do I get help with FLOSS?

You can get help with FLOSS in several ways:

1. Read the documentation in the `docs/` directory
2. Use the `floss --help` command for command-line help
3. Check the [GitHub issues](https://github.com/nom-nom-hub/floss/issues) to see if others have encountered the same issue
4. Join the community chat or forum for real-time help
5. Contact support if you have a commercial license

## How can I contribute to FLOSS?

FLOSS is an open-source project, and contributions are welcome! You can contribute in several ways:

1. Report bugs and suggest features on GitHub
2. Submit pull requests with bug fixes or new features
3. Improve the documentation
4. Help answer questions from other users
5. Spread the word about FLOSS

See the [Contributing Guide](../developer/contributing.md) for more details on how to contribute.

## Is FLOSS secure?

FLOSS takes security seriously. Here are some key security features:

1. **Local-first**: FLOSS runs locally on your machine, so your code stays on your machine.

2. **Controlled data sharing**: You control what information is sent to AI providers.

3. **Secure credential storage**: FLOSS stores credentials securely using your system's keychain.

4. **Open source**: FLOSS is open source, so you can audit the code for security issues.

However, you should still be careful about what information you share with AI providers, as that information will be subject to their privacy policies.

## How do I report a security issue?

If you discover a security issue with FLOSS, please report it by emailing the maintainers directly at security@nom-nom-hub.land. Do not report security issues in public forums or GitHub issues.

## Can I use FLOSS in a commercial setting?

Yes, FLOSS can be used in commercial settings. FLOSS is released under the FSL-1.1-MIT license, which allows for commercial use. However, you should review the license to understand your rights and obligations.

## How do I cite FLOSS in academic work?

If you use FLOSS in academic work, please cite it as follows:

```
FLOSS Development Team. (2025). FLOSS: A Terminal-Based AI Coding Assistant [Computer software]. https://github.com/nom-nom-hub/floss
```

Please also cite any specific AI providers you used through FLOSS.

## Where can I find more information about FLOSS?

You can find more information about FLOSS in several places:

1. The [README](../../README.md) file
2. The documentation in the `docs/` directory
3. The [GitHub repository](https://github.com/nom-nom-hub/floss)
4. The community chat or forum