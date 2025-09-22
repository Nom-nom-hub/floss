package agent

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/nom-nom-hub/floss/internal/csync"
	"github.com/nom-nom-hub/floss/internal/llm/agent"
	"github.com/nom-nom-hub/floss/internal/message"
	"github.com/nom-nom-hub/floss/internal/session"
)

// WorkflowStatus represents the status of a workflow
type WorkflowStatus string

const (
	WorkflowStatusPending   WorkflowStatus = "pending"
	WorkflowStatusRunning   WorkflowStatus = "running"
	WorkflowStatusCompleted WorkflowStatus = "completed"
	WorkflowStatusFailed    WorkflowStatus = "failed"
	WorkflowStatusCancelled WorkflowStatus = "cancelled"
)

// WorkflowInstance represents an instance of a workflow execution
type WorkflowInstance struct {
	ID            string                 `json:"id"`
	WorkflowID    string                 `json:"workflow_id"`
	Name          string                 `json:"name"`
	Description   string                 `json:"description"`
	Status        WorkflowStatus         `json:"status"`
	CurrentStep   int                    `json:"current_step"`
	Steps         []WorkflowStepInstance `json:"steps"`
	CreatedAt     time.Time              `json:"created_at"`
	StartedAt     *time.Time             `json:"started_at,omitempty"`
	CompletedAt   *time.Time             `json:"completed_at,omitempty"`
	SessionID     string                 `json:"session_id"`
	Context       map[string]interface{} `json:"context,omitempty"`
}

// WorkflowStepInstance represents an instance of a workflow step execution
type WorkflowStepInstance struct {
	StepNumber      int          `json:"step_number"`
	ResponsibleAgent AgentRole   `json:"responsible_agent"`
	Action          string       `json:"action"`
	Status          WorkflowStatus `json:"status"`
	StartedAt       *time.Time   `json:"started_at,omitempty"`
	CompletedAt     *time.Time   `json:"completed_at,omitempty"`
	Result          string       `json:"result,omitempty"`
	Error           string       `json:"error,omitempty"`
	SessionID       string       `json:"session_id,omitempty"`
}

// AgentOrchestrator manages the execution of agent workflows
type AgentOrchestrator struct {
	agentServices   *csync.Map[AgentRole, agent.Service]
	sessionService  session.Service
	messageService  message.Service
	workflows       *csync.Map[string, AgentWorkflow]
	workflowInstances *csync.Map[string, WorkflowInstance]
	mutex           sync.RWMutex
}

// NewAgentOrchestrator creates a new agent orchestrator
func NewAgentOrchestrator(
	agentServices *csync.Map[AgentRole, agent.Service],
	sessionService session.Service,
	messageService message.Service,
) *AgentOrchestrator {
	return &AgentOrchestrator{
		agentServices:     agentServices,
		sessionService:    sessionService,
		messageService:    messageService,
		workflows:         csync.NewMap[string, AgentWorkflow](),
		workflowInstances: csync.NewMap[string, WorkflowInstance](),
	}
}

// RegisterWorkflow registers a workflow with the orchestrator
func (ao *AgentOrchestrator) RegisterWorkflow(workflow AgentWorkflow) {
	ao.workflows.Set(workflow.ID, workflow)
	slog.Info("Registered workflow", "workflow_id", workflow.ID, "name", workflow.Name)
}

// StartWorkflow starts a new workflow instance
func (ao *AgentOrchestrator) StartWorkflow(ctx context.Context, workflowID string, contextData map[string]interface{}) (*WorkflowInstance, error) {
	workflow, exists := ao.workflows.Get(workflowID)
	if !exists {
		return nil, fmt.Errorf("workflow %s not found", workflowID)
	}

	// Create a session for this workflow
	sessionTitle := fmt.Sprintf("Workflow: %s", workflow.Name)
	sess, err := ao.sessionService.Create(ctx, sessionTitle)
	if err != nil {
		return nil, fmt.Errorf("failed to create session for workflow: %w", err)
	}

	// Create workflow instance
	workflowInstance := WorkflowInstance{
		ID:          fmt.Sprintf("wf_%s_%d", workflowID, time.Now().Unix()),
		WorkflowID:  workflowID,
		Name:        workflow.Name,
		Description: workflow.Description,
		Status:      WorkflowStatusPending,
		Steps:       make([]WorkflowStepInstance, len(workflow.Steps)),
		CreatedAt:   time.Now(),
		SessionID:   sess.ID,
		Context:     contextData,
	}

	// Initialize step instances
	for i, step := range workflow.Steps {
		workflowInstance.Steps[i] = WorkflowStepInstance{
			StepNumber:      step.StepNumber,
			ResponsibleAgent: step.ResponsibleAgent,
			Action:          step.Action,
			Status:          WorkflowStatusPending,
		}
	}

	ao.workflowInstances.Set(workflowInstance.ID, workflowInstance)
	slog.Info("Started workflow instance", "instance_id", workflowInstance.ID, "workflow_id", workflowID)

	// Start the workflow execution asynchronously
	go ao.executeWorkflow(ctx, workflowInstance.ID)

	return &workflowInstance, nil
}

