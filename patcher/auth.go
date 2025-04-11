package patcher

type PortalUser struct {
	Name         string `json:"name"`
	Email        string `json:"email"`
	PasswordHash string `json:"passwordHash"`
}

func (user PortalUser) Login(users []PortalUser) bool {
	for _, u := range users {
		if u.Email == user.Email && u.PasswordHash == user.PasswordHash {
			return true
		}
	}

	return false
}
