package user

import (
	"context"
	"errors"
	"fmt"
	"log"

	"film-library/src/internal/config"
	"film-library/src/internal/tools"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrPasswordIncorrect = errors.New("password is incorrect")
)

type Service struct {
	repo       UserRepository
	signingKey string
}

func NewService(repo UserRepository, cfg *config.Config) *Service {
	return &Service{
		repo:       repo,
		signingKey: cfg.SigningKey,
	}
}

func (s *Service) CreateUser(ctx context.Context, req *CreateUserRequest) (*CreateUserResponse, error) {
	const op = "user.Service.CreateUser"

	vErr := ValidateCreateUserReuqest(req)
	if vErr != nil {
		log.Printf("ERROR: failed request validation")
		return nil, fmt.Errorf("%s: %w", op, vErr)
	}

	passhash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("ERROR: failed to generate password hash\n")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	u := &User{
		Username: req.Username,
		PassHash: string(passhash),
		IsAdmin:  false,
	}

	u, err = s.repo.CreateUser(ctx, u)
	if err != nil {
		log.Printf("ERROR: failed to create user record in repository\n")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &CreateUserResponse{
		ID:       int(u.ID),
		Username: u.Username,
	}, nil
}

func (s *Service) Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error) {
	const op = "user.Service.Login"

	vErr := ValidateLoginReuqest(req)
	if vErr != nil {
		log.Printf("ERROR: failed request validation")
		return nil, fmt.Errorf("%s: %w", op, vErr)
	}

	u, err := s.repo.GetUserByUsername(ctx, req.Username)
	if err != nil {
		log.Printf("ERROR: failed to get user record from repository\n")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(u.PassHash), []byte(req.Password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			log.Printf("ERROR: mismathced password and hash\n")
			return nil, fmt.Errorf("%s: %w", op, ErrPasswordIncorrect)
		}

		log.Printf("ERROR: failed password and hash comparison\n")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	claims := tools.NewUserClaims(int(u.ID), u.Username, u.IsAdmin)
	ss, err := tools.NewJWTSignedString(claims, s.signingKey)
	if err != nil {
		log.Printf("ERROR: failed to sign user claims\n")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &LoginResponse{
		AccessToken: ss,
	}, nil
}
