package models

// REQUESTS

type LoginRequest struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type OAuthGoogleLoginRequest struct {
	Token string `json:"token" validate:"required"`
}

type OAuthFacebookLoginRequest struct {
	Token string `json:"token" validate:"required"`
}

type CheckEmailRequest struct {
	Email string `json:"email" validate:"required"`
}

type RegisterRequest struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email" validate:"required"`
}

// RESPONSES

type LoginResponse struct {
	Success bool `json:"success"`
}

type OAuthGoogleLoginResponse struct {
	Success bool `json:"success"`
}

type OAuthFacebookLoginResponse struct {
	Success bool `json:"success"`
}

type CheckEmailResponse struct {
	AccountExist bool `json:"account_exist"`
}

type RegisterResponse struct {
	Success bool `json:"success"`
}

type ForgotPasswordResponse struct {
	Success bool `json:"success"`
}
