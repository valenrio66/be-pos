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

func (r *userRepositoryImpl) FindByID(ctx context.Context, id int64) (*domain.User, error) {
	user := new(domain.User)
	err := r.db.NewSelect().Model(user).Where("id = ?", id).Scan(ctx)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *userRepositoryImpl) SaveRefreshToken(ctx context.Context, rt *domain.RefreshToken) error {
	_, err := r.db.NewInsert().Model(rt).Exec(ctx)
	return err
}

func (r *userRepositoryImpl) FindRefreshToken(ctx context.Context, token string) (*domain.RefreshToken, error) {
	rt := new(domain.RefreshToken)
	err := r.db.NewSelect().Model(rt).Where("token = ?", token).Scan(ctx)
	if err != nil {
		return nil, err
	}
	return rt, nil
}

func (r *userRepositoryImpl) DeleteRefreshToken(ctx context.Context, token string) error {
	_, err := r.db.NewDelete().Model((*domain.RefreshToken)(nil)).Where("token = ?", token).Exec(ctx)
	return err
}
