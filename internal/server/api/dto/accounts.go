package dto

import "time"

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

type Account struct {
	FirstName  string    `json:"firstName"`
	LastName   string    `json:"lastName"`
	OtherNames string    `json:"otherNames" omitEmpty:"true"`
	Phone      string    `json:"phone" omitEmpty:"true"`
	Bio        string    `json:"bio" omitEmpty:"true"`
	CreatedAt  time.Time `json:"joinedAt"`
}

type AccountDetailRequest struct {
	BaseRequest
}

type AccountDetailResponseBody struct {
	Account Account `json:"profile"`
	BaseResponse
}

type AccountDetailResponse struct {
	Body AccountDetailResponseBody
}

type AccountUpdateRequestBody struct {
	FirstName  string `json:"firstName" nullable:"true" omitEmpty:"true" required:"false"`
	LastName   string `json:"lastName" nullable:"true" omitEmpty:"true" required:"false"`
	OtherNames string `json:"otherNames" nullable:"true" omitEmpty:"true" required:"false"`
	Phone      string `json:"phone" nullable:"true" omitEmpty:"true" required:"false"`
	Bio        string `json:"bio" nullable:"true" omitEmpty:"true" required:"false"`
}

type AccountUpdateRequest struct {
	Body AccountUpdateRequestBody
}

type AccountUpdateResponseBody struct {
	BaseResponse
	Account Account `json:"account"`
}

type AccountUpdateResponse struct {
	Body AccountUpdateResponseBody
}
