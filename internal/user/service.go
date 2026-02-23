package user

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log/slog"
	"time"

	dbmodel "mindoh-service/internal/db"
	"mindoh-service/internal/mailer"

	"gorm.io/gorm"
)

// UserService handles business logic for users
type UserService struct {
	DB     *gorm.DB
	Repo   *UserRepository
	Mailer mailer.IMailer
	AppURL string // Frontend base URL for links in emails
}

// NewUserService creates a new user service
func NewUserService(db *gorm.DB, m mailer.IMailer, appURL string) *UserService {
	return &UserService{
		DB:     db,
		Repo:   NewUserRepository(db),
		Mailer: m,
		AppURL: appURL,
	}
}

// GetUserByID retrieves a user by their ID
func (s *UserService) GetUserByID(id uint) (*dbmodel.User, error) {
	return s.Repo.GetByID(id)
}

// CreateUser creates a new user and sends a verification email
func (s *UserService) CreateUser(user *dbmodel.User) error {
	token, err := generateToken()
	if err != nil {
		return err
	}
	user.IsEmailVerified = false
	user.EmailVerifyToken = token
	user.EmailVerifyExpiry = time.Now().Add(24 * time.Hour)

	if err := s.Repo.Create(user); err != nil {
		return err
	}

	// Send verification email (non-blocking — don't fail registration if email fails)
	go s.sendVerifyEmail(user.Email, token)
	return nil
}

// UpdateUser updates an existing user
func (s *UserService) UpdateUser(user *dbmodel.User) error {
	return s.Repo.Update(user)
}

// DeleteUser deletes a user by their ID
func (s *UserService) DeleteUser(id uint) error {
	return s.Repo.Delete(id)
}

// ValidateCredentials validates user login credentials.
func (s *UserService) ValidateCredentials(username, password string) (*dbmodel.User, error) {
	user, err := s.Repo.GetByUsername(username)
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}
	if !CheckPasswordHash(password, user.PasswordHash) {
		return nil, fmt.Errorf("invalid credentials")
	}
	return user, nil
}

// VerifyEmail confirms an email address using the given token.
func (s *UserService) VerifyEmail(token string) error {
	user, err := s.Repo.GetByEmailVerifyToken(token)
	if err != nil {
		return fmt.Errorf("invalid or expired token")
	}
	if time.Now().After(user.EmailVerifyExpiry) {
		return fmt.Errorf("token expired")
	}
	return s.Repo.UpdateFields(user.ID, map[string]interface{}{
		"is_email_verified":  true,
		"email_verify_token": "",
	})
}

// ResendVerificationEmail generates a fresh token and re-sends the verification email.
func (s *UserService) ResendVerificationEmail(email string) error {
	user, err := s.Repo.GetByEmail(email)
	if err != nil {
		// Return nil to not leak whether email exists
		return nil
	}
	if user.IsEmailVerified {
		return nil
	}
	token, err := generateToken()
	if err != nil {
		return err
	}
	if err := s.Repo.UpdateFields(user.ID, map[string]interface{}{
		"email_verify_token":  token,
		"email_verify_expiry": time.Now().Add(24 * time.Hour),
	}); err != nil {
		return err
	}
	go s.sendVerifyEmail(user.Email, token)
	return nil
}

// ForgotPassword generates a reset token and sends a password reset email.
func (s *UserService) ForgotPassword(email string) error {
	user, err := s.Repo.GetByEmail(email)
	if err != nil {
		// Return nil to not leak whether email exists
		return nil
	}
	// Only send reset email to verified accounts
	if !user.IsEmailVerified {
		return nil
	}
	token, err := generateToken()
	if err != nil {
		return err
	}
	if err := s.Repo.UpdateFields(user.ID, map[string]interface{}{
		"password_reset_token":  token,
		"password_reset_expiry": time.Now().Add(1 * time.Hour),
	}); err != nil {
		return err
	}
	go s.sendPasswordResetEmail(user.Email, token)
	return nil
}

// ChangePassword verifies the current password then updates to the new one.
func (s *UserService) ChangePassword(userID uint, currentPassword, newPassword string) error {
	user, err := s.Repo.GetByID(userID)
	if err != nil {
		return fmt.Errorf("user not found")
	}
	if !CheckPasswordHash(currentPassword, user.PasswordHash) {
		return fmt.Errorf("current password is incorrect")
	}
	hash, err := HashPassword(newPassword)
	if err != nil {
		return err
	}
	return s.Repo.UpdateFields(userID, map[string]interface{}{
		"password_hash": hash,
	})
}

// ResetPassword validates the reset token and sets a new password.
func (s *UserService) ResetPassword(token, newPassword string) error {
	user, err := s.Repo.GetByPasswordResetToken(token)
	if err != nil {
		return fmt.Errorf("invalid or expired token")
	}
	if time.Now().After(user.PasswordResetExpiry) {
		return fmt.Errorf("token expired")
	}
	hash, err := HashPassword(newPassword)
	if err != nil {
		return err
	}
	return s.Repo.UpdateFields(user.ID, map[string]interface{}{
		"password_hash":        hash,
		"password_reset_token": "",
	})
}

// --- helpers ---

func generateToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func (s *UserService) sendVerifyEmail(to, token string) {
	link := fmt.Sprintf("%s/verify-email?token=%s", s.AppURL, token)
	html := fmt.Sprintf(`
<p>Hi there,</p>
<p>Thanks for signing up for <strong>Mindoh</strong>! Please verify your email by clicking the button below.</p>
<p><a href="%s" style="background:#28C76F;color:#fff;padding:12px 24px;border-radius:6px;text-decoration:none;font-weight:600;">Verify Email</a></p>
<p>Or paste this link into your browser:<br><a href="%s">%s</a></p>
<p>This link expires in 24 hours.</p>
`, link, link, link)
	if err := s.Mailer.Send(to, "Verify your Mindoh email", html); err != nil {
		slog.Error("failed to send verification email", "to", to, "error", err)
	}
}

func (s *UserService) sendPasswordResetEmail(to, token string) {
	link := fmt.Sprintf("%s/reset-password?token=%s", s.AppURL, token)
	html := fmt.Sprintf(`
<p>Hi,</p>
<p>We received a request to reset your <strong>Mindoh</strong> password. Click the button below to set a new password.</p>
<p><a href="%s" style="background:#28C76F;color:#fff;padding:12px 24px;border-radius:6px;text-decoration:none;font-weight:600;">Reset Password</a></p>
<p>Or paste this link into your browser:<br><a href="%s">%s</a></p>
<p>This link expires in 1 hour. If you did not request a password reset, you can safely ignore this email.</p>
`, link, link, link)
	if err := s.Mailer.Send(to, "Reset your Mindoh password", html); err != nil {
		slog.Error("failed to send password reset email", "to", to, "error", err)
	}
}
