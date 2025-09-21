/**
 * @file Agent Communication System
 * @description Manages communication channels between all agents in the system
 * @license Apache-2.0
 */

import { writeFileSync, readFileSync, existsSync, mkdirSync } from 'fs';
import { join } from 'path';

// Ensure the communications directory exists
const COMMUNICATIONS_DIR = join('.qwen', 'communications');
if (!existsSync(COMMUNICATIONS_DIR)) {
  mkdirSync(COMMUNICATIONS_DIR, { recursive: true });
}

// Message types
const MESSAGE_TYPES = {
  DIRECT: 'direct',
  BROADCAST: 'broadcast',
  STATUS_UPDATE: 'status_update',
  TASK_NOTIFICATION: 'task_notification',
  APPROVAL_REQUEST: 'approval_request',
  EMERGENCY_ALERT: 'emergency_alert'
};

// Priority levels
const PRIORITY_LEVELS = {
  LOW: 'low',
  NORMAL: 'normal',
  HIGH: 'high',
  URGENT: 'urgent'
};

/**
 * Message structure
 * @typedef {Object} Message
 * @property {string} id - Unique message identifier
 * @property {string} sender - Sender agent name
 * @property {string|string[]} recipients - Recipient agent name(s)
 * @property {string} type - Message type
 * @property {string} content - Message content
 * @property {string} timestamp - ISO timestamp
 * @property {string} priority - Priority level
 * @property {boolean} responseRequired - Whether a response is required
 * @property {string} [status] - Message status (pending, delivered, read, failed)
 * @property {string} [parentId] - Parent message ID for threaded conversations
 */

/**
 * Send a message between agents
 * @param {string} sender - Sender agent name
 * @param {string|string[]} recipients - Recipient agent name(s)
 * @param {string} type - Message type
 * @param {string} content - Message content
 * @param {Object} [options] - Additional options
 * @param {string} [options.priority] - Priority level
 * @param {boolean} [options.responseRequired] - Whether a response is required
 * @param {string} [options.parentId] - Parent message ID for threaded conversations
 * @returns {string} Message ID
 */
export function sendMessage(sender, recipients, type, content, options = {}) {
  const messageId = generateMessageId();
  const timestamp = new Date().toISOString();
  
  const message = {
    id: messageId,
    sender,
    recipients,
    type,
    content,
    timestamp,
    priority: options.priority || PRIORITY_LEVELS.NORMAL,
    responseRequired: options.responseRequired || false,
    status: 'pending',
    parentId: options.parentId
  };
  
  // Save message to file
  const messageFile = join(COMMUNICATIONS_DIR, `${messageId}.json`);
  writeFileSync(messageFile, JSON.stringify(message, null, 2));
  
  // Update sender's sent messages
  updateAgentMessages(sender, messageId, 'sent');
  
  // Update recipient(s) inbox
  const recipientList = Array.isArray(recipients) ? recipients : [recipients];
  recipientList.forEach(recipient => {
    updateAgentMessages(recipient, messageId, 'inbox');
  });
  
  return messageId;
}

/**
 * Retrieve messages for an agent
 * @param {string} agentName - Agent name
 * @param {string} [folder] - Folder to retrieve from (inbox, sent, archived)
 * @returns {Message[]} Array of messages
 */
export function getMessages(agentName, folder = 'inbox') {
  const agentMessagesFile = join(COMMUNICATIONS_DIR, `${agentName}_messages.json`);
  
  if (!existsSync(agentMessagesFile)) {
    return [];
  }
  
  const agentMessages = JSON.parse(readFileSync(agentMessagesFile, 'utf8'));
  const messages = [];
  
  if (agentMessages[folder]) {
    agentMessages[folder].forEach(messageId => {
      const messageFile = join(COMMUNICATIONS_DIR, `${messageId}.json`);
      if (existsSync(messageFile)) {
        const message = JSON.parse(readFileSync(messageFile, 'utf8'));
        messages.push(message);
      }
    });
  }
  
  return messages;
}

/**
 * Mark a message as read
 * @param {string} agentName - Agent name
 * @param {string} messageId - Message ID
 */
