package database

func isEmailUnique(email string, users map[string]User) bool {
	for _, user := range users {
		if user.Email == email {
			return false
		}
	}
	return true
}
