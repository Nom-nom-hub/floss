# Installation Guide

This guide provides detailed instructions for installing FLOSS on different operating systems.

## Prerequisites

Before installing FLOSS, ensure you have:

- A terminal or command prompt
- Git (for version control integration)
- Access to at least one AI provider (OpenAI, Anthropic, etc.)

## Installation Methods

### Package Managers

#### Homebrew (macOS/Linux)

```bash
brew install nom-nom-hub/tap/floss
```

#### NPM

```bash
npm install -g @nom-nom-hub/floss
```

#### Arch Linux

```bash
yay -S floss-bin
```

#### Nix

```bash
nix run github:numtide/nix-ai-tools#floss
```

### Windows

#### Winget

```bash
winget install nom-nom-hub.floss
```

#### Scoop

```bash
scoop bucket add floss https://github.com/nom-nom-hub/scoop-bucket.git
scoop install floss
```

### Linux Package Managers

#### Debian/Ubuntu

```bash
sudo mkdir -p /etc/apt/keyrings
curl -fsSL https://repo.nom-nom-hub.sh/apt/gpg.key | sudo gpg --dearmor -o /etc/apt/keyrings/floss.gpg
echo "deb [signed-by=/etc/apt/keyrings/floss.gpg] https://repo.nom-nom-hub.sh/apt/ * *" | sudo tee /etc/apt/sources.list.d/floss.list
sudo apt update && sudo apt install floss
```

#### Fedora/RHEL

```bash
echo '[floss]
name=Floss
baseurl=https://repo.nom-nom-hub.sh/yum/
enabled=1
gpgcheck=1
gpgkey=https://repo.nom-nom-hub.sh/yum/gpg.key' | sudo tee /etc/yum.repos.d/floss.repo
sudo yum install floss
```

### Manual Installation

Download pre-built binaries from the [releases page](https://github.com/nom-nom-hub/floss/releases) for:

- Linux
- macOS
- Windows
- FreeBSD
- OpenBSD
- NetBSD

### Installation with Go

If you have Go installed, you can install FLOSS directly:

```bash
go install github.com/nom-nom-hub/floss@latest
```

## Post-Installation Setup

After installing FLOSS, you'll need to configure it with your AI provider API keys. See the [Configuration Guide](configuration.md) for details.

## Verifying Installation

To verify that FLOSS is installed correctly, run:

```bash
floss --version
```

This should display the version of FLOSS that you installed.

## Troubleshooting

If you encounter issues during installation, see the [Troubleshooting Guide](troubleshooting.md) for solutions to common problems.