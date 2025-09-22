package agent

import (
	"github.com/nom-nom-hub/floss/internal/csync"
)

// CommunicationChannel represents a channel for agent communication
type CommunicationChannel string

const (
	ChannelDirectMessage CommunicationChannel = "direct_message"
	ChannelTeamChannel   CommunicationChannel = "team_channel"
	ChannelDepartment    CommunicationChannel = "department_channel"
	ChannelCompanyWide   CommunicationChannel = "company_wide"
)

// AgentCommunication defines how agents can communicate with each other
type AgentCommunication struct {
	Channel     CommunicationChannel `json:"channel"`
	Participants []AgentRole         `json:"participants,omitempty"`
	Subject     string              `json:"subject,omitempty"`
	Message     string              `json:"message"`
	Priority    int                 `json:"priority,omitempty"` // 1-10, 10 being highest priority
}

// CollaborationProtocol defines how agents work together on tasks
type CollaborationProtocol struct {
	InitiatingAgent AgentRole           `json:"initiating_agent"`
	TargetAgents    []AgentRole         `json:"target_agents"`
	Task            string              `json:"task"`
	Context         string              `json:"context,omitempty"`
	ExpectedOutcome string              `json:"expected_outcome"`
	Deadline        string              `json:"deadline,omitempty"` // ISO 8601 format
	Priority        int                 `json:"priority,omitempty"` // 1-10, 10 being highest priority
}

// AgentWorkflow represents a sequence of agent interactions to complete a complex task
type AgentWorkflow struct {
	ID          string                `json:"id"`
	Name        string                `json:"name"`
	Description string                `json:"description"`
	Steps       []WorkflowStep        `json:"steps"`
}

// WorkflowStep represents a single step in an agent workflow
type WorkflowStep struct {
	StepNumber      int             `json:"step_number"`
	ResponsibleAgent AgentRole      `json:"responsible_agent"`
	Action          string          `json:"action"`
	Input           string          `json:"input,omitempty"`
	ExpectedOutput  string          `json:"expected_output"`
	Dependencies    []int           `json:"dependencies,omitempty"` // Step numbers this step depends on
	NextSteps       []int           `json:"next_steps,omitempty"`   // Step numbers that follow this step
}

// DefaultCommunicationProtocols returns the default communication protocols for agent interactions
func DefaultCommunicationProtocols() map[string]AgentCommunication {
	return map[string]AgentCommunication{
		"standup": {
			Channel:     ChannelTeamChannel,
			Subject:     "Daily Standup",
			Message:     "Please provide your daily standup update: What did you do yesterday? What will you do today? Any blockers?",
			Priority:    5,
		},
		"code_review": {
			Channel:     ChannelDirectMessage,
			Subject:     "Code Review Request",
			Message:     "I've submitted a pull request that requires your review. Please review the code and provide feedback.",
			Priority:    7,
		},
		"bug_report": {
			Channel:     ChannelTeamChannel,
			Subject:     "Bug Report",
			Message:     "A bug has been identified. Please investigate and provide a fix.",
			Priority:    8,
		},
		"feature_request": {
			Channel:     ChannelDepartment,
			Subject:     "Feature Request",
			Message:     "A new feature has been requested. Please evaluate and provide implementation plan.",
			Priority:    6,
		},
		"security_alert": {
			Channel:     ChannelCompanyWide,
			Subject:     "Security Alert",
			Message:     "A security vulnerability has been identified. Immediate attention required.",
			Priority:    10,
		},
		"performance_issue": {
			Channel:     ChannelDirectMessage,
			Subject:     "Performance Issue",
			Message:     "A performance issue has been identified. Please investigate and optimize.",
			Priority:    7,
		},
	}
}

