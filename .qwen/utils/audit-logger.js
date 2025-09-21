/**
 * @file Audit Logging System
 * @description Maintains audit logs of all agent activities for compliance and monitoring
 * @license Apache-2.0
 */

import { writeFileSync, readFileSync, existsSync, mkdirSync, appendFileSync } from 'fs';
import { join } from 'path';
import { createHash, createCipheriv, createDecipheriv } from 'crypto';

// Ensure the audit directory exists
const AUDIT_DIR = join('.qwen', 'audit');
if (!existsSync(AUDIT_DIR)) {
  mkdirSync(AUDIT_DIR, { recursive: true });
}

// Encryption setup (in a real system, these would be securely managed)
const ENCRYPTION_KEY = process.env.AUDIT_ENCRYPTION_KEY || 'qwen_audit_key_32bytes_long!!';
const IV_LENGTH = 16;

/**
 * Audit log entry structure
 * @typedef {Object} AuditLogEntry
 * @property {string} id - Unique entry identifier
 * @property {string} timestamp - ISO timestamp
 * @property {string} agent - Agent name
 * @property {string} action - Action performed
 * @property {string} outcome - Outcome of the action
 * @property {Object} details - Additional details about the action
 * @property {string} [user] - User associated with the action
 * @property {string} [sessionId] - Session identifier
 * @property {string} [ipAddress] - IP address of the requester
 * @property {string} [resourceId] - Identifier of the resource affected
 * @property {boolean} [encrypted] - Whether the entry is encrypted
 */

/**
 * Log an audit event
 * @param {string} agent - Agent name
 * @param {string} action - Action performed
 * @param {string} outcome - Outcome of the action
 * @param {Object} details - Additional details about the action
 * @param {Object} [context] - Context information
 * @returns {string} Log entry ID
 */
export function logAuditEvent(agent, action, outcome, details, context = {}) {
  const entryId = generateEntryId();
  const timestamp = new Date().toISOString();
  
  const entry = {
    id: entryId,
    timestamp,
    agent,
    action,
    outcome,
    details,
    ...context
  };
  
  // Encrypt sensitive entries
  let logEntry = entry;
  if (shouldEncryptEntry(action, details)) {
    logEntry = encryptEntry(entry);
  }
  
  // Write to daily log file
  const logFile = join(AUDIT_DIR, `audit_${new Date().toISOString().split('T')[0]}.log`);
  appendFileSync(logFile, JSON.stringify(logEntry) + '\n');
  
  // Update index for quick lookup
  updateAuditIndex(entryId, timestamp, agent, action);
  
  return entryId;
}

/**
 * Retrieve audit logs
 * @param {Object} [filters] - Filters to apply
 * @param {string} [filters.agent] - Filter by agent name
 * @param {string} [filters.action] - Filter by action
 * @param {string} [filters.startDate] - Filter by start date
 * @param {string} [filters.endDate] - Filter by end date
 * @param {number} [filters.limit] - Limit number of results
 * @returns {AuditLogEntry[]} Array of audit log entries
 */
export function getAuditLogs(filters = {}) {
  // This is a simplified implementation
  // In a real system, this would query the index and read relevant log files
  const logs = [];
  
  // Read recent log files (last 30 days)
  const now = new Date();
  for (let i = 0; i < 30; i++) {
    const date = new Date(now);
    date.setDate(date.getDate() - i);
    const logFile = join(AUDIT_DIR, `audit_${date.toISOString().split('T')[0]}.log`);
    
    if (existsSync(logFile)) {
      const content = readFileSync(logFile, 'utf8');
      const lines = content.split('\n').filter(line => line.trim() !== '');
      
      lines.forEach(line => {
        try {
          const entry = JSON.parse(line);
          
          // Apply filters
          if (filters.agent && entry.agent !== filters.agent) return;
          if (filters.action && entry.action !== filters.action) return;
          if (filters.startDate && entry.timestamp < filters.startDate) return;
          if (filters.endDate && entry.timestamp > filters.endDate) return;
          
          // Decrypt if necessary
          const logEntry = entry.encrypted ? decryptEntry(entry) : entry;
          logs.push(logEntry);
        } catch (error) {
          // Skip malformed entries
        }
      });
    }
  }
  
  // Apply limit
  if (filters.limit) {
    return logs.slice(0, filters.limit);
  }
  
  return logs;
}

/**
 * Generate a unique entry ID
 * @returns {string} Unique entry ID
 */
function generateEntryId() {
  return `audit_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`;
}

/**
 * Determine if an entry should be encrypted
 * @param {string} action - Action performed
 * @param {Object} details - Additional details
 * @returns {boolean} Whether the entry should be encrypted
 */
