package impl

import "golang.org/x/crypto/bcrypt"

func (r *repository) BcryptComparePassword(hash, password []byte) error {
	return bcrypt.CompareHashAndPassword(hash, password)
}
