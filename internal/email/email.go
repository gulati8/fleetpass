package email

import (
	"log"
)

// Service defines the interface for sending emails
type Service interface {
	SendVerificationEmail(to, token string) error
	SendPasswordResetEmail(to, token string) error
	SendWelcomeEmail(to, firstName string) error
}

// MockService is a mock email service that logs to console
type MockService struct{}

// NewMockService creates a new mock email service
func NewMockService() *MockService {
	return &MockService{}
}

// SendVerificationEmail logs a verification email to console
func (s *MockService) SendVerificationEmail(to, token string) error {
	log.Println("========================================")
	log.Println("📧 EMAIL: Verification Email")
	log.Println("========================================")
	log.Printf("To: %s\n", to)
	log.Println("Subject: Verify Your Email Address")
	log.Println("----------------------------------------")
	log.Println("Welcome to FleetPass!")
	log.Println()
	log.Println("Please click the link below to verify your email address:")
	log.Printf("http://localhost:3000/verify-email?token=%s\n", token)
	log.Println()
	log.Println("This link will expire in 24 hours.")
	log.Println()
	log.Println("If you didn't create an account, please ignore this email.")
	log.Println("========================================")
	return nil
}

// SendPasswordResetEmail logs a password reset email to console
func (s *MockService) SendPasswordResetEmail(to, token string) error {
	log.Println("========================================")
	log.Println("📧 EMAIL: Password Reset")
	log.Println("========================================")
	log.Printf("To: %s\n", to)
	log.Println("Subject: Reset Your Password")
	log.Println("----------------------------------------")
	log.Println("You requested to reset your password.")
	log.Println()
	log.Println("Click the link below to reset your password:")
	log.Printf("http://localhost:3000/reset-password?token=%s\n", token)
	log.Println()
	log.Println("This link will expire in 1 hour.")
	log.Println()
	log.Println("If you didn't request this, please ignore this email.")
	log.Println("Your password will not be changed.")
	log.Println("========================================")
	return nil
}

// SendWelcomeEmail logs a welcome email to console
func (s *MockService) SendWelcomeEmail(to, firstName string) error {
	log.Println("========================================")
	log.Println("📧 EMAIL: Welcome")
	log.Println("========================================")
	log.Printf("To: %s\n", to)
	log.Println("Subject: Welcome to FleetPass!")
	log.Println("----------------------------------------")
	log.Printf("Hi %s,\n", firstName)
	log.Println()
	log.Println("Welcome to FleetPass! Your email has been verified.")
	log.Println()
	log.Println("You can now log in and start managing your fleet.")
	log.Println()
	log.Println("If you have any questions, feel free to reach out to our support team.")
	log.Println()
	log.Println("Best regards,")
	log.Println("The FleetPass Team")
	log.Println("========================================")
	return nil
}

// TODO: Implement real email service (SendGrid, AWS SES, etc.)
// Example:
//
// type SendGridService struct {
//     apiKey string
//     client *sendgrid.Client
// }
//
// func NewSendGridService(apiKey string) *SendGridService {
//     return &SendGridService{
//         apiKey: apiKey,
//         client: sendgrid.NewSendClient(apiKey),
//     }
// }

// Helper function to get email service (easily swap mock for real)
func GetEmailService() Service {
	// For now, return mock service
	// In production, check environment variable and return appropriate service
	// if os.Getenv("EMAIL_SERVICE") == "sendgrid" {
	//     return NewSendGridService(os.Getenv("SENDGRID_API_KEY"))
	// }
	return NewMockService()
}
