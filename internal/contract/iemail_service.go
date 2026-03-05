package contract

// IEmailService defines email-related capabilities (e.g., sending OTP codes).
type IEmailService interface {
	SendOTP(toEmail string, otpCode string) error
}
