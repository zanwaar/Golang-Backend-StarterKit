package service

import (
	"errors"
	"golang-backend/dto"
	"golang-backend/entity"
	"golang-backend/repository"
	"golang-backend/utils"
	"time"
)

type UserService interface {
	Login(req dto.UserLoginRequest) (string, error)
	Register(req dto.UserRegisterRequest) (*dto.UserResponse, error)
	VerifyEmail(req dto.VerifyEmailRequest) error
	ForgotPassword(email string) error
	ResetPassword(req dto.ResetPasswordRequest) error
	GetUsers(filters map[string]interface{}, page, perPage int) (*utils.PaginationResult, error)
}

type userService struct {
	repo repository.UserRepository
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

	if err := utils.CheckPassword(req.Password, user.Password); err != nil {
		return "", errors.New("invalid email or password")
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
