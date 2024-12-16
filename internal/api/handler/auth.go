package handler

import (
	"auth_rest/internal/models"
	"auth_rest/internal/services"
	"auth_rest/internal/utils"
	"encoding/json"
	"net/http"
)

// Login godoc
// @Summary Login
// @Security ApiKeyAuth
// @Tags songs
// @Description authorization + token gen
// @ID user-login
// @Accept  json
// @Produce  json
// @Param input body models.User true "user info"
// @Success 200 {object} services.TokenResponse "generated tokens"
// @Failure 400,404 {object} models.ErrorResponse "bad Request"
// @Failure 500 {object} models.ErrorResponse "internal server error"
// @Router /auth/login/ [post]
func Login(service services.AuthService, ctx *utils.AppContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var input models.User

		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			http.Error(w, "invalid JSON input", http.StatusBadRequest)
			ctx.Logger.Error(err)
			return
		}

		tokensResp, err := service.CreateSession(input)
		if err != nil {
			http.Error(w, "failed to create session", http.StatusInternalServerError)
			ctx.Logger.Error(err, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		response := map[string]interface{}{
			"access_token":  tokensResp.AccessToken,
			"refresh_token": tokensResp.RefreshToken,
		}
		if err = json.NewEncoder(w).Encode(response); err != nil {
			ctx.Logger.Error(err, http.StatusInternalServerError)
			return
		}
	}
}

// Refresh godoc
// @Summary Refresh
// @Security ApiKeyAuth
// @Tags songs
// @Description refresh tokens
// @ID refresh-tokens
// @Accept  json
// @Produce  json
// @Param input body models.RefreshTokenRequest true "user info + refresh token"
// @Success 200 {object} services.TokenResponse "generated tokens"
// @Failure 400,404 {object} models.ErrorResponse "bad Request"
// @Failure 500 {object} models.ErrorResponse "internal server error"
// @Router /auth/refresh/ [post]
func Refresh(service services.AuthService, ctx *utils.AppContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var input models.RefreshTokenRequest

		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			ctx.Logger.Error(err, http.StatusBadRequest)
			return
		}

		token, err := service.RefreshTokens(input)
		if err != nil {
			ctx.Logger.Error(err, http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		response := map[string]interface{}{
			"token": token,
		}
		if err = json.NewEncoder(w).Encode(response); err != nil {
			ctx.Logger.Error(err, http.StatusBadRequest)
			return
		}
	}
}