export function markMessageAsRead(agentName, messageId) {
  const messageFile = join(COMMUNICATIONS_DIR, `${messageId}.json`);
  
  if (existsSync(messageFile)) {
    const message = JSON.parse(readFileSync(messageFile, 'utf8'));
    message.status = 'read';
    writeFileSync(messageFile, JSON.stringify(message, null, 2));
  }
}

/**
 * Archive a message
 * @param {string} agentName - Agent name
 * @param {string} messageId - Message ID
 */
export function archiveMessage(agentName, messageId) {
  const agentMessagesFile = join(COMMUNICATIONS_DIR, `${agentName}_messages.json`);
  
  if (existsSync(agentMessagesFile)) {
    const agentMessages = JSON.parse(readFileSync(agentMessagesFile, 'utf8'));
    
    // Remove from inbox
    if (agentMessages.inbox) {
      agentMessages.inbox = agentMessages.inbox.filter(id => id !== messageId);
    }
    
    // Add to archived
    if (!agentMessages.archived) {
      agentMessages.archived = [];
    }
    agentMessages.archived.push(messageId);
    
    writeFileSync(agentMessagesFile, JSON.stringify(agentMessages, null, 2));
  }
}

/**
 * Update agent's message tracking
 * @param {string} agentName - Agent name
 * @param {string} messageId - Message ID
 * @param {string} folder - Folder to add message to
 */
function updateAgentMessages(agentName, messageId, folder) {
  const agentMessagesFile = join(COMMUNICATIONS_DIR, `${agentName}_messages.json`);
  let agentMessages = {};
  
  if (existsSync(agentMessagesFile)) {
    agentMessages = JSON.parse(readFileSync(agentMessagesFile, 'utf8'));
  }
  
  if (!agentMessages[folder]) {
    agentMessages[folder] = [];
  }
  
  // Avoid duplicates
  if (!agentMessages[folder].includes(messageId)) {
    agentMessages[folder].push(messageId);
    writeFileSync(agentMessagesFile, JSON.stringify(agentMessages, null, 2));
  }
}

/**
 * Generate a unique message ID
 * @returns {string} Unique message ID
 */
function generateMessageId() {
  return `msg_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`;
}

/**
 * Send a status update from an agent
 * @param {string} agentName - Agent name
 * @param {string} status - Status message
 * @param {Object} [details] - Additional status details
 */
export function sendStatusUpdate(agentName, status, details = {}) {
  const content = {
    status,
    ...details
  };
  
  // Broadcast to orchestrator and CEO
  sendMessage(
    agentName,
    ['Orchestrator', 'CEO'],
    MESSAGE_TYPES.STATUS_UPDATE,
    JSON.stringify(content),
    { priority: PRIORITY_LEVELS.NORMAL }
  );
}

/**
 * Send an emergency alert
 * @param {string} agentName - Agent name
 * @param {string} alert - Alert message
 * @param {Object} [details] - Additional alert details
 */
export function sendEmergencyAlert(agentName, alert, details = {}) {
  const content = {
    alert,
    ...details
  };
  
  // Broadcast to all agents
  const allAgents = [
    'CEO', 'CEO_Advisor', 'Orchestrator', 'Tech_Lead',
    'Developer', 'QA', 'Security_Engineer', 'DevOps', 'Documentation'
  ];
  
  sendMessage(
    agentName,
    allAgents,
    MESSAGE_TYPES.EMERGENCY_ALERT,
    JSON.stringify(content),
    { priority: PRIORITY_LEVELS.URGENT }
  );
}

/**
 * Request approval for a task
 * @param {string} agentName - Agent name
 * @param {string} taskId - Task ID
 * @param {string} description - Task description
 * @param {string[]} approvers - List of approvers
 */
export function requestApproval(agentName, taskId, description, approvers) {
  const content = {
    taskId,
    description,
    requestedBy: agentName
  };
  
  sendMessage(
    agentName,
    approvers,
    MESSAGE_TYPES.APPROVAL_REQUEST,
    JSON.stringify(content),
    { priority: PRIORITY_LEVELS.HIGH, responseRequired: true }
  );
}

export default {
  sendMessage,
  getMessages,
  markMessageAsRead,
  archiveMessage,
  sendStatusUpdate,
  sendEmergencyAlert,
  requestApproval,
  MESSAGE_TYPES,
  PRIORITY_LEVELS
};