package services

import (
	"context"
	"errors"
	"go-api/internal/db"
	"go-api/internal/repository"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Register(ctx context.Context, username, email, password string) (db.User, error)
	Login(ctx context.Context, email, password string) (string, string, error)
}

type authService struct {
	repo repository.UserRepository
}

func NewAuthService(repo repository.UserRepository) AuthService {
	return &authService{
		repo: repo,
	}
}

func (s *authService) Register(ctx context.Context, username, email, password string) (db.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return db.User{}, errors.New("gagal memproses password")
	}
	arg := db.CreateUserParams{
		Username:     username,
		Email:        email,
		PasswordHash: string(hashedPassword),
	}
	user, err := s.repo.CreateUser(ctx, arg)
	if err != nil {
		return db.User{}, err 
	}
	return user, nil
}

func (s *authService) Login(ctx context.Context, email, password string) (string, string, error) {
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return "", "", errors.New("email atau password salah") 
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return "", "", errors.New("email atau password salah")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(1 * time.Hour).Unix(),
	})
	
	accessToken, err := token.SignedString([]byte("RAHASIA_NEGARA"))
	if err != nil {
		return "", "", errors.New("gagal membuat access token")
	}

	sessionID := uuid.New()
	refreshToken := uuid.New().String()

	var pgSessionID pgtype.UUID
	_ = pgSessionID.Scan(sessionID.String())
	_, err = s.repo.CreateSession(ctx, db.CreateSessionParams{
		ID:           pgSessionID,
		UserID:       user.ID,
		RefreshToken: refreshToken,
		ExpiresAt:    pgtype.Timestamp{Time: time.Now().Add(7 * 24 * time.Hour), Valid: true}, 
	})
	if err != nil {
		return "", "", errors.New("gagal menyimpan sesi")
	}

	return accessToken, refreshToken, nil
}

