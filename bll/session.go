package bll

import (
	"bytes"
	"fmt"
	"regexp"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/yzx9/otodo/dal"
	"github.com/yzx9/otodo/entity"
)

// Config
// TODO configurable
var accessTokenExpiresIn = 15 * time.Minute
var refreshTokenExpiresIn = 15 * 24 * time.Hour
var accessTokenRefreshThreshold = 5 * time.Minute

// Constans
var accessTokenExpiresInSeconds = int64(accessTokenExpiresIn.Seconds())
var tokenType = "Bearer"
var authorizationRegexString = "^[Bb]earer (?P<token>[\\w-]+.[\\w-]+.[\\w-]+)$"
var authorizationRegex = regexp.MustCompile(authorizationRegexString)

type AuthTokenResult struct {
	AccessToken  string
	TokenType    string
	ExpiresIn    int64
	RefreshToken string
}

type SessionTokenClaims struct {
	TokenClaims
	RefreshTokenID string `json:"rti,omitempty"`
	UserNickname   string `json:"nickname,omitempty"`
}

func Login(userName, password string) (AuthTokenResult, error) {
	user, err := dal.GetUserByUserName(userName)
	if err != nil {
		return AuthTokenResult{}, fmt.Errorf("user not found: %v", userName)
	}

	if !bytes.Equal(user.Password, GetCryptoPassword(password)) {
		return AuthTokenResult{}, fmt.Errorf("invalid credential")
	}

	refreshToken, refreshTokenID := newRefreshToken(user, refreshTokenExpiresIn)
	return AuthTokenResult{
		AccessToken:  newAccessToken(user, refreshTokenID, accessTokenExpiresIn),
		TokenType:    tokenType,
		ExpiresIn:    accessTokenExpiresInSeconds,
		RefreshToken: refreshToken,
	}, nil
}

func Logout(userID, refreshTokenID string) error {
	_, err := CreateInvalidUserRefreshToken(userID, refreshTokenID)
	return err
}

func NewAccessToken(userID, refreshTokenID string) (AuthTokenResult, error) {
	user, err := dal.GetUser(userID)
	if err != nil {
		return AuthTokenResult{}, fmt.Errorf("fails to get user, %w", err)
	}

	return AuthTokenResult{
		AccessToken: newAccessToken(user, refreshTokenID, accessTokenExpiresIn),
		TokenType:   tokenType,
		ExpiresIn:   accessTokenExpiresInSeconds,
	}, nil
}

func ParseSessionToken(token string) (*jwt.Token, error) {
	return ParseToken(token, &SessionTokenClaims{})
}

func ParseAccessToken(authorization string) (*jwt.Token, error) {
	matches := authorizationRegex.FindStringSubmatch(authorization)
	if len(matches) != 2 {
		return nil, fmt.Errorf("unauthorized")
	}

	token, err := ParseToken(matches[1], &SessionTokenClaims{})
	if err != nil {
		return nil, fmt.Errorf("fails to parse access token: %w", err)
	}

	claims, ok := token.Claims.(*SessionTokenClaims)
	_, err = uuid.Parse(claims.UserID)
	if !ok || err != nil {
		return nil, fmt.Errorf("invalid access token")
	}

	return token, nil
}

func ShouldRefreshAccessToken(oldAccessToken *jwt.Token) bool {
	if !oldAccessToken.Valid {
		return false
	}

	claims, ok := oldAccessToken.Claims.(*SessionTokenClaims)
	if !ok || claims.ExpiresAt == 0 {
		return false
	}

	return time.Now().Add(accessTokenRefreshThreshold).Unix() > claims.ExpiresAt
}

func newAccessToken(user entity.User, refreshTokenID string, exp time.Duration) string {
	claims := SessionTokenClaims{
		TokenClaims:    NewClaims(user.ID, exp),
		UserNickname:   user.Nickname,
		RefreshTokenID: refreshTokenID,
	}
	return NewToken(claims)
}

func newRefreshToken(user entity.User, exp time.Duration) (string, string) {
	id := uuid.NewString()
	claims := SessionTokenClaims{TokenClaims: NewClaims(user.ID, exp)}
	claims.Id = id
	return NewToken(claims), id
}
