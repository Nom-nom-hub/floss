/**
 * @file Git Workflow Enforcement Utility
 * @description Utility functions to enforce Git workflow rules across all agents
 * @license Apache-2.0
 */

import { execSync } from 'child_process';
import { existsSync, readFileSync, writeFileSync } from 'fs';
import { join } from 'path';

/**
 * Validates branch name against required convention
 * @param {string} branchName - The branch name to validate
 * @returns {boolean} - Whether the branch name is valid
 */
export function validateBranchName(branchName) {
  // Branch naming convention: role/feature-short-description
  const validRoles = [
    'ceo', 'ceo_advisor', 'orchestrator', 'tech_lead', 
    'developer', 'qa', 'security_engineer', 'devops', 'documentation'
  ];
  
  // Check if branch follows the pattern role/feature-short-description
  const parts = branchName.split('/');
  if (parts.length < 2) {
    return false;
  }
  
  const role = parts[0];
  return validRoles.includes(role);
}

/**
 * Validates commit message format
 * @param {string} commitMessage - The commit message to validate
 * @returns {boolean} - Whether the commit message is valid
 */
export function validateCommitMessage(commitMessage) {
  // Commit message format: [ROLE] task_summary
  const validRoles = [
    'CEO', 'CEO_ADVISOR', 'ORCHESTRATOR', 'TECH_LEAD', 
    'DEVELOPER', 'QA', 'SECURITY_ENGINEER', 'DEVOPS', 'DOCUMENTATION'
  ];
  
  // Check if commit message starts with a valid role tag
  const roleMatch = commitMessage.match(/^\[([A-Z_]+)\]/);
  if (!roleMatch) {
    return false;
  }
  
  const role = roleMatch[1];
  return validRoles.includes(role);
}

/**
 * Prevents direct pushes to remote repositories
 * @param {string} command - The Git command being executed
 * @returns {boolean} - Whether the command should be blocked
 */
export function shouldBlockPush(command) {
  // Block direct push commands
  return command.includes('git push') && !command.includes('--dry-run');
}

/**
 * Installs Git hooks to enforce workflow rules
 * @param {string} repoPath - Path to the Git repository
 */
export function installGitHooks(repoPath) {
  const hooksDir = join(repoPath, '.git', 'hooks');
  
  // Pre-commit hook to validate commit messages
  const preCommitHook = `#!/bin/sh
# Git Workflow Enforcement - Pre-commit Hook

# Get the commit message
commitMessage=$(cat $1)

# Validate commit message format
if ! node -e "
const { validateCommitMessage } = require('./.qwen/utils/git-workflow-enforcer.js');
process.exit(validateCommitMessage(process.argv[1]) ? 0 : 1);
" "$commitMessage"; then
  echo "ERROR: Invalid commit message format"
  echo "Commit messages must follow the format: [ROLE] task_summary"
  echo "Valid roles: CEO, CEO_ADVISOR, ORCHESTRATOR, TECH_LEAD, DEVELOPER, QA, SECURITY_ENGINEER, DEVOPS, DOCUMENTATION"
  exit 1
fi
`;
  
  // Pre-push hook to prevent direct pushes
  const prePushHook = `#!/bin/sh
# Git Workflow Enforcement - Pre-push Hook

# Prevent all direct pushes
echo "ERROR: Direct pushes to remote repositories are not allowed"
echo "Please submit your changes through the proper PR process via the Orchestrator"
exit 1
`;
  
  // Write hook files
  writeFileSync(join(hooksDir, 'commit-msg'), preCommitHook);
  writeFileSync(join(hooksDir, 'pre-push'), prePushHook);
  
  // Make hooks executable
  execSync(`chmod +x ${join(hooksDir, 'commit-msg')}`);
  execSync(`chmod +x ${join(hooksDir, 'pre-push')}`);
}

/**
 * Validates current branch and provides corrective actions if needed
 * @param {string} repoPath - Path to the Git repository
 * @returns {Object} - Validation result with status and message
 */
export function validateCurrentBranch(repoPath) {
  try {
    const branchName = execSync('git rev-parse --abbrev-ref HEAD', {
      cwd: repoPath,
      stdio: ['pipe', 'pipe', 'ignore']
    }).toString().trim();
    
    if (validateBranchName(branchName)) {
      return {
        valid: true,
        message: `Branch "${branchName}" follows the required naming convention`
      };
    } else {
      return {
        valid: false,
        message: `Branch "${branchName}" does not follow the required naming convention`,
        suggestion: 'Branch names must follow the format: role/feature-short-description'
      };
    }
  } catch (error) {
    return {
      valid: false,
      message: 'Unable to determine current branch',
      error: error.message
    };
  }
}

/**
 * Creates a properly formatted commit
 * @param {string} role - The role making the commit
 * @param {string} message - The commit message summary
 * @param {string} repoPath - Path to the Git repository
 */
export function createFormattedCommit(role, message, repoPath) {
  const formattedMessage = `[${role.toUpperCase()}] ${message}`;
  
  if (!validateCommitMessage(formattedMessage)) {
    throw new Error('Unable to create valid commit message format');
  }
  
  try {
    execSync(`git commit -m "${formattedMessage}"`, {
      cwd: repoPath,
      stdio: 'inherit'
    });
  } catch (error) {
    throw new Error(`Failed to create commit: ${error.message}`);
  }
}

/**
 * Creates a branch with proper naming convention
 * @param {string} role - The role creating the branch
 * @param {string} feature - Short description of the feature
 * @param {string} repoPath - Path to the Git repository
 */
export function createFormattedBranch(role, feature, repoPath) {
  const branchName = `${role.toLowerCase()}/${feature.toLowerCase().replace(/\s+/g, '-')}`;
  
  if (!validateBranchName(branchName)) {
    throw new Error('Unable to create valid branch name');
  }
  
  try {
    execSync(`git checkout -b ${branchName}`, {
      cwd: repoPath,
      stdio: 'inherit'
    });
  } catch (error) {
    throw new Error(`Failed to create branch: ${error.message}`);
  }
}

export default {
  validateBranchName,
  validateCommitMessage,
  shouldBlockPush,
  installGitHooks,
  validateCurrentBranch,
  createFormattedCommit,
  createFormattedBranch
};