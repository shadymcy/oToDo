package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/yzx9/otodo/bll"
	"github.com/yzx9/otodo/utils"
	webUtils "github.com/yzx9/otodo/web/utils"
)

// Ping Test
func GetSessionHandler(c *gin.Context) {
	c.String(http.StatusOK, "hello")
}

// Login
func PostSessionHandler(c *gin.Context) {
	payload := struct {
		UserName string `json:"user_name"`
		Password string `json:"password"`
	}{}
	if err := c.ShouldBind(&payload); err != nil {
		webUtils.AbortWithError(c, err)
		return
	}

	tokens, err := bll.Login(payload.UserName, payload.Password)
	if err != nil {
		webUtils.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, struct {
		AccessToken  string `json:"access_token"`
		TokenType    string `json:"token_type"`
		ExpiresIn    int64  `json:"expires_in"`
		RefreshToken string `json:"refresh_token"`
	}{tokens.AccessToken, tokens.TokenType, tokens.ExpiresIn, tokens.RefreshToken})
}

// Logout, unactive refresh token if exists
func DeleteSessionHandler(c *gin.Context) {
	if token, err := parseRefreshToken(c); err == nil {
		bll.Logout(token)
	}

	c.String(http.StatusOK, "See you!")
}

// Create New Access Token by Refresh Token
func PostAccessTokenHandler(c *gin.Context) {
	token, err := parseRefreshToken(c)
	if err != nil {
		webUtils.AbortWithError(c, err)
		return
	}

	claims, ok := token.Claims.(*utils.TokenClaims)
	if !ok || !token.Valid {
		webUtils.AbortWithJson(c, "invalid token")
		return
	}

	newToken, err := bll.NewAccessToken(claims.UserID)
	if err != nil {
		webUtils.AbortWithJson(c, fmt.Sprintf("fails to refresh an token, %v", err.Error()))
		return
	}

	c.JSON(http.StatusOK, struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		ExpiresIn   int64  `json:"expires_in"`
	}{newToken.AccessToken, newToken.TokenType, newToken.ExpiresIn})
}

func parseRefreshToken(c *gin.Context) (*jwt.Token, error) {
	obj := &struct {
		RefreshToken string `json:"refresh_token"`
	}{}
	if err := c.ShouldBind(&obj); err != nil {
		return nil, err
	}

	token, err := utils.ParseJWT(obj.RefreshToken)
	if err != nil || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return token, nil
}
