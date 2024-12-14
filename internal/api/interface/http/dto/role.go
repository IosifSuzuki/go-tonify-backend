package dto

type Role string

var (
	FreelancerRole Role = "freelancer"
	ClientRole     Role = "client"
)

func (r Role) Valid() bool {
	switch r {
	case FreelancerRole, ClientRole:
		return true
	default:
		return false
	}
}