// DefaultCollaborationProtocols returns the default collaboration protocols for multi-agent tasks
func DefaultCollaborationProtocols() map[string]CollaborationProtocol {
	return map[string]CollaborationProtocol{
		"new_project": {
			InitiatingAgent: AgentRoleProductManager,
			TargetAgents:    []AgentRole{AgentRoleCTO, AgentRoleTechLead, AgentRoleProjectManager},
			Task:            "Plan and initiate a new project",
			Context:         "A new project has been approved by the CEO. We need to plan the technical implementation and project timeline.",
			ExpectedOutcome: "Technical architecture document, project plan, and resource allocation plan",
			Priority:        9,
		},
		"feature_implementation": {
			InitiatingAgent: AgentRoleTechLead,
			TargetAgents:    []AgentRole{AgentRoleSeniorDeveloper, AgentRoleJuniorDeveloper},
			Task:            "Implement a new feature",
			Context:         "A new feature has been prioritized in the product roadmap. Implementation is needed.",
			ExpectedOutcome: "Implemented feature with tests and documentation",
			Priority:        7,
		},
		"security_audit": {
			InitiatingAgent: AgentRoleSecurityEngineer,
			TargetAgents:    []AgentRole{AgentRoleCTO, AgentRoleDevOpsEngineer},
			Task:            "Conduct a security audit",
			Context:         "Regular security audit is required to ensure system security and compliance.",
			ExpectedOutcome: "Security audit report with identified vulnerabilities and recommendations",
			Priority:        8,
		},
		"performance_optimization": {
			InitiatingAgent: AgentRoleDevOpsEngineer,
			TargetAgents:    []AgentRole{AgentRoleSeniorDeveloper, AgentRoleDatabaseAdmin},
			Task:            "Optimize system performance",
			Context:         "Performance metrics indicate the system is not meeting performance targets.",
			ExpectedOutcome: "Performance optimization plan and implemented improvements",
			Priority:        7,
		},
		"user_experience_review": {
			InitiatingAgent: AgentRoleUIDesigner,
			TargetAgents:    []AgentRole{AgentRoleProductManager, AgentRoleSeniorDeveloper},
			Task:            "Review user experience",
			Context:         "User feedback indicates potential improvements to the user experience.",
			ExpectedOutcome: "UX improvement recommendations and implementation plan",
			Priority:        6,
		},
	}
}

