package auth

import "github.com/gin-gonic/gin"

// Saves the refresh token to a cookie using the gin context.
func SaveRefreshTokenCookie(c *gin.Context, refreshToken string) {
	c.SetCookie("refresh_token", refreshToken, MaxRefreshCookieMaxAge, "/api/auth", "localhost", true, true)
}
