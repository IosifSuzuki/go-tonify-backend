package model

type Role string

const (
	Client    Role = "client"
	Freelance Role = "freelance"
)

func (r Role) Opposite() Role {
	if r == Client {
		return Freelance
	}
	return Client
}
