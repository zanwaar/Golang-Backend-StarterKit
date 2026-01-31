package service

import (
	"bytes"
	"encoding/base64"
	"errors"
	"golang-backend/dto"
	"golang-backend/entity"
	"golang-backend/repository"
	"golang-backend/utils"
	"image/png"
	"time"

	"github.com/pquerna/otp/totp"
)

type UserService interface {
	Login(req dto.UserLoginRequest) (string, error)
	Register(req dto.UserRegisterRequest) (*dto.UserResponse, error)
	VerifyEmail(req dto.VerifyEmailRequest) error
	ForgotPassword(email string) error
	ResetPassword(req dto.ResetPasswordRequest) error
	ResendVerificationCode(email string) error
	ResendResetPasswordCode(email string) error
	GetUsers(filters map[string]interface{}, page, perPage int) (*utils.PaginationResult, error)
	Setup2FA(userID string) (*dto.Setup2FAResponse, error)
	Verify2FA(userID string, code string) error
	GetMe(userID string) (*dto.UserResponse, error)
}

type userService struct {
	repo repository.UserRepository
}

func (s *userService) GetMe(userID string) (*dto.UserResponse, error) {
	user, err := s.repo.FindByID(userID)
	if err != nil {
		return nil, err
	}

	return &dto.UserResponse{
		ID:             user.ID,
		Name:           user.Name,
		Email:          user.Email,
		IsTwoFAEnabled: user.IsTwoFAEnabled,
	}, nil
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) GetUsers(filters map[string]interface{}, page, perPage int) (*utils.PaginationResult, error) {
	return s.repo.Paginate(filters, page, perPage)
}

func (s *userService) Login(req dto.UserLoginRequest) (string, error) {
	user, err := s.repo.FindByEmail(req.Email)
	if err != nil {
		return "", errors.New("invalid email or password")
	}

	if !user.IsVerified {
		return "", errors.New("email not verified")
	}

	err = utils.CheckPassword(req.Password, user.Password)
	if err != nil {
		return "", errors.New("invalid email or password")
	}

	if user.IsTwoFAEnabled {
		if req.TwoFACode == "" {
			return "", errors.New("2FA code required")
		}
		if !totp.Validate(req.TwoFACode, user.TwoFASecret) {
			return "", errors.New("invalid 2FA code")
		}
	}

	token, err := utils.GenerateToken(user.ID)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *userService) Register(req dto.UserRegisterRequest) (*dto.UserResponse, error) {
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	code := utils.GenerateRandomCode(6)

	user := &entity.User{
		Name:             req.Name,
		Email:            req.Email,
		Password:         hashedPassword,
		VerificationCode: code,
		IsVerified:       false,
	}

	if err := s.repo.Create(user); err != nil {
		return nil, err
	}

	// Send email asynchronously
	go func() {
		_ = utils.SendVerificationEmail(user.Email, code)
	}()

	return &dto.UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}, nil
}

func (s *userService) VerifyEmail(req dto.VerifyEmailRequest) error {
	user, err := s.repo.FindByEmail(req.Email)
	if err != nil {
		return errors.New("user not found")
	}

	if user.IsVerified {
		return errors.New("email already verified")
	}

	if user.VerificationCode != req.Code {
		return errors.New("invalid verification code")
	}

	user.IsVerified = true
	user.VerificationCode = ""

	return s.repo.Update(user)
}

func (s *userService) ForgotPassword(email string) error {
	user, err := s.repo.FindByEmail(email)
	if err != nil {
		return errors.New("user not found")
	}

	code := utils.GenerateRandomCode(6)
	user.ResetToken = code
	user.ResetTokenExpiry = time.Now().Add(15 * time.Minute)

	if err := s.repo.Update(user); err != nil {
		return err
	}

	// Send email asynchronously
	go func() {
		_ = utils.SendResetPasswordEmail(user.Email, code)
	}()

	return nil
}

func (s *userService) ResetPassword(req dto.ResetPasswordRequest) error {
	user, err := s.repo.FindByEmail(req.Email)
	if err != nil {
		return errors.New("user not found")
	}

	if user.ResetToken != req.Code {
		return errors.New("invalid or expired reset code")
	}

	if time.Now().After(user.ResetTokenExpiry) {
		return errors.New("reset code expired")
	}

	hashedPassword, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		return err
	}

	user.Password = hashedPassword
	user.ResetToken = ""
	user.ResetTokenExpiry = time.Time{}

	return s.repo.Update(user)
}

func (s *userService) ResendVerificationCode(email string) error {
	user, err := s.repo.FindByEmail(email)
	if err != nil {
		return errors.New("user not found")
	}

	if user.IsVerified {
		return errors.New("email already verified")
	}

	code := utils.GenerateRandomCode(6)
	user.VerificationCode = code

	if err := s.repo.Update(user); err != nil {
		return err
	}

	// Send email asynchronously
	go func() {
		_ = utils.SendVerificationEmail(user.Email, code)
	}()

	return nil
}

func (s *userService) ResendResetPasswordCode(email string) error {
	// Functionally same as ForgotPassword
	return s.ForgotPassword(email)
}

func (s *userService) Setup2FA(userID string) (*dto.Setup2FAResponse, error) {
	user, err := s.repo.FindByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "GolangBackend",
		AccountName: user.Email,
	})
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	img, err := key.Image(200, 200)
	if err != nil {
		return nil, err
	}
	png.Encode(&buf, img)
	qrCodeURL := "data:image/png;base64," + base64.StdEncoding.EncodeToString(buf.Bytes())

	user.TwoFASecret = key.Secret()
	if err := s.repo.Update(user); err != nil {
		return nil, err
	}

	return &dto.Setup2FAResponse{
		Secret:    key.Secret(),
		QRCodeURL: qrCodeURL,
	}, nil
}

func (s *userService) Verify2FA(userID string, code string) error {
	user, err := s.repo.FindByID(userID)
	if err != nil {
		return errors.New("user not found")
	}

	if user.TwoFASecret == "" {
		return errors.New("2FA not setup")
	}

	if !totp.Validate(code, user.TwoFASecret) {
		return errors.New("invalid code")
	}

	user.IsTwoFAEnabled = true
	return s.repo.Update(user)
}
