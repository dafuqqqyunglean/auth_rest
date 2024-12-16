package services

import (
	"auth_rest/internal/lib/tokens"
	"auth_rest/internal/models"
	"auth_rest/internal/storage/sql"
	"auth_rest/internal/utils"
	"crypto/sha256"
	"encoding/base64"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type AuthService interface {
	CreateSession(user models.User) (TokenResponse, error)
	RefreshTokens(request models.RefreshTokenRequest) (TokenResponse, error)
	generateTokens(user models.User) (TokenResponse, error)
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type Auth struct {
	repo       sql.AuthRepository
	tokens     tokens.TokenManager
	ctx        *utils.AppContext
	AccessTTL  time.Duration
	RefreshTTL time.Duration
}

func NewAuthService(repo sql.AuthRepository, tokens tokens.TokenManager, ctx *utils.AppContext, accessTTL time.Duration, refreshTTL time.Duration) *Auth {
	return &Auth{
		repo:       repo,
		ctx:        ctx,
		tokens:     tokens,
		AccessTTL:  accessTTL,
		RefreshTTL: refreshTTL,
	}
}

func (a *Auth) CreateSession(user models.User) (TokenResponse, error) {
	a.ctx.Logger.Info("generating session")

	tokenResp, err := a.generateTokens(user)
	if err != nil {
		a.ctx.Logger.Error("failed to generate tokens")
		return TokenResponse{}, err
	}

	hashToken := sha256.Sum256([]byte(tokenResp.RefreshToken))

	hashedRefreshToken, err := bcrypt.GenerateFromPassword(hashToken[:], bcrypt.DefaultCost)
	if err != nil {
		a.ctx.Logger.Error("failed to hash refresh token")
		return TokenResponse{}, err
	}

	err = a.repo.PostRefreshToken(user.GUID, string(hashedRefreshToken))
	if err != nil {
		a.ctx.Logger.Error("failed to set refresh token")
		return TokenResponse{}, err
	}

	tokenResp.RefreshToken = base64.StdEncoding.EncodeToString([]byte(tokenResp.RefreshToken))

	a.ctx.Logger.Info("session created")

	return tokenResp, nil
}

func (a *Auth) RefreshTokens(request models.RefreshTokenRequest) (TokenResponse, error) {
	a.ctx.Logger.Info("refreshing tokens")

	decodedToken, err := base64.StdEncoding.DecodeString(request.RefreshToken)
	if err != nil {
		return TokenResponse{}, err
	}

	claims, err := a.tokens.ParseRefresh(string(decodedToken))
	if err != nil {
		return TokenResponse{}, err
	}
	hashToken := sha256.Sum256(decodedToken)

	dbHashToken, err := a.repo.GetRefreshToken(claims.GUID)
	if err != nil {
		return TokenResponse{}, err
	}

	if err = bcrypt.CompareHashAndPassword([]byte(dbHashToken), hashToken[:]); err != nil {
		a.ctx.Logger.Error("invalid token for guid")
		return TokenResponse{}, err
	}

	if claims.IP != request.User.IP {
		a.ctx.Logger.Warn("ip changed")
	}

	err = a.repo.DeleteRefreshToken(request.User.GUID)
	if err != nil {
		return TokenResponse{}, err
	}

	tokensResp, err := a.CreateSession(request.User)
	if err != nil {
		a.ctx.Logger.Error("failed to generate tokens")
		return TokenResponse{}, err
	}

	return tokensResp, nil
}

func (a *Auth) generateTokens(user models.User) (TokenResponse, error) {
	var (
		result TokenResponse
		err    error
	)

	result.AccessToken, err = a.tokens.NewAccess(user, a.AccessTTL)
	if err != nil {
		a.ctx.Logger.Error("error occurred while creating JWT token")
		return result, err
	}

	result.RefreshToken, err = a.tokens.NewRefresh(user, a.RefreshTTL)
	if err != nil {
		a.ctx.Logger.Error("error occurred while creating refresh token")
		return result, err
	}

	return result, nil
}