// executeWorkflow executes a workflow instance
func (ao *AgentOrchestrator) executeWorkflow(ctx context.Context, workflowInstanceID string) {
	workflowInstance, exists := ao.workflowInstances.Get(workflowInstanceID)
	if !exists {
		slog.Error("Workflow instance not found", "instance_id", workflowInstanceID)
		return
	}

	workflow, exists := ao.workflows.Get(workflowInstance.WorkflowID)
	if !exists {
		slog.Error("Workflow template not found", "workflow_id", workflowInstance.WorkflowID)
		return
	}

	// Update workflow status to running
	workflowInstance.Status = WorkflowStatusRunning
	now := time.Now()
	workflowInstance.StartedAt = &now
	ao.workflowInstances.Set(workflowInstanceID, workflowInstance)

	slog.Info("Executing workflow", "instance_id", workflowInstanceID, "workflow_id", workflowInstance.WorkflowID)

	// Execute steps in order
	for i := 0; i < len(workflow.Steps); i++ {
		step := workflow.Steps[i]
		stepInstance := &workflowInstance.Steps[i]

		// Check if step dependencies are satisfied
		if !ao.areDependenciesSatisfied(&workflowInstance, step) {
			// Find next step that can be executed
			nextStepIndex := ao.findNextExecutableStep(&workflowInstance, workflow)
			if nextStepIndex == -1 {
				// No more steps can be executed
				break
			}
			// Swap current step with next executable step
			workflowInstance.Steps[i], workflowInstance.Steps[nextStepIndex] = 
				workflowInstance.Steps[nextStepIndex], workflowInstance.Steps[i]
			step = workflow.Steps[i]
			stepInstance = &workflowInstance.Steps[i]
		}

		// Execute the step
		err := ao.executeStep(ctx, &workflowInstance, stepInstance, step)
		if err != nil {
			slog.Error("Failed to execute step", "error", err, "step", step.StepNumber)
			workflowInstance.Status = WorkflowStatusFailed
			now := time.Now()
			workflowInstance.CompletedAt = &now
			ao.workflowInstances.Set(workflowInstanceID, workflowInstance)
			return
		}

		// Update workflow instance
		ao.workflowInstances.Set(workflowInstanceID, workflowInstance)
	}

	// Mark workflow as completed
	workflowInstance.Status = WorkflowStatusCompleted
	now = time.Now()
	workflowInstance.CompletedAt = &now
	ao.workflowInstances.Set(workflowInstanceID, workflowInstance)

	slog.Info("Workflow completed", "instance_id", workflowInstanceID)
}

// areDependenciesSatisfied checks if all dependencies for a step are satisfied
func (ao *AgentOrchestrator) areDependenciesSatisfied(workflowInstance *WorkflowInstance, step WorkflowStep) bool {
	for _, dep := range step.Dependencies {
		// Find the dependent step
		depStepIndex := -1
		for i, s := range workflowInstance.Steps {
			if s.StepNumber == dep {
				depStepIndex = i
				break
			}
		}
		if depStepIndex == -1 {
			// Dependency not found
			return false
		}
		// Check if dependency is completed
		if workflowInstance.Steps[depStepIndex].Status != WorkflowStatusCompleted {
			return false
		}
	}
	return true
}

// findNextExecutableStep finds the next step that can be executed
func (ao *AgentOrchestrator) findNextExecutableStep(workflowInstance *WorkflowInstance, workflow AgentWorkflow) int {
	for i := workflowInstance.CurrentStep + 1; i < len(workflow.Steps); i++ {
		step := workflow.Steps[i]
		if ao.areDependenciesSatisfied(workflowInstance, step) {
			return i
		}
	}
	return -1
}

