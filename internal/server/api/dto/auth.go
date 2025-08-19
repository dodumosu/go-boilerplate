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

type AccountLite struct {
	ID                      string `json:"id"`
	FirstName               string `json:"firstName"`
	LastName                string `json:"lastName"`
	Username                string `json:"username"`
	Email                   string `json:"email"`
	PasswordChangeRequested bool   `json:"changePasswordOnLogin"`
	IsSuperuser             bool   `json:"isSuperuser"`
}

type SignInResponseBody struct {
	BaseResponse
	Account     AccountLite `json:"account"`
	AccessToken string      `json:"accessToken"`
}

type SignInReponse struct {
	Body SignInResponseBody
}

type SignOutRequest struct {
	Body BaseRequest
}

type SignOutResponse struct {
	Body BaseResponse
}
