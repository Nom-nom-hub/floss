/**
 * @file Test script for Goal Driven Development agents loop prevention
 * @description Test to verify that agents won't get stuck in infinite loops
 * @license Apache-2.0
 */

import { readFileSync } from 'fs';
import { join } from 'path';

// Test that all agent files have the loop prevention rules
const agentsDir = join('.qwen', 'agents');

console.log('Testing Goal Driven Development agents for loop prevention...\\n');

const agentFiles = [
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

let allAgentsHaveLoopPrevention = true;

for (const agentFile of agentFiles) {
  const filePath = join(agentsDir, agentFile);
  try {
    const content = readFileSync(filePath, 'utf8');
    if (content.includes('## Important Rules for Task Completion:') && 
        content.includes('Do NOT delegate tasks back to the same agent') &&
        content.includes('Do NOT create infinite loops')) {
      console.log(`  ✅ ${agentFile} has loop prevention rules`);
    } else {
      console.log(`  ❌ ${agentFile} is missing loop prevention rules`);
      allAgentsHaveLoopPrevention = false;
    }
  } catch (error) {
    console.log(`  ❌ ${agentFile} could not be read: ${error.message}`);
    allAgentsHaveLoopPrevention = false;
  }
}

if (allAgentsHaveLoopPrevention) {
  console.log('\\n✅ All agents have loop prevention rules');
  console.log('The agents should now work properly without getting stuck in infinite loops.');
} else {
  console.log('\\n❌ Some agents are missing loop prevention rules');
  console.log('Please add the loop prevention rules to all agents.');
}

console.log('\\nTest completed. If all agents have loop prevention rules, they should work properly with Qwen Code.');
