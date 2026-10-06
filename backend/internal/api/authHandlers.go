package api

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"math/rand/v2"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/vosram/filesender/backend/internal/auth"
	"github.com/vosram/filesender/backend/internal/database"
)

const fiveMinutes = 5 * time.Minute

// GET /api/auth/email/login
func (api *apiConfig) EmailLogin(c *gin.Context) {
	var reqBody auth.EmailLoginBody
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Request was formatted wrong"})
		return
	}

	token := auth.CreateOpaqueToken()
	tokenHash, err := auth.CreateTokenHash(token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create email token"})
		return
	}

	// TODO: save tokenHash in DB and print token to the console
	err = api.db.SaveEmailLoginToken(c, database.SaveEmailLoginTokenParams{
		TokenHash: tokenHash,
		Email:     reqBody.Email,
		ExpiresAt: pgtype.Timestamptz{Time: time.Now().Add(fiveMinutes), Valid: true},
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Internal Server Error"})
		return
	}

	fmt.Printf("Email Login is token is: %s\n", token)
	c.JSON(http.StatusOK, gin.H{"data": reqBody.Email})
}

// GET /api/auth/email/verify
func (api *apiConfig) EmailLoginVerify(c *gin.Context) {
	var reqBody auth.EmailVerifyBody
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Request was formatted incorrectly",
		})
		return
	}
	tokenHash, _ := auth.CreateTokenHash(reqBody.Token)
	loginToken, err := api.db.ConsumeEmailLoginToken(c, tokenHash)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to find login token", "token": loginToken})
		return
	}
	user, err := api.db.FindUserByEmail(c, loginToken.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// User needs account registration
			newUser, dbErr := api.db.CreateUserAndReturnUser(c, database.CreateUserAndReturnUserParams{
				Name:           fmt.Sprintf("Momo%d", rand.IntN(100000)),
				Email:          loginToken.Email,
				CredentialType: auth.CredentialEmail,
			})
			if dbErr != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "internal server error",
				})
				return
			}
			refreshToken := auth.CreateOpaqueToken()
			refreshTokenHash, err := auth.CreateTokenHash(refreshToken)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create refresh token hash"})
				return
			}
			err = api.db.SaveRefreshToken(c, database.SaveRefreshTokenParams{
				TokenHash: refreshTokenHash,
				UserID:    newUser.ID,
				ExpiresAt: pgtype.Timestamptz{
					Time:  time.Now().Add(auth.MaxRefreshCookieMaxAge * time.Second),
					Valid: true,
				},
			})
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create refresh token session"})
				return
			}
			accessToken, err := auth.CreateAccessJWT(newUser, api.JWTSecret)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create access token"})
				return
			}
			auth.SaveRefreshTokenCookie(c, refreshToken)
			c.JSON(http.StatusOK, gin.H{"accessToken": accessToken})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}

	// User has account
	accessToken, err := auth.CreateAccessJWT(user, api.JWTSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create access token"})
		return
	}
	refreshToken := auth.CreateOpaqueToken()
	refreshTokenHash, err := auth.CreateTokenHash(refreshToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create refresh token hash"})
		return
	}
	err = api.db.SaveRefreshToken(c, database.SaveRefreshTokenParams{
		TokenHash: refreshTokenHash,
		UserID:    user.ID,
		ExpiresAt: pgtype.Timestamptz{
			Time:  time.Now().Add(auth.MaxRefreshCookieMaxAge * time.Second),
			Valid: true,
		},
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save refresh token in db"})
		return
	}
	err = api.db.UpdateLastLoginByUserID(c, user.ID)
	if err != nil {
		log.Printf("failed to update last_login for user: %s", err.Error())
	}
	auth.SaveRefreshTokenCookie(c, refreshToken)
	c.JSON(http.StatusOK, gin.H{"accessToken": accessToken})
}

// POST /api/auth/refresh
func (api *apiConfig) RefreshJWT(c *gin.Context) {
	refreshToken, err := auth.GetRefreshTokenFromCookie(c)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "refresh token missing"})
		return
	}

	refreshTokenHash, err := auth.CreateTokenHash(refreshToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash token"})
		return
	}
	dbToken, err := api.db.ConsumeRefreshToken(c, refreshTokenHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusForbidden, gin.H{"error": "token invalid"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to connect to db"})
		return
	}

	if time.Now().After(dbToken.ExpiresAt.Time) {
		c.JSON(http.StatusForbidden, gin.H{"error": "token expired"})
		return
	}
	user, err := api.db.FindUserById(c, dbToken.UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusForbidden, gin.H{"error": "couldn't find user"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to connect to db"})
		return
	}
	if user.BannedUntil.Valid && time.Now().Before(user.BannedUntil.Time) {
		c.JSON(http.StatusForbidden, gin.H{"error": "user is currently banned", "banned": user.BannedUntil.Time.Format(time.RFC3339)})
		return
	}

	newRefreshToken := auth.CreateOpaqueToken()
	newRefreshTokenHash, err := auth.CreateTokenHash(newRefreshToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash new refresh token"})
		return
	}
	accessToken, err := auth.CreateAccessJWT(user, api.JWTSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create access JWT"})
		return
	}
	err = api.db.SaveRefreshToken(c, database.SaveRefreshTokenParams{
		TokenHash: newRefreshTokenHash,
		UserID:    user.ID,
		ExpiresAt: pgtype.Timestamptz{
			Time:  time.Now().Add(auth.MaxRefreshCookieMaxAge * time.Second),
			Valid: true,
		},
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save new refresh token"})
		return
	}
	auth.SaveRefreshTokenCookie(c, newRefreshToken)
	c.JSON(http.StatusOK, gin.H{"accessToken": accessToken})
}
