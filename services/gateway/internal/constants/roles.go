package constants

const (
	RoleAdministrator string = "administrator"
	RoleInstructor    string = "instructor"
	RoleSchool        string = "school"
	RoleStudent       string = "student"
)

var (
	RolesAll []string = []string{
		RoleAdministrator,
		RoleInstructor,
		RoleSchool,
		RoleStudent,
	}

	RolesAdmin []string = []string{
		RoleAdministrator,
		RoleSchool,
	}

	RolesLMS []string = []string{
		RoleInstructor,
		RoleStudent,
	}
)
