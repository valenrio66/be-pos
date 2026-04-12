package postgres

import (
	"context"

	"github.com/uptrace/bun"
	"github.com/valenrio66/be-pos/internal/domain"
)

type userRepositoryImpl struct {
	db *bun.DB
}

func NewUserRepository(db *bun.DB) domain.UserRepository {
	return &userRepositoryImpl{db: db}
}

func (r *userRepositoryImpl) Create(ctx context.Context, user *domain.User) error {
	_, err := r.db.NewInsert().Model(user).Exec(ctx)
	return err
}

func (r *userRepositoryImpl) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	user := new(domain.User)

	err := r.db.NewSelect().Model(user).Where("email = ?", email).Scan(ctx)
	if err != nil {
		return nil, err
	}

	return user, nil
}
