package passwordhasherinterface

type PasswordHasher interface {
	Hash(rawPassword string) (string, error)
	Compare(passwordHash, rawPassword string) error
}
