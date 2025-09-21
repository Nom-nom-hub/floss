# How to Use the Goal Driven Development Agents

This guide explains how to use the Goal Driven Development agents with Qwen Code.

## Prerequisites

1. Make sure you have Qwen Code installed:
   ```bash
   npm install -g @qwen-code/qwen-code@latest
   ```

2. Ensure the `.qwen` directory with agents is in your project root

## Using the Agents

### 1. Start Qwen Code
```bash
cd your-project-directory
qwen
```

### 2. Create a New Project for Testing
```bash
mkdir test-project
cd test-project
```

### 3. Copy the Goal Driven Development Agents
Copy the `.qwen` directory from the qwen-spec project to your test project:
```bash
cp -r ../qwen-spec/.qwen .
```

### 4. Start Qwen Code
```bash
qwen
```

### 5. Use the Agents

#### Method 1: Interactive Usage
Once Qwen Code is running, you can use the agents in several ways:

1. **Ask the CEO to set a goal:**
   ```
   > I want to create an offline notes application
   ```

2. **Directly invoke a specific agent:**
   ```
   > Use the developer agent to create a simple HTML file with a notes UI
   ```

3. **Use a workflow:**
   ```
   > Follow the api_development workflow to create a notes API
   ```

#### Method 2: Command Line Usage
You can also use the agents from the command line:

```bash
# Use the CEO agent to set a goal
echo "Create a REST API for a todo application" | qwen --prompt "Act as the CEO and break this down into tasks"

# Use the developer agent to implement a feature
echo "Create a simple Express server with a /todos endpoint" | qwen --prompt "Act as the developer and implement this feature"
```

## Available Agents

### Core Agents
1. **ceo** - Sets project goals and makes high-level decisions
2. **ceo_advisor** - Provides technical guidance to the CEO
3. **orchestrator** - Breaks down goals into tasks and coordinates the team
4. **tech_lead** - Reviews code quality and architecture
5. **developer** - Implements features and functionality
6. **qa** - Tests features and validates correctness
7. **security_engineer** - Reviews code for security issues
8. **devops** - Handles deployment and infrastructure
9. **documentation** - Updates documentation

### Utility Agents
1. **git_workflow_enforcer** - Ensures Git workflow compliance
2. **agent_communication_system** - Manages communication between agents
3. **audit_logger** - Maintains audit logs

### Workflows
1. **user_login_mfa** - Implements user authentication with MFA
2. **api_development** - Creates REST APIs with full testing
3. **ml_model_deployment** - Develops and deploys ML models

## Example Usage

### Creating an Offline Notes App
```bash
mkdir notes-app
cd notes-app
qwen
```

Then in the Qwen interface:
```
> Act as the CEO and create an offline notes application that works without internet connectivity
```

The system will automatically:
1. Have the CEO set the goal
2. Have the CEO_Advisor provide technical recommendations
3. Have the Orchestrator break down the tasks
4. Have Developers implement the features
5. Have QA test the functionality
6. Have Security_Engineer review for vulnerabilities
7. Have Tech_Lead review the architecture
8. Have Documentation update the docs
9. Have DevOps prepare for deployment

### Using a Specific Workflow
```bash
qwen
> Follow the user_login_mfa workflow to implement authentication
```

This will automatically coordinate all agents to implement a complete user authentication system with multi-factor authentication.

## Best Practices

1. **Start with the CEO**: Always start by setting a clear goal with the CEO agent
2. **Let the system coordinate**: Allow the agents to automatically coordinate with each other
3. **Monitor progress**: Watch the real-time updates as agents work
4. **Provide feedback**: Give guidance when needed, especially for business decisions
5. **Review deliverables**: Check the final output before accepting it

## Troubleshooting

### Agents Not Found
If Qwen Code doesn't recognize the agents:
1. Ensure the `.qwen` directory is in your project root
2. Verify that agent files have the correct format (YAML frontmatter + content)
3. Restart Qwen Code to refresh the agent list

### Agent Not Responding
If an agent isn't responding as expected:
1. Check that the agent file has proper YAML frontmatter
2. Ensure the agent has the necessary tools listed
3. Verify that the system prompt is clear and actionable

### Git Workflow Issues
If there are Git workflow compliance issues:
1. Ensure the git_workflow_enforcer agent is available
2. Check that branch naming and commit message conventions are followed
3. Verify that direct pushes to remote repositories are prevented

## Customization

You can customize the agents by editing their Markdown files:
1. Modify the system prompt to change agent behavior
2. Add or remove tools from the tools list
3. Adjust the agent descriptions for better discoverability

For example, to make the Developer agent more focused on React development:
1. Edit `.qwen/agents/developer.md`
2. Update the system prompt to emphasize React patterns
3. Add React-specific tools if needed