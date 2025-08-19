package api

import (
	"context"
	"errors"
	"go-boilerplate/internal/db/repositories"
	"go-boilerplate/internal/lib"
	"go-boilerplate/internal/server/api/dto"
	"strings"
	"time"

	"github.com/danielgtaylor/huma/v2"
)

func (w *APIWrapper) SignIn(ctx context.Context, input *dto.PasswordSignInRequest) (*dto.SignInReponse, error) {
	accountInfo, err := w.authRepo.SigninWithPassword(ctx, input.Body.Identifier, input.Body.Password)
	if err != nil {
		correlationID := GetCorrelationID(ctx)
		if errors.Is(err, repositories.ErrAccountDoesNotExist) {
			w.logger.Error("sign in failed (invalid credentials)", "correlation_id", correlationID)
			return nil, huma.Error401Unauthorized("Invalid credentials")
		}

		if errors.Is(err, repositories.ErrAccountIsNotActive) {
			return nil, huma.Error403Forbidden("Account is inactive")
		}

		if errors.Is(err, repositories.ErrInvalidCredentials) {
			return nil, huma.Error401Unauthorized("Invalid credentials")
		}

		w.logger.Error("error occurred while signing in", "correlation_id", correlationID, "error", err)
		return nil, huma.Error500InternalServerError("Unable to sign in")
	}

	if accountInfo.SuspendedUntil.Valid && accountInfo.SuspendedUntil.Time.After(time.Now()) {
		return nil, huma.Error403Forbidden("Account is suspended")
	}

	if accountInfo.BannedAt.Valid && accountInfo.BannedAt.Time.Before(time.Now()) {
		return nil, huma.Error403Forbidden("Account is banned")
	}

	if accountInfo.DeactivateAt.Valid && accountInfo.DeactivateAt.Time.Before(time.Now()) {
		return nil, huma.Error403Forbidden("Account is inactive")
	}

	if !accountInfo.IsActive {
		return nil, huma.Error403Forbidden("Account is inactive")
	}

	if !accountInfo.IsVerified {
		return nil, huma.Error403Forbidden("Account has not been verified")
	}

	if !accountInfo.PasswordHash.Valid {
		return nil, huma.Error403Forbidden("Please log in using social login")
	}

	result, err := lib.VerifyPassword(input.Body.Password, accountInfo.PasswordHash.String)
	if err != nil {
		correlationID := GetCorrelationID(ctx)
		w.logger.Error("password verification failed", "correlation_id", correlationID, "error", err)
		return nil, huma.Error500InternalServerError("Unable to log user in with password")
	}

	if !result {
		return nil, huma.Error401Unauthorized("Invalid credentials")
	}

	accessToken, err := w.authenticator.GenerateAccessToken(accountInfo.ID, accountInfo.IsSuperuser, w.config.Auth.TokenLifetime)
	if err != nil {
		correlationID := GetCorrelationID(ctx)
		w.logger.Error("token generation failed", "correlation_id", correlationID, "error", err)
		return nil, huma.Error500InternalServerError("Unable to log user in with password")
	}

	responseBody := dto.SignInResponseBody{
		Account: dto.AccountLite{
			ID:                      accountInfo.ID,
			Email:                   accountInfo.Email,
			Username:                accountInfo.Username.String,
			FirstName:               accountInfo.FirstName,
			LastName:                accountInfo.LastName,
			PasswordChangeRequested: accountInfo.PasswordChangeOnLogin,
			IsSuperuser:             accountInfo.IsSuperuser,
		},
		AccessToken: accessToken,
	}
	responseBody.Message = "ok"
	responseBody.StatusMessage = dto.Success

	return &dto.SignInReponse{
		Body: responseBody,
	}, nil
}

func (w *APIWrapper) SignOut(ctx context.Context, input *struct{}) (*dto.SignOutResponse, error) {
	request, ok := GetRequest(ctx)

	if !ok {
		return nil, huma.Error500InternalServerError("A server error occured")
	}

	authHeader := strings.TrimSpace(request.Header.Get("Authorization"))
	if authHeader == "" {
		return nil, huma.Error400BadRequest("Invalid request sent")
	}
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return nil, huma.Error400BadRequest("Invalid auth header")
	}
	tokenStr := parts[1]

	claims, err := w.authenticator.ValidateAndParseAccessToken(tokenStr)
	if err != nil {
		return nil, huma.Error400BadRequest("Invalid auth header")
	}

	if claims.ID == "" {
		w.logger.Warn("Logout attempt with token missing JTI, cannot add to blocklist", "user_id", claims.Subject)
	} else {
		if claims.ExpiresAt != nil && claims.ExpiresAt.Time.After(time.Now()) {
			err = w.redisService.AddToBlocklist(ctx, claims.ID, claims.ExpiresAt.Time)
			if err != nil {
				w.logger.Error("Failed to add token to blocklist during logout", "jti", claims.ID, "error", err)
				// Non-critical error for logout itself, but important to log
			}
		} else {
			w.logger.Debug("Token already expired, not adding to blocklist on logout", "jti", claims.ID)
		}
	}

	return &dto.SignOutResponse{
		Body: dto.BaseResponse{
			Message:       "ok",
			StatusMessage: dto.Success,
		},
	}, nil
}
