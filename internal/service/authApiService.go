// Package service implements the service interface for the given service
package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	middlewr "github.com/eugenshima/trading-api/internal/middleware"
	"github.com/eugenshima/trading-api/internal/model"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

// ProfileService represents a service that is responsible for managing profiles
type ProfileService struct {
	rps ProfileAPIRepository
}

// NewProfileService creates a new ProfileService
func NewProfileService(rps ProfileAPIRepository) *ProfileService {
	return &ProfileService{rps: rps}
}

// ProfileAPIRepository represents a profile-api-repository
type ProfileAPIRepository interface {
	CreateProfile(context.Context, *model.User) error
	Login(context.Context, string, []byte) (uuid.UUID, error)
	UpdateProfile(context.Context, uuid.UUID, []byte) error
	GetRefreshTokenByID(ctx context.Context, ID uuid.UUID) ([]byte, error)
	DeleteProfileByID(context.Context, uuid.UUID) error
}

const (
	key             = "ew4t137tr1eyfg1ryg4ryerg2743gr2"
	accessTokenTTL  = 60 * time.Minute
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
		RefreshToken: []byte(refresh),
	}
	return jwtResponse, nil
}

// SignUp method to sign up a user with the provided credentials
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

// RefreshTokenPair method to refresh access and refresh tokens
func (s *ProfileService) RefreshTokenPair(ctx context.Context, ID uuid.UUID, access string, refresh []byte) (*model.JWTResponse, error) {
	savedRefreshToken, err := s.rps.GetRefreshTokenByID(ctx, ID)
	if err != nil {
		return nil, fmt.Errorf("GetRefreshTokenByID: %w", err)
	}
	hashedRefreshToken, err := HashRefreshToken(string(refresh))
	if err != nil {
		return nil, fmt.Errorf("HashRefreshToken: %w", err)
	}
	if !CompareHashedTokens(savedRefreshToken, hashedRefreshToken) {
		return nil, fmt.Errorf("CompareHashedTokens: %w", err)
	}
	tokenID, err := middlewr.GetPayloadFromToken(access)
	if err != nil {
		return nil, fmt.Errorf("GetPayloadFromToken: %w", err)
	}
	compID, err := CompareTokenIDs(access, string(refresh), key)
	if err != nil {
		return nil, fmt.Errorf("CompareTokenIDs: %w", err)
	}
	if !compID {
		return nil, fmt.Errorf("invalid token(campare error): %w", err)
	}
	// GenerateAccessAndRefreshTokens
	newAccess, newRefresh, err := generateAccessAndRefreshTokens(key, tokenID)
	if err != nil {
		return nil, fmt.Errorf("GenerateAccessAndRefreshTokens: %w", err)
	}
	hashedRefreshToken, err = HashRefreshToken(newRefresh)
	if err != nil {
		return nil, fmt.Errorf("HashRefreshToken: %w", err)
	}
	err = s.rps.UpdateProfile(ctx, ID, hashedRefreshToken)
	if err != nil {
		return nil, fmt.Errorf("UpdateProfile: %w", err)
	}
	jwtResponse := &model.JWTResponse{
		ID:           ID,
		AccessToken:  newAccess,
		RefreshToken: []byte(newRefresh),
	}
	return jwtResponse, nil
}

// DeleteProfile function deletes the profile associated with the given ID
func (s *ProfileService) DeleteProfile(ctx context.Context, ID uuid.UUID) error {
	return s.rps.DeleteProfileByID(ctx, ID)
}

// generateAccessAndRefreshTokens func returns generated access & refresh tokens
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

// CompareHashedTokens func compairs hashed tokens from database and request
func CompareHashedTokens(token1, token2 []byte) bool {
	return sha256.Sum256(token1) == sha256.Sum256(token2)
}

// CompareTokenIDs func compares token ids
func CompareTokenIDs(accessToken, refreshToken, key string) (bool, error) {
	accessID, err := ExtractIDFromToken(accessToken, key)
	if err != nil {
		return false, fmt.Errorf("ExtractIDFromToken: %w", err)
	}

	refreshID, err := ExtractIDFromToken(refreshToken, key)
	if err != nil {
		return false, fmt.Errorf("ExtractIDFromToken: %w", err)
	}
	return accessID == refreshID, nil
}

// ExtractIDFromToken extracts the identifier (ID) from the payload (claims) of the token.
func ExtractIDFromToken(tokenString, key string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(key), nil
	})
	if err != nil {
		return "", fmt.Errorf("Parse(): %w", err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if id, ok := claims["jti"].(string); ok {
			return id, nil
		}
	}

	return "", fmt.Errorf("error extracting ID from token: %v", token)
}
