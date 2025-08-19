package handlers

import (
	"context"
	"encoding/json"
	"go-boilerplate/internal/contracts"
	"go-boilerplate/internal/job"

	"github.com/hibiken/asynq"
)

func HandlePasswordResetEmail(ctx context.Context, handlerCtx *job.JobHandlerContext, task *asynq.Task) error {
	var p contracts.PasswordResetEmailPayload
	if err := json.Unmarshal(task.Payload(), &p); err != nil {
		return err
	}

	handlerCtx.EmailService.Logger.Info("sending password reset email", "correlation_id", p.CorrelationID)

	return handlerCtx.EmailService.SendPasswordResetEmail(p.Email, p.Name, p.ResetLink)
}

func HandleVerificationEmail(ctx context.Context, handlerCtx *job.JobHandlerContext, task *asynq.Task) error {
	var p contracts.VerificationEmailPayload
	if err := json.Unmarshal(task.Payload(), &p); err != nil {
		return err
	}

	handlerCtx.EmailService.Logger.Info("sending verification email", "correlation_id", p.CorrelationID)

	return handlerCtx.EmailService.SendWelcomeEmail(p.Email, p.Name, p.VerificationLink)
}
