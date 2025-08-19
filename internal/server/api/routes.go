package api

import (
	"fmt"
	"go-boilerplate/internal/config"
	"net/url"
)

const (
	APIBase   = "/api"
	APIV1Base = APIBase + "/v1"
	APIDocs   = APIBase + "/docs"
)

const (
	AccountBase               = APIV1Base + "/accounts"
	AccountSignUp             = AccountBase + "/signup"
	AccountSignIn             = AccountBase + "/signin"
	AccountSignOut            = AccountBase + "/signout"
	AccountVerify             = AccountBase + "/verify"
	AccountResendVerification = AccountBase + "/resend-verification"
	AccountForgotPassword     = AccountBase + "/forgot-password"
	AccountResetPassword      = AccountBase + "/reset-password"
	AccountChangePassword     = AccountBase + "/change-password"
	AccountAvailability       = AccountBase + "/check-availability"
)

const (
	FrontendAuthVerifyPath        = "/verify-account" // Example frontend path
	FrontendAuthResetPasswordPath = "/reset-password" // Example frontend path
)

// BuildAuthVerifyURL constructs the full URL for email verification.
func BuildAuthVerifyURL(configProvider config.RoutingConfigProvider, verificationToken string) string {
	frontendBaseURL := configProvider.GetFrontendBaseURL() // Use new method
	queryParams := url.Values{}
	queryParams.Set("token", verificationToken)
	// Use a frontend-specific path, not the API path AuthVerify
	return fmt.Sprintf("%s%s?%s", frontendBaseURL, FrontendAuthVerifyPath, queryParams.Encode())
}

// BuildPasswordResetURL constructs the full URL for password reset.
func BuildPasswordResetURL(configProvider config.RoutingConfigProvider, resetToken string) string {
	frontendBaseURL := configProvider.GetFrontendBaseURL() // Use new method
	queryParams := url.Values{}
	queryParams.Set("token", resetToken)
	// Use a frontend-specific path, not the API path AuthResetPassword
	return fmt.Sprintf("%s%s?%s", frontendBaseURL, FrontendAuthResetPasswordPath, queryParams.Encode())
}
