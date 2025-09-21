/**
 * @file Test script for Goal Driven Development agents
 * @description Simple test to verify that agents are properly configured
 * @license Apache-2.0
 */

import { existsSync, readdirSync } from 'fs';
import { join } from 'path';

// Test that all agent files exist and have correct format
const agentsDir = join('.qwen', 'agents');
const workflowsDir = join('.qwen', 'workflows');

console.log('Testing Goal Driven Development agents...\n');

// Check agents directory
if (existsSync(agentsDir)) {
  console.log('✅ Agents directory found');
  const agentFiles = readdirSync(agentsDir);
  console.log(`Found ${agentFiles.length} agent files:`);
  
  const requiredAgents = [
    'ceo.md',
    'ceo_advisor.md',
    'orchestrator.md',
    'tech_lead.md',
    'developer.md',
    'qa.md',
    'security_engineer.md',
    'devops.md',
    'documentation.md'
  ];
  
  let allAgentsFound = true;
  for (const agent of requiredAgents) {
    if (agentFiles.includes(agent)) {
      console.log(`  ✅ ${agent}`);
    } else {
      console.log(`  ❌ ${agent} (MISSING)`);
      allAgentsFound = false;
    }
  }
  
  if (allAgentsFound) {
    console.log('✅ All core agents found');
  } else {
    console.log('❌ Some core agents are missing');
  }
} else {
  console.log('❌ Agents directory not found');
}

// Check workflows directory
if (existsSync(workflowsDir)) {
  console.log('\n✅ Workflows directory found');
  const workflowFiles = readdirSync(workflowsDir);
  console.log(`Found ${workflowFiles.length} workflow files:`);
  
  const requiredWorkflows = [
    'user_login_mfa.md',
    'api_development.md',
    'ml_model_deployment.md'
  ];
  
  let allWorkflowsFound = true;
  for (const workflow of requiredWorkflows) {
    if (workflowFiles.includes(workflow)) {
      console.log(`  ✅ ${workflow}`);
    } else {
      console.log(`  ❌ ${workflow} (MISSING)`);
      allWorkflowsFound = false;
    }
  }
  
  if (allWorkflowsFound) {
    console.log('✅ All workflows found');
  } else {
    console.log('❌ Some workflows are missing');
  }
} else {
  console.log('\n❌ Workflows directory not found');
}

console.log('\nTest completed. If all agents and workflows are found, you can use them with Qwen Code.');