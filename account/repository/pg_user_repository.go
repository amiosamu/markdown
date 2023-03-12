package repository

import (
	"context"
	"github.com/amiosamu/markdown/account/model"
	"github.com/amiosamu/markdown/account/model/apperrors"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"log"
)

type PGUserRepository struct {
	DB *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) model.UserRepository {
	return &PGUserRepository{
		DB: db,
	}
}

func (r *PGUserRepository) FindByID(ctx context.Context, uid uuid.UUID) (*model.User, error) {
	user := &model.User{}
	q := "SELECT * FROM users WHERE uid=&1"
	if err := r.DB.Get(user, q, uid); err != nil {
		return user, apperrors.NewNotFound("uid", uid.String())
	}
	return user, nil
}

func (r *PGUserRepository) Create(ctx context.Context, u *model.User) error {
	q := "INSERT INTO users(email,password) VALUES ($1,$2) RETURNING *"
	if err := r.DB.Get(u, q, u.Email, u.Password); err != nil {
		if err, ok := err.(*pq.Error); ok && err.Code.Name() == "unique violation" {
			log.Printf("could not create user with email: %v. reason: %v\n", u.Email, err.Code.Name())
			return apperrors.NewConflict("email", u.Email)
		}
		log.Printf("Could not create a user with email: %v. Reason: %v\n", u.Email, err)
		return apperrors.NewInternalServerError()
	}
	return nil
}
