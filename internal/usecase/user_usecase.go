package usecase

import (
	"context"
	"database/sql"
	"errors"

	"github.com/valenrio66/be-pos/config"
	"github.com/valenrio66/be-pos/internal/delivery/http/dto"
	"github.com/valenrio66/be-pos/internal/domain"
	"github.com/valenrio66/be-pos/pkg/utils"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type UserUsecase struct {
	repo domain.UserRepository
	cfg  *config.Config
	log  *zap.Logger
}

func NewUserUsecase(repo domain.UserRepository, cfg *config.Config, log *zap.Logger) *UserUsecase {
	return &UserUsecase{
		repo: repo,
		cfg:  cfg,
		log:  log,
	}
}

func (u *UserUsecase) CreateUser(ctx context.Context, name, email, password string) (*domain.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		u.log.Error("Gagal hash password", zap.Error(err))
		return nil, err
	}

	user := &domain.User{
		Name:     name,
		Email:    email,
		Password: string(hashedPassword),
		Role:     "user",
	}

	err = u.repo.Create(ctx, user)
	if err != nil {
		u.log.Error("Gagal insert user ke database", zap.Error(err))
		return nil, err
	}

	u.log.Info("User berhasil dibuat", zap.String("email", email))
	return user, nil
}

func (u *UserUsecase) Login(ctx context.Context, email, password string) (*dto.LoginResponse, error) {
	user, err := u.repo.FindByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("email atau password salah")
		}
		u.log.Error("Database error saat login", zap.Error(err))
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, errors.New("email atau password salah")
	}

	token, err := utils.GenerateToken(user.ID, user.Role, u.cfg.JWTSecret)
	if err != nil {
		u.log.Error("Gagal generate token", zap.Error(err))
		return nil, errors.New("gagal memproses login")
	}

	u.log.Info("User berhasil login", zap.String("email", email))

	return &dto.LoginResponse{
		Token: token,
		Role:  user.Role,
	}, nil
}
