package model

import "github.com/google/uuid"

type User struct {
	UID      uuid.UUID `db:"uid" json:"uid"`
	Name     string    `db:"name" json:"name"`
	Email    string    `db:"email" json:"email"`
	Password string    `db:"password" json:"password"`
}
