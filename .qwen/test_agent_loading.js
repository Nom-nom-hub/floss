/**
 * @file Test script to verify agent loading
 * @description Simple test to verify that agents can be loaded by Qwen Code
 * @license Apache-2.0
 */

import { readFileSync } from 'fs';
import { join } from 'path';

// Test parsing of a simple agent file
const testAgentPath = join('.qwen', 'agents', 'test-agent.md');

console.log('Testing agent loading...\n');

try {
  const content = readFileSync(testAgentPath, 'utf8');
  console.log('✅ Test agent file found and readable');
  
  // Check for required components
  if (content.startsWith('---')) {
    console.log('✅ YAML frontmatter delimiter found');
  } else {
    console.log('❌ YAML frontmatter delimiter missing');
  }
  
  if (content.includes('name:')) {
    console.log('✅ Name field found');
  } else {
    console.log('❌ Name field missing');
  }
  
  if (content.includes('description:')) {
    console.log('✅ Description field found');
  } else {
    console.log('❌ Description field missing');
  }
  
  if (content.includes('---\n\n')) {
    console.log('✅ System prompt section found');
  } else {
    console.log('❌ System prompt section missing');
  }
  
  console.log('\nTest agent content:');
  console.log('==================');
  console.log(content);
  
} catch (error) {
  console.log('❌ Error reading test agent file:', error.message);
}

console.log('\nIf all checks pass, the agent should be compatible with Qwen Code.');
