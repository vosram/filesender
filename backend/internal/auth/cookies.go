package auth

import "github.com/gin-gonic/gin"

const (
	RefreshTokenCookieName string = "refresh_token"
)

// Saves the refresh token to a cookie using the gin context.
func SaveRefreshTokenCookie(c *gin.Context, refreshToken string) {
	c.SetCookie(RefreshTokenCookieName, refreshToken, MaxRefreshCookieMaxAge, "/api/auth", "localhost", true, true)
}

// Get refresh token from cookie in gin context.
func GetRefreshTokenFromCookie(c *gin.Context) (string, error) {
	return c.Cookie(RefreshTokenCookieName)
}