// executeStep executes a single workflow step
func (ao *AgentOrchestrator) executeStep(ctx context.Context, workflowInstance *WorkflowInstance, stepInstance *WorkflowStepInstance, step WorkflowStep) error {
	slog.Info("Executing step", "step", step.StepNumber, "agent", step.ResponsibleAgent)

	// Update step status to running
	stepInstance.Status = WorkflowStatusRunning
	now := time.Now()
	stepInstance.StartedAt = &now

	// Get the agent service for this step
	agentService, exists := ao.agentServices.Get(step.ResponsibleAgent)
	if !exists {
		stepInstance.Status = WorkflowStatusFailed
		stepInstance.Error = fmt.Sprintf("Agent service for role %s not found", step.ResponsibleAgent)
		now := time.Now()
		stepInstance.CompletedAt = &now
		return fmt.Errorf("agent service for role %s not found", step.ResponsibleAgent)
	}

	// Create a session for this step if one doesn't exist
	var sessionID string
	if stepInstance.SessionID == "" {
		sessionTitle := fmt.Sprintf("Workflow Step: %s - %s", workflowInstance.Name, step.Action)
		sess, err := ao.sessionService.Create(ctx, sessionTitle)
		if err != nil {
			stepInstance.Status = WorkflowStatusFailed
			stepInstance.Error = fmt.Sprintf("Failed to create session: %v", err)
			now := time.Now()
			stepInstance.CompletedAt = &now
			return fmt.Errorf("failed to create session: %w", err)
		}
		sessionID = sess.ID
		stepInstance.SessionID = sessionID
	} else {
		sessionID = stepInstance.SessionID
	}

	// Prepare the prompt for the agent
	prompt := fmt.Sprintf("Workflow Step: %s\n\nAction: %s\n\nContext: %+v\n\nPlease complete this step and provide a detailed report of your work.", 
		workflowInstance.Name, step.Action, workflowInstance.Context)

	// Run the agent
	events, err := agentService.Run(ctx, sessionID, prompt)
	if err != nil {
		stepInstance.Status = WorkflowStatusFailed
		stepInstance.Error = fmt.Sprintf("Failed to run agent: %v", err)
		now := time.Now()
		stepInstance.CompletedAt = &now
		return fmt.Errorf("failed to run agent: %w", err)
	}

	// Wait for the agent to complete
	var finalEvent agent.AgentEvent
	for event := range events {
		finalEvent = event
	}

	// Check if the agent completed successfully
	if finalEvent.Error != nil {
		stepInstance.Status = WorkflowStatusFailed
		stepInstance.Error = finalEvent.Error.Error()
		now := time.Now()
		stepInstance.CompletedAt = &now
		return finalEvent.Error
	}

	// Get the result from the messages
	messages, err := ao.messageService.List(ctx, sessionID)
	if err != nil {
		stepInstance.Status = WorkflowStatusFailed
		stepInstance.Error = fmt.Sprintf("Failed to get messages: %v", err)
		now := time.Now()
		stepInstance.CompletedAt = &now
		return fmt.Errorf("failed to get messages: %w", err)
	}

	// Extract the result from the last assistant message
	var result string
	for i := len(messages) - 1; i >= 0; i-- {
		if messages[i].Role == message.Assistant {
			result = messages[i].Content().String()
			break
		}
	}

	// Update step instance
	stepInstance.Status = WorkflowStatusCompleted
	stepInstance.Result = result
	now = time.Now()
	stepInstance.CompletedAt = &now

	slog.Info("Step completed", "step", step.StepNumber, "agent", step.ResponsibleAgent)
	return nil
}

// GetWorkflowInstance retrieves a workflow instance by ID
func (ao *AgentOrchestrator) GetWorkflowInstance(instanceID string) (*WorkflowInstance, bool) {
	workflowInstance, exists := ao.workflowInstances.Get(instanceID)
	if !exists {
		return nil, false
	}
	return &workflowInstance, true
}

// ListWorkflowInstances lists all workflow instances
func (ao *AgentOrchestrator) ListWorkflowInstances() []WorkflowInstance {
	var instances []WorkflowInstance
	for _, instance := range ao.workflowInstances.Seq2() {
		instances = append(instances, instance)
	}
	return instances
}

// CancelWorkflow cancels a running workflow instance
func (ao *AgentOrchestrator) CancelWorkflow(instanceID string) error {
	workflowInstance, exists := ao.workflowInstances.Get(instanceID)
	if !exists {
		return fmt.Errorf("workflow instance %s not found", instanceID)
	}

	// Cancel all agent services involved in the workflow
	for _, step := range workflowInstance.Steps {
		if step.Status == WorkflowStatusRunning {
			agentService, exists := ao.agentServices.Get(step.ResponsibleAgent)
			if exists {
				agentService.Cancel(step.SessionID)
			}
		}
	}

	// Update workflow status
	workflowInstance.Status = WorkflowStatusCancelled
	now := time.Now()
	workflowInstance.CompletedAt = &now
	ao.workflowInstances.Set(instanceID, workflowInstance)

	slog.Info("Cancelled workflow", "instance_id", instanceID)
	return nil
}

// GetAgentService retrieves an agent service by role
func (ao *AgentOrchestrator) GetAgentService(role AgentRole) (agent.Service, bool) {
	return ao.agentServices.Get(role)
}

// RegisterAgentService registers an agent service with the orchestrator
func (ao *AgentOrchestrator) RegisterAgentService(role AgentRole, service agent.Service) {
	ao.agentServices.Set(role, service)
	slog.Info("Registered agent service", "role", role)
}