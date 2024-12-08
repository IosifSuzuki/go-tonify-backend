package entity

type Role struct {
	value string
}

var (
	UnknownRole    = Role{value: "unknown"}
	ClientRole     = Role{value: "client"}
	FreelancerRole = Role{value: "freelancer"}
)

func (r Role) Opposite() Role {
	switch r {
	case ClientRole:
		return FreelancerRole
	case FreelancerRole:
		return ClientRole
	default:
		return UnknownRole
	}
}

func RoleFromString(text string) (Role, error) {
	switch text {
	case ClientRole.value:
		return ClientRole, nil
	case FreelancerRole.value:
		return FreelancerRole, nil
	default:
		return UnknownRole, UnknownValueError
	}
}

func (r Role) String() string {
	return r.value
}