// DefaultWorkflows returns the default workflows for common multi-step processes
func DefaultWorkflows() map[string]AgentWorkflow {
	return map[string]AgentWorkflow{
		"product_development": {
			ID:          "product_development",
			Name:        "Product Development Process",
			Description: "Complete workflow for developing a new product feature from concept to deployment",
			Steps: []WorkflowStep{
				{
					StepNumber:      1,
					ResponsibleAgent: AgentRoleProductManager,
					Action:          "Define requirements",
					ExpectedOutput:  "Detailed product requirements document",
				},
				{
					StepNumber:      2,
					ResponsibleAgent: AgentRoleBusinessAnalyst,
					Action:          "Analyze requirements",
					ExpectedOutput:  "Requirements analysis report",
					Dependencies:    []int{1},
				},
				{
					StepNumber:      3,
					ResponsibleAgent: AgentRoleCTO,
					Action:          "Architect solution",
					ExpectedOutput:  "Technical architecture document",
					Dependencies:    []int{2},
				},
				{
					StepNumber:      4,
					ResponsibleAgent: AgentRoleTechLead,
					Action:          "Create implementation plan",
					ExpectedOutput:  "Detailed implementation plan with timeline",
					Dependencies:    []int{3},
				},
				{
					StepNumber:      5,
					ResponsibleAgent: AgentRoleProjectManager,
					Action:          "Coordinate resources",
					ExpectedOutput:  "Resource allocation plan and team assignments",
					Dependencies:    []int{4},
				},
				{
					StepNumber:      6,
					ResponsibleAgent: AgentRoleSeniorDeveloper,
					Action:          "Implement core features",
					ExpectedOutput:  "Implemented core features with unit tests",
					Dependencies:    []int{5},
				},
				{
					StepNumber:      7,
					ResponsibleAgent: AgentRoleJuniorDeveloper,
					Action:          "Implement supporting features",
					ExpectedOutput:  "Implemented supporting features with unit tests",
					Dependencies:    []int{5},
				},
				{
					StepNumber:      8,
					ResponsibleAgent: AgentRoleQAEngineer,
					Action:          "Test implementation",
					ExpectedOutput:  "Test reports and bug reports",
					Dependencies:    []int{6, 7},
				},
				{
					StepNumber:      9,
					ResponsibleAgent: AgentRoleSecurityEngineer,
					Action:          "Security audit",
					ExpectedOutput:  "Security audit report",
					Dependencies:    []int{6, 7},
				},
				{
					StepNumber:      10,
					ResponsibleAgent: AgentRoleTechLead,
					Action:          "Code review and integration",
					ExpectedOutput:  "Integrated and reviewed codebase",
					Dependencies:    []int{8, 9},
				},
				{
					StepNumber:      11,
					ResponsibleAgent: AgentRoleDevOpsEngineer,
					Action:          "Deploy to staging",
					ExpectedOutput:  "Deployed to staging environment",
					Dependencies:    []int{10},
				},
				{
					StepNumber:      12,
					ResponsibleAgent: AgentRoleQAEngineer,
					Action:          "Staging validation",
					ExpectedOutput:  "Staging validation report",
					Dependencies:    []int{11},
				},
				{
					StepNumber:      13,
					ResponsibleAgent: AgentRoleProductManager,
					Action:          "User acceptance testing",
					ExpectedOutput:  "UAT sign-off",
					Dependencies:    []int{12},
				},
				{
					StepNumber:      14,
					ResponsibleAgent: AgentRoleDevOpsEngineer,
					Action:          "Deploy to production",
					ExpectedOutput:  "Deployed to production environment",
					Dependencies:    []int{13},
				},
				{
					StepNumber:      15,
					ResponsibleAgent: AgentRoleTechnicalWriter,
					Action:          "Update documentation",
					ExpectedOutput:  "Updated user and technical documentation",
					Dependencies:    []int{14},
				},
			},
		},
		"incident_response": {
			ID:          "incident_response",
			Name:        "Incident Response Process",
			Description: "Workflow for responding to system incidents and outages",
			Steps: []WorkflowStep{
				{
					StepNumber:      1,
					ResponsibleAgent: AgentRoleDevOpsEngineer,
					Action:          "Identify and assess incident",
					ExpectedOutput:  "Incident report with severity classification",
				},
				{
					StepNumber:      2,
					ResponsibleAgent: AgentRoleCTO,
					Action:          "Activate incident response",
					ExpectedOutput:  "Incident response team assembled",
					Dependencies:    []int{1},
				},
				{
					StepNumber:      3,
					ResponsibleAgent: AgentRoleDevOpsEngineer,
					Action:          "Implement immediate fixes",
					ExpectedOutput:  "System stabilized",
					Dependencies:    []int{2},
				},
				{
					StepNumber:      4,
					ResponsibleAgent: AgentRoleSecurityEngineer,
					Action:          "Security assessment",
					ExpectedOutput:  "Security impact assessment",
					Dependencies:    []int{3},
				},
				{
					StepNumber:      5,
					ResponsibleAgent: AgentRoleSeniorDeveloper,
					Action:          "Root cause analysis",
					ExpectedOutput:  "Root cause analysis report",
					Dependencies:    []int{3},
				},
				{
					StepNumber:      6,
					ResponsibleAgent: AgentRoleTechLead,
					Action:          "Develop permanent fix",
					ExpectedOutput:  "Permanent fix implementation plan",
					Dependencies:    []int{5},
				},
				{
					StepNumber:      7,
					ResponsibleAgent: AgentRoleSeniorDeveloper,
					Action:          "Implement permanent fix",
					ExpectedOutput:  "Implemented permanent fix",
					Dependencies:    []int{6},
				},
				{
					StepNumber:      8,
					ResponsibleAgent: AgentRoleQAEngineer,
					Action:          "Test fix",
					ExpectedOutput:  "Test results confirming fix",
					Dependencies:    []int{7},
				},
				{
					StepNumber:      9,
					ResponsibleAgent: AgentRoleDevOpsEngineer,
					Action:          "Deploy fix",
					ExpectedOutput:  "Fix deployed to production",
					Dependencies:    []int{8},
				},
				{
					StepNumber:      10,
					ResponsibleAgent: AgentRoleCTO,
					Action:          "Post-incident review",
					ExpectedOutput:  "Post-incident review report with action items",
					Dependencies:    []int{9},
				},
			},
		},
	}
}

// AgentCollaborationSystem manages agent interactions and workflows
type AgentCollaborationSystem struct {
	Agents          *csync.Map[string, AgentDefinition]     `json:"-"`
	Communications  *csync.Map[string, AgentCommunication]  `json:"-"`
	Collaborations  *csync.Map[string, CollaborationProtocol] `json:"-"`
	Workflows       *csync.Map[string, AgentWorkflow]       `json:"-"`
}

// NewAgentCollaborationSystem creates a new agent collaboration system
func NewAgentCollaborationSystem() *AgentCollaborationSystem {
	return &AgentCollaborationSystem{
		Agents:         csync.NewMap[string, AgentDefinition](),
		Communications: csync.NewMap[string, AgentCommunication](),
		Collaborations: csync.NewMap[string, CollaborationProtocol](),
		Workflows:      csync.NewMap[string, AgentWorkflow](),
	}
}