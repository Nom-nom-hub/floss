# Goal Driven Development - Agent System Configuration

This directory contains the configuration files for the Goal Driven Development agent system.

## Directory Structure

```
.qwen/
├── agents/
│   ├── ceo.md
│   ├── ceo_advisor.md
│   ├── orchestrator.md
│   ├── tech_lead.md
│   ├── developer.md
│   ├── qa.md
│   ├── security_engineer.md
│   ├── devops.md
│   ├── documentation.md
│   ├── git_workflow_enforcer.md
│   ├── agent_communication_system.md
│   └── audit_logger.md
├── workflows/
│   ├── user_login_mfa.md
│   ├── api_development.md
│   └── ml_model_deployment.md
├── utils/
│   ├── git-workflow-enforcer.js
│   ├── agent-communication-system.js
│   └── audit-logger.js
├── config.json
└── README.md
```

## Agents

The following agents are defined in this system:

1. **CEO** (`ceo`) - Sets project goals, approves final merges and pushes
2. **CEO_Advisor** (`ceo_advisor`) - Guides CEO on prioritization and technical decisions
3. **Orchestrator** (`orchestrator`) - Breaks goals into tasks, assigns agents, manages Git workflow
4. **Tech_Lead** (`tech_lead`) - Reviews PRs, ensures architecture and code quality
5. **Developer** (`developer`) - Implements features, commits locally, drafts PRs
6. **QA** (`qa`) - Tests features, validates against goals
7. **Security_Engineer** (`security_engineer`) - Reviews code for security and best practices
8. **DevOps** (`devops`) - Deploys code after approvals
9. **Documentation** (`documentation`) - Updates docs after merges

## Utilities

- **GitWorkflowEnforcer** (`git_workflow_enforcer`) - Enforces Git workflow rules
- **AgentCommunicationSystem** (`agent_communication_system`) - Manages communication between agents
- **AuditLogger** (`audit_logger`) - Maintains audit logs for compliance

## Workflows

Example workflows demonstrating the system capabilities:

1. **UserLoginMFAWorkflow** (`user_login_mfa`) - User login with multi-factor authentication
2. **APIDevelopmentWorkflow** (`api_development`) - REST API development with full testing
3. **MLModelDeploymentWorkflow** (`ml_model_deployment`) - Machine learning model development and deployment

## Usage

To use these agents with Qwen Code, ensure this directory is in your project's `.qwen` folder. The agents will be automatically discovered and available for use.