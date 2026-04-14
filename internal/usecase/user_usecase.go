package usecase

import (
	"context"
	"database/sql"
	"errors"
	"time"

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
		u.log.Error("Password hashing failed", zap.Error(err))
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
		u.log.Error("Failed to add the user to the database", zap.Error(err))
		return nil, err
	}

	u.log.Info("The user has been successfully created", zap.String("email", email))
	return user, nil
}

func (u *UserUsecase) Login(ctx context.Context, email, password string) (*dto.LoginResponse, error) {
	user, err := u.repo.FindByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("incorrect email address or password")
		}
		u.log.Error("Database error during login", zap.Error(err))
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, errors.New("incorrect email address or password")
	}

	accessToken, err := utils.GenerateToken(user.ID, user.Role, u.cfg.JWTSecret)
	if err != nil {
		return nil, errors.New("failed to process token")
	}

	refreshTokenStr, err := utils.GenerateRefreshToken()
	if err != nil {
		return nil, errors.New("failed to create login session")
	}

	rt := &domain.RefreshToken{
		UserID:    user.ID,
		Token:     refreshTokenStr,
		ExpiredAt: time.Now().AddDate(0, 0, 7),
	}
	if err := u.repo.SaveRefreshToken(ctx, rt); err != nil {
		u.log.Error("Gagal simpan refresh token", zap.Error(err))
		return nil, errors.New("gagal memproses login")
	}

	return &dto.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshTokenStr,
		Role:         user.Role,
	}, nil
}

func (u *UserUsecase) RefreshAccessToken(ctx context.Context, refreshToken string) (*dto.LoginResponse, error) {
	rt, err := u.repo.FindRefreshToken(ctx, refreshToken)
	if err != nil {
		return nil, errors.New("the session is invalid or has expired")
	}

	if time.Now().After(rt.ExpiredAt) {
		_ = u.repo.DeleteRefreshToken(ctx, refreshToken)
		return nil, errors.New("your session has expired. please log in again")
	}

	user, err := u.repo.FindByID(ctx, rt.UserID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	newAccessToken, err := utils.GenerateToken(user.ID, user.Role, u.cfg.JWTSecret)
	if err != nil {
		return nil, errors.New("failed to process the new token")
	}

	return &dto.LoginResponse{
		AccessToken:  newAccessToken,
		RefreshToken: refreshToken,
		Role:         user.Role,
	}, nil
}

func (u *UserUsecase) GetProfile(ctx context.Context, userID int64) (*domain.User, error) {
	user, err := u.repo.FindByID(ctx, userID)
	if err != nil {
		u.log.Error("Failed to retrieve user profile", zap.Error(err))
		return nil, errors.New("user not found")
	}

	return user, nil
}
