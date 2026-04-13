package constants

const (
	RoleAdministrator string = "administrator"
	RoleInstructor    string = "instructor"
	RoleSchool        string = "school"
	RoleStudent       string = "student"
)

var (
	// Roles: Administrator, Instructor, School, and Student
	AllRoles []string = []string{
		RoleAdministrator,
		RoleInstructor,
		RoleSchool,
		RoleStudent,
	}

	// Roles: Administrator, Instructor, and Student
	UserRoles []string = []string{
		RoleAdministrator,
		RoleInstructor,
		RoleStudent,
	}

	// Roles: Administrator and School
	AdminRoles []string = []string{
		RoleAdministrator,
		RoleSchool,
	}

	// Roles: Instructor and Student
	LMSRoles []string = []string{
		RoleInstructor,
		RoleStudent,
	}
)