function shouldEncryptEntry(action, details) {
  // Encrypt entries with sensitive information
  const sensitiveActions = [
    'user_authentication',
    'data_access',
    'file_modification',
    'configuration_change',
    'secrets_management'
  ];
  
  const sensitiveDataKeys = [
    'password',
    'token',
    'key',
    'secret',
    'credential',
    'api_key'
  ];
  
  if (sensitiveActions.includes(action)) {
    return true;
  }
  
  // Check if details contain sensitive data
  const detailsString = JSON.stringify(details).toLowerCase();
  return sensitiveDataKeys.some(key => detailsString.includes(key));
}

/**
 * Encrypt an audit entry
 * @param {AuditLogEntry} entry - Audit log entry to encrypt
 * @returns {Object} Encrypted entry
 */
function encryptEntry(entry) {
  try {
    // Create a hash of the original entry for integrity verification
    const entryString = JSON.stringify(entry);
    const hash = createHash('sha256').update(entryString).digest('hex');
    
    // Encrypt the entry
    const iv = Buffer.from(ENCRYPTION_KEY.slice(0, IV_LENGTH));
    const cipher = createCipheriv('aes-256-cbc', Buffer.from(ENCRYPTION_KEY), iv);
    
    let encrypted = cipher.update(entryString, 'utf8', 'hex');
    encrypted += cipher.final('hex');
    
    return {
      id: entry.id,
      timestamp: entry.timestamp,
      encrypted: true,
      hash,
      data: encrypted
    };
  } catch (error) {
    // If encryption fails, log an error and return the original entry
    console.error('Failed to encrypt audit entry:', error);
    return entry;
  }
}

/**
 * Decrypt an audit entry
 * @param {Object} encryptedEntry - Encrypted audit log entry
 * @returns {AuditLogEntry} Decrypted entry
 */
function decryptEntry(encryptedEntry) {
  try {
    if (!encryptedEntry.encrypted) {
      return encryptedEntry;
    }
    
    // Decrypt the entry
    const iv = Buffer.from(ENCRYPTION_KEY.slice(0, IV_LENGTH));
    const decipher = createDecipheriv('aes-256-cbc', Buffer.from(ENCRYPTION_KEY), iv);
    
    let decrypted = decipher.update(encryptedEntry.data, 'hex', 'utf8');
    decrypted += decipher.final('utf8');
    
    const entry = JSON.parse(decrypted);
    
    // Verify integrity
    const entryString = JSON.stringify(entry);
    const hash = createHash('sha256').update(entryString).digest('hex');
    
    if (hash !== encryptedEntry.hash) {
      throw new Error('Audit entry integrity check failed');
    }
    
    return entry;
  } catch (error) {
    console.error('Failed to decrypt audit entry:', error);
    throw new Error('Failed to decrypt audit entry');
  }
}

/**
 * Update the audit index for quick lookup
 * @param {string} entryId - Entry ID
 * @param {string} timestamp - Timestamp
 * @param {string} agent - Agent name
 * @param {string} action - Action performed
 */
function updateAuditIndex(entryId, timestamp, agent, action) {
  // In a real system, this would update a database index
  // For this implementation, we'll just maintain a simple index file
  const indexFile = join(AUDIT_DIR, 'index.json');
  let index = {};
  
  if (existsSync(indexFile)) {
    index = JSON.parse(readFileSync(indexFile, 'utf8'));
  }
  
  if (!index.entries) {
    index.entries = [];
  }
  
  index.entries.push({
    id: entryId,
    timestamp,
    agent,
    action
  });
  
  // Keep only recent entries in the index
  const cutoffDate = new Date();
  cutoffDate.setDate(cutoffDate.getDate() - 30);
  index.entries = index.entries.filter(entry => 
    new Date(entry.timestamp) > cutoffDate
  );
  
  writeFileSync(indexFile, JSON.stringify(index, null, 2));
}

/**
 * Generate a compliance report
 * @param {string} [startDate] - Start date for the report
 * @param {string} [endDate] - End date for the report
 * @returns {Object} Compliance report
 */
export function generateComplianceReport(startDate, endDate) {
  const logs = getAuditLogs({ startDate, endDate });
  
  const report = {
    period: {
      start: startDate || new Date(Date.now() - 30 * 24 * 60 * 60 * 1000).toISOString(),
      end: endDate || new Date().toISOString()
    },
    summary: {
      totalEntries: logs.length,
      uniqueAgents: [...new Set(logs.map(log => log.agent))].length,
      actions: {}
    },
    securityEvents: logs.filter(log => 
      log.action.includes('authentication') || 
      log.action.includes('authorization') ||
      log.action.includes('security')
    ),
    systemEvents: logs.filter(log => 
      log.action.includes('system') || 
      log.action.includes('configuration') ||
      log.action.includes('initialization')
    )
  };
  
  // Count actions
  logs.forEach(log => {
    if (!report.summary.actions[log.action]) {
      report.summary.actions[log.action] = 0;
    }
    report.summary.actions[log.action]++;
  });
  
  return report;
}

export default {
  logAuditEvent,
  getAuditLogs,
  generateComplianceReport
};