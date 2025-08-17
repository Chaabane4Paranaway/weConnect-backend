package dto

type GlobalSuccess struct {
	CodeStatus int         `json:"codeStatus" example:"200"`
	Message    string      `json:"message" example:"Succès"`
	Data       interface{} `json:"data"`
}

type GlobalError struct {
	CodeStatus       int    `json:"codeStatus" example:"400"`
	Message          string `json:"message" example:"Erreur"`
	TechnicalMessage string `json:"technical" example:"Internal Server Error"`
}

// Auth

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email" example:"user@example.com"`
	Password string `json:"password" binding:"required,min=8,max=100" example:"StrongP@ssword123"`
}

type RegisterResponse struct {
	Email    string `json:"email" example:"user@example.com"`
	Verified bool   `json:"verified" example:"false"`
	Token    string `json:"token" example:"jwt.token.here"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email" example:"user@example.com"`
	Password string `json:"password" binding:"required,min=8,max=100" example:"StrongP@ssword123"`
}

type LoginResponse struct {
	Token    string `json:"token" example:"jwt.token.here"`
	Email    string `json:"email" example:"user@example.com"`
	Verified bool   `json:"verified" example:"true"`
}

// OTP

type VerifyOTPRequest struct {
	Email string `json:"email" binding:"required,email" example:"user@example.com"`
	OTP   string `json:"otp" binding:"required,len=6" example:"123456"`
}

// Me

type MeResponse struct {
	ID        string `json:"id" example:"9f1fc230-b312-4c3e-a9ff-5c8e12345678"`
	Email     string `json:"email" example:"user@example.com"`
	Verified  bool   `json:"verified" example:"true"`
	CreatedAt string `json:"created_at" example:"2025-07-12T14:30:00Z"`
}
