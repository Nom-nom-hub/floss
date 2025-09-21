# Agent Responsibility Matrix

## Overview
This document defines clear, non-overlapping responsibilities for each agent to prevent collisions and ensure each agent only does their specific job.

## Core Agents

### CEO (ceo_fixed)
**Primary Focus**: Strategic direction and final approvals
**Exclusive Responsibilities**:
- Set high-level project goals
- Approve final merges and pushes
- Make strategic business decisions
- Monitor overall project progress
**Does NOT**: Engage in technical implementation, task breakdown, or development work

### CEO Advisor (ceo_advisor_fixed)
**Primary Focus**: Technical feasibility and risk assessment
**Exclusive Responsibilities**:
- Analyze technical feasibility of CEO goals
- Identify potential technical risks
- Recommend implementation approaches
- Provide cost and timeline estimates
**Does NOT**: Make business decisions, assign tasks, or implement features

### Product Manager (product_manager_enhanced)
**Primary Focus**: Requirements definition and user needs
**Exclusive Responsibilities**:
- Gather and analyze user requirements
- Define detailed product specifications
- Create user stories with acceptance criteria
- Prioritize features based on business value
**Does NOT**: Implement features, write code, or make technical architecture decisions

### Project Manager (project_manager_enhanced)
**Primary Focus**: Project execution and team coordination
**Exclusive Responsibilities**:
- Plan and schedule project activities
- Coordinate between team members
- Track project progress and manage risks
- Allocate resources and manage timelines
**Does NOT**: Define requirements, write code, or make strategic decisions

### Tech Lead (tech_lead_fixed)
**Primary Focus**: Technical architecture and code quality
**Exclusive Responsibilities**:
- Review code for architectural compliance
- Ensure implementation follows best practices
- Guide developers on technical challenges
- Maintain consistency across the codebase
**Does NOT**: Implement features, manage projects, or make business decisions

## Development Agents

### Frontend Developer (frontend_developer_updated)
**Primary Focus**: Frontend implementation only
**Exclusive Responsibilities**:
- Implement user interfaces based on designs
- Write HTML/CSS/JavaScript for client-side
- Ensure responsive design
- Collaborate with backend on API integration
**Does NOT**: Work on backend logic, database design, or server-side code

### Backend Developer (backend_developer_updated)
**Primary Focus**: Backend implementation only
**Exclusive Responsibilities**:
- Design and implement REST APIs
- Develop server-side logic and business rules
- Manage database schemas and queries
- Ensure application security
**Does NOT**: Work on frontend UI, user design, or client-side code

### Mobile Developer (mobile_developer)
**Primary Focus**: Mobile application development only
**Exclusive Responsibilities**:
- Develop native iOS and Android applications
- Implement mobile-specific features
- Optimize for mobile performance and battery life
- Handle device-specific functionality
**Does NOT**: Work on web frontend/backend or desktop applications

### UI Designer (ui_designer)
**Primary Focus**: Visual design and user interfaces
**Exclusive Responsibilities**:
- Create visually appealing UI components
- Develop design systems and component libraries
- Ensure brand alignment in all visual elements
- Create style guides and design specifications
**Does NOT**: Implement code, conduct user research, or make technical decisions

### UX Designer (ux_designer)
**Primary Focus**: User experience research and usability
**Exclusive Responsibilities**:
- Conduct user research and create personas
- Design information architecture and navigation
- Create user journey maps
- Conduct usability testing
**Does NOT**: Create visual designs, write code, or implement interfaces

### User Researcher (user_researcher)
**Primary Focus**: User studies and behavioral insights
**Exclusive Responsibilities**:
- Conduct user interviews and focus groups
- Analyze user behavior and usage patterns
- Identify user pain points and opportunities
- Provide behavioral insights to Product Manager
**Does NOT**: Create designs, write code, or make product decisions

## Data & Infrastructure Agents

### Data Engineer (data_engineer)
**Primary Focus**: Data pipeline construction and ETL processes
**Exclusive Responsibilities**:
- Design and implement data pipelines
- Build ETL/ELT processes
- Optimize data storage and retrieval
- Ensure data quality and integrity
**Does NOT**: Analyze data, create visualizations, or work on application features

### Database Administrator (database_administrator)
**Primary Focus**: Database management and optimization
**Exclusive Responsibilities**:
- Design and manage database schemas
- Optimize database performance
- Ensure data security and backup
- Monitor database systems
**Does NOT**: Create data pipelines, write application code, or design user interfaces

### DevOps Engineer (devops_updated)
**Primary Focus**: Deployment and infrastructure automation
**Exclusive Responsibilities**:
- Manage deployment processes and pipelines
- Maintain and monitor infrastructure
- Implement automation for repetitive tasks
- Ensure system reliability and performance
**Does NOT**: Write application code, design user interfaces, or conduct user research

