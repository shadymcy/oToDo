package bll

import (
	"crypto/sha256"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/yzx9/otodo/dal"
	"github.com/yzx9/otodo/entity"
)

// Config
// TODO configurable
var passwordNonce = []byte("test_nonce")

// User
type CreateUserPayload struct {
	UserName string `json:"user_name"`
	Password string `json:"password"`
	Nickname string `json:"nickname"`
}

func CreateUser(payload CreateUserPayload) (entity.User, error) {
	user, err := dal.InsertUser(entity.User{
		ID:       uuid.New(),
		Name:     payload.UserName,
		Nickname: payload.Nickname,
		Password: GetCryptoPassword(payload.Password),
	})

	// TODO create base todo list

	return user, err
}

func GetUser(userID uuid.UUID) (entity.User, error) {
	return dal.GetUser(userID)
}

// Invalid User Refresh Token

func CreateInvalidUserRefreshToken(userID uuid.UUID, tokenID uuid.UUID) (entity.UserRefreshToken, error) {
	model, err := dal.InsertInvalidUserRefreshToken(entity.UserRefreshToken{
		ID:        uuid.New(),
		UserID:    userID,
		TokenID:   tokenID,
		CreatedAt: time.Now(),
	})
	if err != nil {
		return entity.UserRefreshToken{}, fmt.Errorf("fails to invalid user refresh token, %w", err)
	}

	return model, nil
}

// Verify is it an valid token.
// Note: This func don't check token expire time
func IsValidRefreshToken(userID string, tokenID string) bool {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return false
	}

	tokenUUID, err := uuid.Parse(tokenID)
	if err != nil {
		return false
	}

	return !dal.ExistInvalidUserRefreshToken(userUUID, tokenUUID)
}

// Password
func GetCryptoPassword(password string) []byte {
	pwd := sha256.Sum256(append([]byte(password), passwordNonce...))
	return pwd[:]
}
