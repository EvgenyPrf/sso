package auth

import (
	"context"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
	"sso/internal/domain/models"
	"sso/internal/lib/jwt"
	"sso/internal/storage"
	"time"
)

type Auth struct {
	log          *slog.Logger
	userSaver    UserSaver
	userProvider UserProvider
	appProvider  AppProvider
	tokenTTL     time.Duration
}

type UserSaver interface {
	SaveUser(
		ctx context.Context,
		email string,
		passHash []byte) (uid int64, err error)
}

type UserProvider interface {
	User(ctx context.Context, email string) (models.User, error)
	IsAdmin(ctx context.Context, userId int64) (bool, error)
}

type AppProvider interface {
	App(ctx context.Context, appId int) (models.App, error)
}

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidUserId      = errors.New("invalid user id")
	ErrUserExists         = errors.New("user exists")
)

func New(
	log *slog.Logger,
	userSaver UserSaver,
	userProvider UserProvider,
	appProvider AppProvider,
	tokenTTL time.Duration) *Auth {
	return &Auth{log: log, userSaver: userSaver, userProvider: userProvider, appProvider: appProvider, tokenTTL: tokenTTL}
}

func (auth *Auth) Login(
	ctx context.Context,
	email, password string,
	appID int) (string, error) {
	const op = "auth.Login"

	log := auth.log.With(
		slog.String("op", op),
		slog.String("username", email),
	)

	log.Info("attempting to login user")
	user, err := auth.userProvider.User(ctx, email)

	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			auth.log.Warn("user not found", err)
			return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		}

		if errors.Is(err, storage.ErrAppNotFound) {
			auth.log.Warn("app not found", err)
			return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		}

		auth.log.Warn("failed to get user", err)
		return "", fmt.Errorf("%s: %w", op, err)
	}

	if err := bcrypt.CompareHashAndPassword(user.PassHash, []byte(password)); err != nil {
		auth.log.Info("Invalid credentials", err)
		return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}

	app, err := auth.appProvider.App(ctx, appID)

	if err != nil {
		auth.log.Info("Error getting app id", err)
		return "", fmt.Errorf("%s: %w", op, err)
	}

	log.Info("User logged in successfully")

	token, err := jwt.NewToken(user, app, auth.tokenTTL)

	if err != nil {
		auth.log.Info("Failed to generate token", err)
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return token, nil
}

func (auth *Auth) RegisterNewUser(
	ctx context.Context,
	email, password string) (int64, error) {
	const op = "auth.RegisterNewUser"

	log := auth.log.With(
		slog.String("op", op),
		slog.String("email", email),
	)

	log.Info("Register new user")

	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		log.Error("failed to generate password hash", err)
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	id, err := auth.userSaver.SaveUser(ctx, email, passHash)

	if err != nil {
		if errors.Is(err, storage.ErrUserExists) {
			auth.log.Warn("user already exists", err)
			return 0, fmt.Errorf("%s: %w", op, ErrUserExists)
		}
		log.Error("failed to save user", err)
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("User registered")

	return id, nil

}

func (auth *Auth) IsAdmin(
	ctx context.Context,
	userId int64) (bool, error) {

	const op = "auth.IsAdmin"

	log := auth.log.With(
		slog.String("op", op),
		slog.Int64("userId", userId),
	)

	log.Info("Checking if user is admin")

	isAdmin, err := auth.IsAdmin(ctx, userId)

	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			auth.log.Warn("user not found", err)
			return false, fmt.Errorf("%s: %w", op, ErrInvalidUserId)
		}
		return false, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("Checking if user is admin", slog.Bool("is_admin", isAdmin))

	return isAdmin, nil
}
