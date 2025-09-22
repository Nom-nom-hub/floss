# Virtual Company Agent System

The FLOSS Agent System simulates a complete software development company with specialized AI agents filling various roles. This system enables complex workflows and collaboration between different specialized agents to simulate real-world software development processes.

## Overview

The agent system consists of:

1. **Organizational Structure** - A complete company hierarchy with departments and teams
2. **Specialized Agents** - AI agents with specific roles and capabilities
3. **Communication Protocols** - Standardized ways for agents to communicate
4. **Collaboration Protocols** - Multi-agent processes for complex tasks
5. **Workflows** - Sequences of steps to accomplish complex goals
6. **Orchestration System** - Management of agent interactions and workflow execution

## Agent Roles

The system includes the following specialized agent roles:

### Executive Leadership
- **CEO** - Sets strategic direction and makes high-level decisions
- **CTO** - Defines technical architecture and standards

### Product & Project Management
- **Product Manager** - Defines product vision and requirements
- **Project Manager** - Coordinates projects and timelines
- **Business Analyst** - Analyzes requirements and processes

### Engineering
- **Tech Lead** - Provides technical guidance and code reviews
- **Senior Developer** - Implements complex features and mentors juniors
- **Junior Developer** - Implements features and learns from seniors

### Quality & Security
- **QA Engineer** - Ensures software quality through testing
- **Security Engineer** - Audits security and identifies vulnerabilities

### Infrastructure & Operations
- **DevOps Engineer** - Manages deployment and infrastructure
- **Database Administrator** - Manages databases and optimization

### Design & Documentation
- **UI/UX Designer** - Creates user interfaces and experiences
- **Technical Writer** - Produces technical documentation

## Workflows

The system includes predefined workflows for common processes:

### Product Development Process
A complete workflow from concept to deployment including:
1. Requirements definition
2. Technical architecture
3. Implementation
4. Testing and QA
5. Security auditing
6. Deployment to staging and production
7. Documentation updates

### Incident Response Process
A workflow for handling system incidents:
1. Incident identification and assessment
2. Response team activation
3. Immediate fixes
4. Root cause analysis
5. Permanent fix development and implementation
6. Testing and deployment
7. Post-incident review

## Configuration

The agent system is configured through the `agents.json` file in the FLOSS configuration directory. This file defines:

- Company structure and hierarchy
- Agent definitions with roles and capabilities
- Communication protocols
- Collaboration protocols
- Workflows

## Commands

The agent system can be managed through the FLOSS CLI:

### List Available Agent Roles
```bash
floss agents list
```

### Get Information About a Specific Agent Role
```bash
floss agents info [role]
```

### List Available Workflows
```bash
floss agents workflows
```

### Initialize Agent System with Default Configuration
```bash
floss agents init
```

### Show Current Agent System Configuration
```bash
floss agents config
```

### Enable/Disable Agent System
```bash
floss agents enable
floss agents disable
```

## Usage Examples

### Starting a Product Development Workflow
```bash
# Start a new product development workflow
floss run "Please start a product development workflow for a new user authentication feature"
```

### Requesting a Security Audit
```bash
# Request a security audit through the security engineer agent
floss run "Please conduct a security audit of our current codebase"
```

### Getting a Progress Update on a Workflow
```bash
# Check the status of a running workflow
floss run "What is the current status of our product development workflow?"
```

## Customization

The agent system can be customized by modifying the `agents.json` configuration file:

1. **Add new agent roles** by defining new entries in the `agent_definitions` section
2. **Modify existing agent capabilities** by changing their allowed tools and prompt templates
3. **Create new workflows** by adding entries to the `workflows` section
4. **Define new communication protocols** in the `communication_protocols` section
5. **Establish new collaboration protocols** in the `collaboration_protocols` section

## Integration with FLOSS Core

The agent system integrates with the core FLOSS functionality:

- **LLM Providers** - Uses the same provider configuration as the main FLOSS system
- **Tools** - Agents have access to the same tools as the main FLOSS system
- **Sessions** - Each agent interaction and workflow step creates its own session
- **Messages** - All agent communications are stored as messages in the database
- **Permissions** - Agents respect the same permission system as the main FLOSS system

## Best Practices

1. **Start with default configuration** - Use `floss agents init` to get started with a complete default setup
2. **Enable only needed agents** - Disable agents you don't need to reduce resource usage
3. **Monitor workflow progress** - Use the CLI to check the status of running workflows
4. **Customize prompts for your domain** - Modify agent prompt templates to better suit your specific use case
5. **Use appropriate models** - Assign larger models to complex roles like CEO and CTO, smaller models to simpler roles