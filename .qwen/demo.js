/**
 * @file Demo script for Goal Driven Development agents
 * @description Demonstration of how to use the agents to build a project
 * @license Apache-2.0
 */

console.log('Goal Driven Development System - Demo');
console.log('=====================================\n');

console.log('To use the Goal Driven Development agents with Qwen Code:\n');

console.log('1. Start Qwen Code in your project directory:');
console.log('   $ qwen\n');

console.log('2. Set a goal with the CEO agent:');
console.log('   > Act as the CEO and create an offline notes application\n');

console.log('3. The system will automatically:');
console.log('   - Have the CEO set the goal');
console.log('   - Have the CEO_Advisor provide technical recommendations');
console.log('   - Have the Orchestrator break down the tasks');
console.log('   - Assign work to Developers, QA, Security_Engineer, etc.');
console.log('   - Automatically coordinate the entire development process\n');

console.log('Example goal for testing:');
console.log('   > Create a simple offline notes app with local storage\n');

console.log('The agents will coordinate automatically and you\'ll see real-time updates like:');
console.log('   [DEVELOPER] Creating HTML structure for notes app...');
console.log('   [QA] Testing offline functionality...');
console.log('   [SECURITY] Reviewing client-side data storage...\n');

console.log('All agents now have loop prevention rules to avoid infinite delegation.');