package model

type User struct {
	Id       int64  `db:"id" json:"id,omitempty"`
	Email    string `db:"email" json:"email,omitempty"`
	Password string `db:"password" json:"password,omitempty"`
	Role     string `db:"role" json:"role,omitempty"`
}
