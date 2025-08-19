package api

import (
	"context"
	"encoding/json"
	"errors"
	"go-boilerplate/internal/contracts"
	"go-boilerplate/internal/db/repositories"
	"go-boilerplate/internal/lib"
	"go-boilerplate/internal/server/api/dto"

	"github.com/danielgtaylor/huma/v2"
	"github.com/hibiken/asynq"
)

func (w *APIWrapper) SignUp(ctx context.Context, input *dto.AccountSignUpRequest) (*dto.AccountSignUpResponse, error) {
	accountCreateParams := repositories.AccountCreateParams{
		Email:                 input.Body.Email,
		FirstName:             input.Body.FirstName,
		IsOAuth:               false,
		LastName:              input.Body.LastName,
		Password:              input.Body.Password,
		Username:              input.Body.Username,
		ChangePasswordOnLogin: false,
		Activate:              true,
	}

	user, profile, err := w.accountRepo.CreateAccount(ctx, accountCreateParams)
	if err != nil {
		// Handle user creation error, maybe the user already exists
		// For now, we return a generic error
		return nil, err
	}

	if !user.IsVerified {
		verificationToken, tokenErr := w.authenticator.GenerateVerificationTokenJWT(user.ID, w.config.Auth.VerificationTokenLifetime)
		if tokenErr != nil {
			w.logger.Error("error generating verification token", "error", tokenErr)

			responseBody := dto.AccountSignUpResponseBody{}
			responseBody.Message = "Account created, but verification email not sent. Please request a resend."
			responseBody.StatusMessage = dto.Fail

			return &dto.AccountSignUpResponse{Body: responseBody}, nil
		}

		verificationLink := BuildAuthVerifyURL(w.config.Routing, verificationToken)
		name := profile.FirstName

		// Create and enqueue a new verification email task.
		w.logger.Info("Sending verification email...")
		correlationID := GetCorrelationID(ctx)
		taskPayload := contracts.VerificationEmailPayload{
			Email:            user.Email,
			Name:             name,
			VerificationLink: verificationLink,
			CorrelationID:    correlationID,
		}
		payload, err := json.Marshal(taskPayload)
		if err != nil {
			return nil, huma.Error500InternalServerError("Unable to send verification email")
		}

		task := asynq.NewTask(contracts.TypeVerificationEmail, payload)
		_, err = w.jobService.Enqueue(ctx, task)
		if err != nil {
			return nil, huma.Error500InternalServerError("Unable to send verification email")
		}
	}

	responseBody := dto.AccountSignUpResponseBody{}
	if user.IsVerified {
		responseBody.Message = "Account created"
	} else {
		responseBody.Message = "Account created. Please check your email to verify your account."
	}
	responseBody.StatusMessage = dto.Success

	return &dto.AccountSignUpResponse{Body: responseBody}, nil
}

func (w *APIWrapper) VerifyAccount(ctx context.Context, input *dto.AccountVerificationRequest) (*dto.AccountVerificationResponse, error) {
	userID, err := w.authenticator.ValidateAndParseVerificationTokenJWT(input.Token)
	if err != nil {
		return nil, huma.Error400BadRequest("Invalid or expired verification token")
	}

	err = w.accountRepo.VerifyAccount(ctx, userID)
	if errors.Is(err, repositories.ErrAccountDoesNotExist) {
		return nil, huma.Error404NotFound("Account does not exist")
	}

	return &dto.AccountVerificationResponse{
		Body: dto.BaseResponse{
			Message:       "Account verified successfully",
			StatusMessage: dto.Success,
		},
	}, nil
}

func (w *APIWrapper) CheckAccountAvailability(ctx context.Context, input *dto.CheckAccountAvailabilityRequest) (*dto.CheckAccountAvailabilityResponse, error) {
	result, err := w.accountRepo.CheckAccountAvailability(ctx, input.Body.Username, input.Body.Email)
	if err != nil {
		return nil, huma.Error500InternalServerError("Account check failed")
	}

	responseBody := dto.CheckAccountAvailabilityResponseBody{}
	responseBody.AccountAvailable = result
	if result {
		responseBody.Message = "Account already exists"
	} else {
		responseBody.Message = "Account does not exist"
	}
	responseBody.StatusMessage = dto.Success

	return &dto.CheckAccountAvailabilityResponse{
		Body: responseBody,
	}, nil
}

func (w *APIWrapper) ResendVerificationEmail(ctx context.Context, input *dto.ResendVerificationRequest) (*dto.ResendVerificationResponse, error) {
	accountInfo, err := w.accountRepo.GetByEmail(ctx, input.Body.Email)
	if err != nil {
		if errors.Is(err, repositories.ErrAccountDoesNotExist) {
			return nil, huma.Error404NotFound("Account not found")
		}
	}

	if !accountInfo.IsVerified {
		verificationToken, tokenErr := w.authenticator.GenerateVerificationTokenJWT(accountInfo.ID, w.config.Auth.TokenLifetime)
		if tokenErr != nil {
			return nil, huma.Error500InternalServerError("Unable to resend email")
		}

		verificationLink := BuildAuthVerifyURL(w.config.Routing, verificationToken)
		correlationID := GetCorrelationID(ctx)
		taskPayload := contracts.VerificationEmailPayload{
			Email:            accountInfo.Email,
			Name:             accountInfo.FirstName,
			VerificationLink: verificationLink,
			CorrelationID:    correlationID,
		}
		payload, err := json.Marshal(taskPayload)
		if err != nil {
			return nil, huma.Error500InternalServerError("Unable to resend email")
		}

		task := asynq.NewTask(contracts.TypeVerificationEmail, payload)
		_, err = w.jobService.Enqueue(ctx, task)
		if err != nil {
			return nil, huma.Error500InternalServerError("Unable to resend email")
		}
	} else {
		return nil, huma.Error409Conflict("Account already verified")
	}

	responseBody := dto.ResendVerificationResponseBody{}
	responseBody.Message = "Please check your email for the verification email"
	responseBody.StatusMessage = dto.Success

	return &dto.ResendVerificationResponse{
		Body: responseBody,
	}, nil
}

func (w *APIWrapper) UpdatePassword(ctx context.Context, input *dto.PasswordChangeRequest) (*dto.PasswordChangeResponse, error) {
	userID, ok := GetUserID(ctx)
	if !ok {
		return nil, huma.Error401Unauthorized("Please provide credentials to continue")
	}

	user, err := w.accountRepo.GetUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, repositories.ErrAccountDoesNotExist) {
			return nil, huma.Error404NotFound("Account not found")
		}

		return nil, huma.Error500InternalServerError("Unable to update password")
	}

	result, err := lib.VerifyPassword(input.Body.OldPassword, user.PasswordHash.String)
	if err != nil {
		return nil, huma.Error500InternalServerError("Unable to update password")
	}
	if !result {
		return nil, huma.Error401Unauthorized("Please ensure that your current password matches")
	}

	newPasswordHash, err := lib.HashPassword(input.Body.NewPassword)
	if err != nil {
		return nil, huma.Error500InternalServerError("Unable to update password")
	}

	err = w.accountRepo.UpdatePassword(ctx, userID, newPasswordHash)
	if err != nil {
		if errors.Is(err, repositories.ErrAccountDoesNotExist) {
			return nil, huma.Error404NotFound("Account was not found")
		}

		return nil, huma.Error500InternalServerError("Unable to update password")
	}

	return &dto.PasswordChangeResponse{
		Body: dto.BaseResponse{
			Message:       "Password updated",
			StatusMessage: dto.Success,
		},
	}, nil
}
