package agent

// AgentRole represents the role of an agent in the virtual company
type AgentRole string

const (
	AgentRoleCEO             AgentRole = "ceo"
	AgentRoleCTO             AgentRole = "cto"
	AgentRoleProductManager  AgentRole = "product_manager"
	AgentRoleTechLead        AgentRole = "tech_lead"
	AgentRoleSeniorDeveloper AgentRole = "senior_developer"
	AgentRoleJuniorDeveloper AgentRole = "junior_developer"
	AgentRoleQAEngineer      AgentRole = "qa_engineer"
	AgentRoleDevOpsEngineer  AgentRole = "devops_engineer"
	AgentRoleUIDesigner      AgentRole = "ui_designer"
	AgentRoleBusinessAnalyst AgentRole = "business_analyst"
	AgentRoleProjectManager  AgentRole = "project_manager"
	AgentRoleSecurityEngineer AgentRole = "security_engineer"
	AgentRoleDatabaseAdmin   AgentRole = "database_admin"
	AgentRoleTechnicalWriter AgentRole = "technical_writer"
)

// AgentSpecialization represents the specialization area of an agent
type AgentSpecialization string

const (
	SpecializationFrontend    AgentSpecialization = "frontend"
	SpecializationBackend     AgentSpecialization = "backend"
	SpecializationMobile      AgentSpecialization = "mobile"
	SpecializationDevOps      AgentSpecialization = "devops"
	SpecializationDataScience AgentSpecialization = "data_science"
	SpecializationAI          AgentSpecialization = "ai"
	SpecializationSecurity    AgentSpecialization = "security"
	SpecializationDatabase    AgentSpecialization = "database"
)

// CompanyStructure defines the organizational hierarchy of the virtual company
type CompanyStructure struct {
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Departments []Department `json:"departments"`
}

// Department represents a department in the virtual company
type Department struct {
	Name        string       `json:"name"`
	Description string       `json:"description"`
	HeadRole    AgentRole    `json:"head_role"`
	Teams       []Team       `json:"teams"`
}

// Team represents a team within a department
type Team struct {
	Name           string              `json:"name"`
	Description    string              `json:"description"`
	LeadRole       AgentRole           `json:"lead_role"`
	MemberRoles    []AgentRole         `json:"member_roles"`
	Specialization AgentSpecialization `json:"specialization"`
}

// DefaultCompanyStructure returns the default organizational structure for the virtual company
func DefaultCompanyStructure() CompanyStructure {
	return CompanyStructure{
		Name:        "Virtual Development Company",
		Description: "A complete virtual software development company with all roles filled by AI agents",
		Departments: []Department{
			{
				Name:        "Executive Leadership",
				Description: "Top-level strategic leadership",
				HeadRole:    AgentRoleCEO,
				Teams:       []Team{},
			},
			{
				Name:        "Technology Leadership",
				Description: "Technical direction and architecture",
				HeadRole:    AgentRoleCTO,
				Teams:       []Team{},
			},
			{
				Name:        "Product Management",
				Description: "Product strategy and requirements",
				HeadRole:    AgentRoleProductManager,
				Teams:       []Team{},
			},
			{
				Name:        "Engineering",
				Description: "Software development and implementation",
				HeadRole:    AgentRoleTechLead,
				Teams: []Team{
					{
						Name:           "Frontend Team",
						Description:    "Frontend development specialists",
						LeadRole:       AgentRoleSeniorDeveloper,
						MemberRoles:    []AgentRole{AgentRoleJuniorDeveloper, AgentRoleJuniorDeveloper},
						Specialization: SpecializationFrontend,
					},
					{
						Name:           "Backend Team",
						Description:    "Backend development specialists",
						LeadRole:       AgentRoleSeniorDeveloper,
						MemberRoles:    []AgentRole{AgentRoleJuniorDeveloper, AgentRoleJuniorDeveloper},
						Specialization: SpecializationBackend,
					},
					{
						Name:           "Mobile Team",
						Description:    "Mobile application development specialists",
						LeadRole:       AgentRoleSeniorDeveloper,
						MemberRoles:    []AgentRole{AgentRoleJuniorDeveloper},
						Specialization: SpecializationMobile,
					},
				},
			},
			{
				Name:        "Quality Assurance",
				Description: "Testing and quality control",
				HeadRole:    AgentRoleQAEngineer,
				Teams:       []Team{},
			},
			{
				Name:        "DevOps",
				Description: "Deployment and infrastructure management",
				HeadRole:    AgentRoleDevOpsEngineer,
				Teams:       []Team{},
			},
			{
				Name:        "Design",
				Description: "User interface and experience design",
				HeadRole:    AgentRoleUIDesigner,
				Teams:       []Team{},
			},
			{
				Name:        "Business Analysis",
				Description: "Requirements analysis and business process optimization",
				HeadRole:    AgentRoleBusinessAnalyst,
				Teams:       []Team{},
			},
			{
				Name:        "Project Management",
				Description: "Project coordination and timeline management",
				HeadRole:    AgentRoleProjectManager,
				Teams:       []Team{},
			},
			{
				Name:        "Security",
				Description: "Security auditing and vulnerability assessment",
				HeadRole:    AgentRoleSecurityEngineer,
				Teams:       []Team{},
			},
			{
				Name:        "Database Administration",
				Description: "Database management and optimization",
				HeadRole:    AgentRoleDatabaseAdmin,
				Teams:       []Team{},
			},
			{
				Name:        "Technical Documentation",
				Description: "Technical writing and documentation",
				HeadRole:    AgentRoleTechnicalWriter,
				Teams:       []Team{},
			},
		},
	}
}