package bll

import (
	"bytes"
	"fmt"
	"regexp"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/yzx9/otodo/dal"
	"github.com/yzx9/otodo/model/dto"
	"github.com/yzx9/otodo/model/entity"
	"github.com/yzx9/otodo/otodo"
	"github.com/yzx9/otodo/util"
)

const tokenType = `bearer`
const authorizationRegexString = `^[Bb]earer (?P<token>[\w-]+.[\w-]+.[\w-]+)$`

var authorizationRegex = regexp.MustCompile(authorizationRegexString)

func Login(userName, password string) (dto.SessionToken, error) {
	write := func() (dto.SessionToken, error) {
		return dto.SessionToken{}, util.NewErrorWithBadRequest("invalid credential")
	}

	user, err := dal.SelectUserByUserName(userName)
	if err != nil || user.Password == nil {
		return write()
	}

	if cryptoPwd := GetCryptoPassword(password); !bytes.Equal(user.Password, cryptoPwd) {
		return write()
	}

	return newSessionToken(user), nil
}

func LoginByGithubOAuth(code, state string) (dto.SessionToken, error) {
	token, err := FetchGithubOAuthToken(code, state)
	if err != nil {
		return dto.SessionToken{}, fmt.Errorf("fails to login: %w", err)
	}

	profile, err := FetchGithubUserPublicProfile(token.Token)
	if err != nil {
		return dto.SessionToken{}, fmt.Errorf("fails to fetch github user: %w", err)
	}

	user, err := getOrRegisterUserByGithub(profile)
	if err != nil {
		return dto.SessionToken{}, fmt.Errorf("fails to get user: %w", err)
	}

	go UpdateThirdPartyOAuthTokenAsync(&token)

	return newSessionToken(user), nil
}

func Logout(userID int64, refreshTokenID string) error {
	_, err := CreateUserInvalidRefreshToken(userID, refreshTokenID)
	return err
}

func NewAccessToken(userID int64, refreshTokenID string) (dto.SessionToken, error) {
	user, err := dal.SelectUser(userID)
	if err != nil {
		return dto.SessionToken{}, fmt.Errorf("fails to get user, %w", err)
	}

	return newAccessToken(user, refreshTokenID), nil
}

func ParseSessionToken(token string) (*jwt.Token, error) {
	return ParseToken(token, &dto.SessionTokenClaims{})
}

func ParseAccessToken(authorization string) (*jwt.Token, error) {
	matches := authorizationRegex.FindStringSubmatch(authorization)
	if len(matches) != 2 {
		return nil, fmt.Errorf("unauthorized")
	}

	token, err := ParseToken(matches[1], &dto.SessionTokenClaims{})
	if err != nil {
		return nil, fmt.Errorf("fails to parse access token: %w", err)
	}

	return token, nil
}

func ShouldRefreshAccessToken(oldAccessToken *jwt.Token) bool {
	if !oldAccessToken.Valid {
		return false
	}

	claims, ok := oldAccessToken.Claims.(*dto.SessionTokenClaims)
	if !ok || claims.ExpiresAt == 0 {
		return false
	}

	thd := otodo.Conf.Session.AccessTokenRefreshThreshold
	dur := time.Duration(thd * int(time.Second))
	return time.Now().Add(dur).Unix() > claims.ExpiresAt
}

// access token only
func newAccessToken(user entity.User, refreshTokenID string) dto.SessionToken {
	exp := otodo.Conf.Session.AccessTokenExpiresIn
	dur := time.Duration(exp * int(time.Second))

	claims := dto.SessionTokenClaims{
		TokenClaims:    NewClaims(user.ID, dur),
		RefreshTokenID: refreshTokenID,
	}
	token := NewToken(claims)

	return dto.SessionToken{
		AccessToken: token,
		TokenType:   tokenType,
		ExpiresIn:   int64(exp),
	}
}

// access token + refresh token
func newSessionToken(user entity.User) dto.SessionToken {
	// refresh token
	exp := otodo.Conf.Session.RefreshTokenExpiresIn
	dur := time.Duration(exp * int(time.Second))

	claims := dto.SessionTokenClaims{TokenClaims: NewClaims(user.ID, dur)}
	claims.Id = uuid.NewString()
	refreshToken := NewToken(claims)
	refreshTokenID := claims.Id

	// access token
	re := newAccessToken(user, refreshTokenID)

	re.RefreshToken = refreshToken
	return re
}
