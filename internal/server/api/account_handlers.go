package api

import (
	"context"
	"encoding/json"
	"go-boilerplate/internal/contracts"
	"go-boilerplate/internal/db/repositories"
	"go-boilerplate/internal/server/api/dto"

	"github.com/hibiken/asynq"
)

func (w *APIWrapper) SignUp(ctx context.Context, input *dto.AccountSignUpRequest) (*dto.AccountSignUpResponse, error) {
	accountCreateParams := repositories.AccountCreateParams{
		Email:     input.Body.Email,
		FirstName: input.Body.FirstName,
		IsOAuth:   false,
		LastName:  input.Body.LastName,
		Password:  input.Body.Password,
		Username:  input.Body.Username,
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
		correlationID, _ := ctx.Value(CorrelationIDKey).(string)
		taskPayload := contracts.VerificationEmailPayload{
			Email:            user.Email,
			Name:             name,
			VerificationLink: verificationLink,
			CorrelationID:    correlationID,
		}
		payload, err := json.Marshal(taskPayload)
		if err != nil {
			return nil, err
		}

		task := asynq.NewTask(contracts.TypeVerificationEmail, payload)
		_, err = w.jobService.Enqueue(ctx, task)
		if err != nil {
			return nil, err
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
