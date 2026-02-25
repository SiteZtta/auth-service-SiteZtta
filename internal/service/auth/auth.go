package auth

import (
	"auth-service-SiteZtta/config"
	"auth-service-SiteZtta/internal/domain/entities"
	"auth-service-SiteZtta/internal/storage"
	"auth-service-SiteZtta/internal/transport/grpc/v1/dto"
	"auth-service-SiteZtta/pkg/jwt"
	"context"
	"errors"
	"fmt"
	"log/slog"

	"golang.org/x/crypto/bcrypt"
)

type UserSaver interface {
	SaveUser(ctx context.Context, user *entities.User) (uid int64, err error)
}

type UserProvider interface {
	GetUserByEmail(ctx context.Context, email string) (user *entities.User, err error)
	GetUserByUsername(ctx context.Context, username string) (user *entities.User, err error)
}

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
)

// type TokenProvider interface {
// 	GenerateToken(ctx context.Context, signInInput *dto.SignInInput) (token string, err error)
// 	ValidateToken(ctx context.Context, token string) (user *entities.User, err error)
// }

type Auth struct {
	log          *slog.Logger
	userCreator  UserSaver
	userProvider UserProvider
	authConf     config.AuthConf
}

// New returns a new instance of the Auth service.
func New(
	log *slog.Logger,
	userCreator UserSaver,
	userProvider UserProvider,
	authConf config.AuthConf,
) *Auth {
	return &Auth{
		log:          log,
		userCreator:  userCreator,
		userProvider: userProvider,
		authConf:     authConf}
}

// CreateUser creates a new user and returns its ID.
func (a *Auth) CreateUser(ctx context.Context, signUpInput *dto.SignUpInput) (uid int64, err error) {
	const fn = "auth-service-SiteZtta.internal.service.auth.createUser"
	log := a.log.With(slog.String("fn", fn))
	log.Info("creating new user")
	passHash, err := bcrypt.GenerateFromPassword([]byte(signUpInput.Password), bcrypt.DefaultCost)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", fn, err)
	}
	userToSave := &entities.User{
		Email:    signUpInput.Email,
		Username: signUpInput.UserName,
		Phone:    signUpInput.Phone,
		PassHash: passHash,
	}
	if uid, err = a.userCreator.SaveUser(ctx, userToSave); err != nil {
		return 0, fmt.Errorf("%s: %w", fn, err)
	}
	log.Info("user created")
	return uid, nil
}

// GenerateToken checks with given credentials exists in the system and returns a token if it does.
//
// If user exists, but password is incorrect, returns an error.
// If user does not exist, returns an error.
func (a *Auth) GenerateToken(ctx context.Context, in dto.SignInInput) (token string, err error) {
	const fn = "auth-service-SiteZtta.internal.service.auth.generateToken"
	log := a.log.With(slog.String("fn", fn))
	log.Info("generating token for user")
	user := &entities.User{}
	user, err = a.userProvider.GetUserByEmail(ctx, in.Login)
	if user == nil {
		user, err = a.userProvider.GetUserByUsername(ctx, in.Login)
	}
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			return "", fmt.Errorf("%s: %w", fn, ErrInvalidCredentials)
		}
		return "", err
	}
	if err = bcrypt.CompareHashAndPassword(user.PassHash, []byte(in.Password)); err != nil {
		return "", fmt.Errorf("%s: %w", fn, ErrInvalidCredentials)
	}
	log.Info("found user", slog.Int64("id", user.ID))
	token, err = jwt.NewToken(*user, a.authConf)
	if err != nil {
		return "", fmt.Errorf("%s: %w", fn, err)
	}
	return token, nil
}

// ValidateToken validates the token and returns the user if it is valid.
func (a *Auth) ValidateToken(ctx context.Context, token string) (authInfo dto.AuthInfo, err error) {
	const fn = "auth-service-SiteZtta.internal.service.auth.validateToken"
	log := a.log.With(slog.String("fn", fn))
	log.Info("validating token")
	claims, err := jwt.ParseToken(token, a.authConf)
	if err != nil {
		return dto.AuthInfo{}, fmt.Errorf("%s: %w", fn, err)
	}
	authInfo = dto.AuthInfo{
		UserId: claims.UserId,
		Role:   claims.Role,
	}
	return authInfo, nil
}
