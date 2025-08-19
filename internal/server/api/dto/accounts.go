package dto

type AccountSignUpRequest struct {
	BaseRequest
	Body AccountSignUpRequestBody
}

type AccountSignUpResponseBody struct {
	BaseResponse
}
type AccountSignUpResponse struct {
	Body AccountSignUpResponseBody
}

type ResendVerificationRequestBody struct {
	Email string `json:"email" format:"email"`
}

type ResendVerificationRequest struct {
	BaseRequest
	Body ResendVerificationRequestBody
}

type ResendVerificationResponseBody struct {
	BaseResponse
}

type ResendVerificationResponse struct {
	Body ResendVerificationResponseBody
}

type AccountVerificationRequest struct {
	BaseRequest
	Token string `query:"token" validate:"required"`
}

type AccountVerificationResponse struct {
	Body BaseResponse
}

type PasswordChangeRequestBody struct {
	OldPassword string `json:"oldPassword"`
	NewPassword string `json:"newPassword"`
}

type PasswordChangeRequest struct {
	BaseRequest
	Body PasswordChangeRequestBody
}

type PasswordChangeResponse struct {
	Body BaseResponse
}

type ForgotPasswordRequestBody struct {
	Email string `json:"email" format:"email"`
}

type ResetPasswordRequestBody struct {
	Email string `json:"email" format:"email"`
}

type CheckAccountAvailabilityRequestBody struct {
	Email    string `json:"email" format:"email"`
	Username string `json:"username"`
}

type CheckAccountAvailabilityRequest struct {
	BaseRequest
	Body CheckAccountAvailabilityRequestBody
}

type CheckAccountAvailabilityResponseBody struct {
	AccountAvailable bool `json:"isAccountAvailable"`
	BaseResponse
}

type CheckAccountAvailabilityResponse struct {
	Body CheckAccountAvailabilityResponseBody
}
