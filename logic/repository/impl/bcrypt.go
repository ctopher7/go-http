package impl

import "golang.org/x/crypto/bcrypt"

func (r *repository) BcryptComparePassword(hash, password []byte) error {
	return bcrypt.CompareHashAndPassword(hash, password)
}

func (r *repository) BcryptGenerateHash(password []byte) ([]byte, error) {
	return bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
}
