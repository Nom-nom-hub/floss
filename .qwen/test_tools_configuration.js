/**
 * @file Test script to verify agent tools configuration
 * @description Simple test to verify that agents have the correct tools configuration
 * @license Apache-2.0
 */

import { readFileSync } from 'fs';
import { join } from 'path';

// Test that agents have the correct tools configuration
const agentsDir = join('.qwen', 'agents');

console.log('Testing Goal Driven Development agents tools configuration...\\n');

const agentFiles = [
  { name: 'ceo.md', shouldHaveTaskTool: false },
  { name: 'ceo_advisor.md', shouldHaveTaskTool: false },
  { name: 'orchestrator.md', shouldHaveTaskTool: false },
  { name: 'tech_lead.md', shouldHaveTaskTool: false },
  { name: 'developer.md', shouldHaveTaskTool: true },
  { name: 'qa.md', shouldHaveTaskTool: false },
  { name: 'security_engineer.md', shouldHaveTaskTool: false },
  { name: 'devops.md', shouldHaveTaskTool: false },
  { name: 'documentation.md', shouldHaveTaskTool: false }
];

let allAgentsCorrect = true;

for (const { name: agentFile, shouldHaveTaskTool } of agentFiles) {
  const filePath = join(agentsDir, agentFile);
  try {
    const content = readFileSync(filePath, 'utf8');
    
    // Check if the file has the tools section by looking for the tools field in the YAML
    const lines = content.split('\\n');
    let inToolsSection = false;
    let hasTaskTool = false;
    
    for (const line of lines) {
      if (line.trim() === 'tools:') {
        inToolsSection = true;
        continue;
      }
      
      if (inToolsSection) {
        if (line.trim().startsWith('- ')) {
          if (line.includes('task')) {
            hasTaskTool = true;
          }
        } else if (line.trim() === '---' || line.trim().startsWith('name:') || line.trim().startsWith('description:')) {
          // End of tools section
          break;
        }
      }
    }
    
    if (hasTaskTool === shouldHaveTaskTool) {
      console.log(`  ✅ ${agentFile} has correct Task tool configuration`);
    } else {
      if (shouldHaveTaskTool) {
        console.log(`  ❌ ${agentFile} should have Task tool but doesn't`);
      } else {
        console.log(`  ❌ ${agentFile} should NOT have Task tool but does`);
      }
      allAgentsCorrect = false;
    }
  } catch (error) {
    console.log(`  ❌ ${agentFile} could not be read: ${error.message}`);
    allAgentsCorrect = false;
  }
}

if (allAgentsCorrect) {
  console.log('\\n✅ All agents have correct tools configuration');
  console.log('The agents should now work properly without getting stuck in delegation loops.');
} else {
  console.log('\\n❌ Some agents have incorrect tools configuration');
  console.log('Please fix the tools configuration for the agents listed above.');
}

console.log('\\nTest completed. If all agents have correct tools configuration, they should work properly with Qwen Code.');