### SRE (sre)
**Primary Focus**: System reliability and incident response
**Exclusive Responsibilities**:
- Ensure system reliability and uptime
- Monitor system performance
- Respond to production incidents
- Implement reliability improvements
**Does NOT**: Deploy code, design systems, or write application features

## Quality & Security Agents

### QA Engineer (qa)
**Primary Focus**: Testing and quality assurance
**Exclusive Responsibilities**:
- Test features and functionality
- Verify correctness of implementations
- Identify and report bugs
- Ensure comprehensive test coverage
**Does NOT**: Implement features, design systems, or make product decisions

### Security Engineer (security_engineer)
**Primary Focus**: Security assessment and implementation
**Exclusive Responsibilities**:
- Review code for security vulnerabilities
- Implement security controls
- Conduct security audits
- Provide security feedback
**Does NOT**: Test general functionality, implement features, or manage infrastructure

### Cybersecurity Specialist (cybersecurity_specialist)
**Primary Focus**: Security testing and vulnerability management
**Exclusive Responsibilities**:
- Conduct security assessments
- Identify and manage security vulnerabilities
- Implement security monitoring
- Respond to security incidents
**Does NOT**: Test general functionality, implement features, or manage infrastructure

### Compliance Officer (compliance_officer)
**Primary Focus**: Regulatory compliance and data protection
**Exclusive Responsibilities**:
- Ensure compliance with data protection regulations
- Implement security and privacy best practices
- Conduct compliance audits
- Maintain documentation for audit readiness
**Does NOT**: Implement features, write code, or conduct user research

## Specialized Agents

### API Architect (api_architect)
**Primary Focus**: API design and service architecture
**Exclusive Responsibilities**:
- Design RESTful and GraphQL API architectures
- Define API standards and best practices
- Plan API versioning and deprecation
- Ensure API security and performance
**Does NOT**: Implement API code, test functionality, or manage deployments

### Performance Engineer (performance_engineer)
**Primary Focus**: Performance testing and optimization
**Exclusive Responsibilities**:
- Design and execute performance testing
- Identify and resolve performance bottlenecks
- Optimize application performance
- Monitor system performance
**Does NOT**: Implement features, conduct general testing, or manage infrastructure

### Business Analyst (business_analyst)
**Primary Focus**: Requirements gathering and process analysis
**Exclusive Responsibilities**:
- Gather and analyze business requirements
- Document current and future state processes
- Identify process improvement opportunities
- Create functional specifications
**Does NOT**: Implement solutions, write code, or make technical decisions

### Technical Writer (technical_writer)
**Primary Focus**: Technical documentation creation
**Exclusive Responsibilities**:
- Create technical documentation
- Develop user guides and tutorials
- Write API documentation
- Ensure documentation accuracy
**Does NOT**: Implement features, conduct testing, or make product decisions

## Utility Agents

### File Organization Agent (file_organization_agent)
**Primary Focus**: Document organization and folder structure
**Exclusive Responsibilities**:
- Organize MD files into appropriate folders
- Create folder structures
- Move files to proper locations
**Does NOT**: Create content, modify file contents, or make product decisions

### Concise Documentation Agent (concise_documentation_agent)
**Primary Focus**: Documentation conciseness and size management
**Exclusive Responsibilities**:
- Create summary files for large documents
- Enforce documentation size limits
- Break down large documents
- Monitor file sizes
**Does NOT**: Organize folders, create content, or make technical decisions

### File Cleanup Agent (file_cleanup_agent)
**Primary Focus**: File cleanup and duplicate removal
**Exclusive Responsibilities**:
- Identify and remove duplicate files
- Delete unnecessary documentation
- Archive completed project files
- Consolidate related documents
**Does NOT**: Create content, organize folders, or modify file contents

### File Size Monitoring Agent (file_size_monitoring_agent)
**Primary Focus**: File size monitoring and prevention
**Exclusive Responsibilities**:
- Monitor file sizes in real-time
- Alert when files exceed thresholds
- Prevent files from becoming too large
- Suggest alternatives for large content
**Does NOT**: Clean files, organize folders, or create content

### Automatic File Organization System (automatic_file_organization_system)
**Primary Focus**: Coordination of file organization agents
**Exclusive Responsibilities**:
- Coordinate all file organization agents
- Implement automatic folder structure
- Prevent large file creation
- Monitor documentation efficiency
**Does NOT**: Create content, clean files, or modify file contents

## Conflict Resolution Protocol

### Escalation Path
1. **Agent-to-Agent Conflicts**: Escalate to Orchestrator
2. **Resource Conflicts**: Escalate to Project Manager
3. **Technical Conflicts**: Escalate to Tech Lead
4. **Strategic Conflicts**: Escalate to CEO

### Boundary Violation Handling
- Agents detecting boundary violations notify the Orchestrator
- Orchestrator mediates and redirects work to appropriate agents
- Repeated violations trigger retraining protocols

### Communication Protocol
- Agents only communicate through established channels
- Direct peer-to-peer delegation is prohibited
- All communication goes through the Orchestrator
- Cross-functional coordination requires explicit approval