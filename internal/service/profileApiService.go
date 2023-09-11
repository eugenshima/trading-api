package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/eugenshima/trading-api/internal/model"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type ProfileService struct {
	rps ProfileApiRepository
}

func NewProfileService(rps ProfileApiRepository) *ProfileService {
	return &ProfileService{rps: rps}
}

type ProfileApiRepository interface {
	CreateProfile(context.Context, *model.User) error
	Login(context.Context, string, []byte) (uuid.UUID, error)
	UpdateProfile(context.Context, uuid.UUID, []byte) error
}

const (
	key             = "ew4t137tr1eyfg1ryg4ryerg2743gr2"
	accessTokenTTL  = 24 * time.Hour
	refreshTokenTTL = 72 * time.Hour
)

// tokenClaims struct contains information about the claims associated with the given token
type tokenClaims struct {
	jwt.StandardClaims
}

// Login functions compares password and generates access&refresh tokens if password is valid
func (s *ProfileService) Login(ctx context.Context, auth *model.Login) (*model.JWTResponse, error) {
	id, err := s.rps.Login(ctx, auth.Login, []byte(auth.Password))
	if err != nil {
		return nil, fmt.Errorf("Login: %w", err)
	}
	access, refresh, err := generateAccessAndRefreshTokens(key, id)
	if err != nil {
		return nil, fmt.Errorf("generateAccessAndRefreshTokens: %w", err)
	}
	hashedRefresh, err := HashRefreshToken(refresh)
	if err != nil {
		return nil, fmt.Errorf("HashRefreshToken: %w", err)
	}
	err = s.rps.UpdateProfile(ctx, id, hashedRefresh)
	if err != nil {
		return nil, fmt.Errorf("UpdateProfile: %w", err)
	}
	jwtResponse := &model.JWTResponse{
		ID:           id,
		AccessToken:  access,
		RefreshToken: refresh,
	}
	return jwtResponse, nil
}

func (s *ProfileService) SignUp(ctx context.Context, user *model.User) error {
	hashedPassword := hashPassword(user.Password)
	user.Password = hashedPassword
	err := s.rps.CreateProfile(ctx, user)
	if err != nil {
		logrus.WithFields(logrus.Fields{"user": user}).Errorf("CreateProfile: %v", err)
		return fmt.Errorf("CreateProfile: %w", err)
	}
	return nil
}

// GenerateAccessAndRefreshTokens func returns access & refresh tokens
func generateAccessAndRefreshTokens(key string, id uuid.UUID) (access, refresh string, err error) {
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(accessTokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
			Id:        id.String(),
		},
	})

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(refreshTokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
			Id:        id.String(),
		},
	})

	access, err = accessToken.SignedString([]byte(key))
	if err != nil {
		return "", "", fmt.Errorf("SignedString(access): %w", err)
	}
	refresh, err = refreshToken.SignedString([]byte(key))
	if err != nil {
		return "", "", fmt.Errorf("SignedString(refresh): %w", err)
	}
	return access, refresh, err
}

// HashPassword func returns hashed password using bcrypt algorithm
func hashPassword(password []byte) []byte {
	hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		return nil
	}
	return hashedPassword
}

// HashRefreshToken func returns hashed refresh token using bcrypt algorithm
func HashRefreshToken(refreshToken string) ([]byte, error) {
	hash := sha256.New()
	hash.Write([]byte(refreshToken))
	hashBytes := hash.Sum(nil)
	hashString := hex.EncodeToString(hashBytes)

	return []byte(hashString), nil
}
