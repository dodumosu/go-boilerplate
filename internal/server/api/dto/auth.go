package dto

type PasswordSignInRequestBody struct {
	Identifier string `json:"identifier"`
	Password   string `json:"password"`
}

type PasswordSignInRequest struct {
	BaseRequest
	Body PasswordSignInRequestBody
}

type AccountSignUpRequestBody struct {
	Email     string `json:"email" format:"email" doc:"Email address"`
	FirstName string `json:"firstName" doc:"First name"`
	LastName  string `json:"lastName" doc:"Last name"`
	Password  string `json:"password" doc:"Password"`
	Username  string `json:"username" doc:"Username"`
}

type AccountSignUpRequest struct {
	BaseRequest
	Body AccountSignUpRequestBody
}

type AccountSignUpResponse struct {
	BaseResponse
}

type AccountLite struct {
	ID        string `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Username  string `json:"username"`
	Email     string `json:"email"`
}

type SignInResponseBody struct {
	BaseResponse
	Account     AccountLite `json:"account"`
	AccessToken string      `json:"accessToken"`
}

type SignInReponse struct {
	Body SignInResponseBody
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
	Body ResendVerificationRequestBody
}

type AccountVerificationRequest struct {
	BaseRequest
	Token string `query:"token" validate:"required"`
}

type AccountVerificationResponse struct {
	Body BaseResponse
}

type SignOutRequest struct {
	Body BaseRequest
}

type SignOutResponse struct {
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
