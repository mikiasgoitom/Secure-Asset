package usecase

import (
	"context"
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"strings"

	"github.com/casbin/casbin/v2"
	"github.com/google/uuid"
	"github.com/mikiasgoitom/Secure-Asset/internal/contract"
	"github.com/mikiasgoitom/Secure-Asset/internal/domain/entity"
	"github.com/mikiasgoitom/Secure-Asset/internal/domain/valueobject"
	"golang.org/x/crypto/bcrypt"
)

type UserUsecase struct {
	userRepository contract.IUserRepository
	logger         contract.ILogger
	jwtservice     contract.IJWTService
	casbinEnforcer  *casbin.Enforcer
	emailService   contract.IEmailService
}

func NewUserUsecase(userRepo contract.IUserRepository, logger contract.ILogger, jwtservice contract.IJWTService, casbinEnforcer *casbin.Enforcer, emailService contract.IEmailService) contract.IUserUsecase {
	return &UserUsecase{
		userRepository: userRepo,
		logger:         logger,
		jwtservice:     jwtservice,
		casbinEnforcer: casbinEnforcer,
		emailService:   emailService,
	}
}

func (u *UserUsecase) Register(ctx context.Context, username, email, password, role string) (*entity.User, error) {
	if existingUser, _ := u.userRepository.FindByEmail(ctx, email); existingUser != nil {
		u.logger.Error("Email already in use", valueobject.Field{Key: "email", Value: email})
		return nil, nil
	}
	if existingUser, _ := u.userRepository.FindByUsername(ctx, username); existingUser != nil {
		u.logger.Error("Username already in use", valueobject.Field{Key: "username", Value: username})
		return nil, nil
	}
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), 5)
	if err != nil {
		u.logger.Error("Failed to generate password hash", valueobject.Field{Key: "error", Value: err})
		return nil, err
	}

	otpCode := generateNumericOTP()
	if u.emailService != nil {
		if err := u.emailService.SendOTP(email, otpCode); err != nil {
			u.logger.Error("Failed to send OTP email", valueobject.Field{Key: "error", Value: err}, valueobject.Field{Key: "email", Value: email})
			return nil, err
		}
	}
	newUser := &entity.User{
		ID:              uuid.New().String(),
		Username:        username,
		Email:           email,
		Password:        string(hashPassword),
		Role:            role,
		Department:      "",
		ClearanceLevel:  valueobject.GetClassificationLevel("Restricted"),
		IsMFAEnabled:    false,
		IsAccountLocked: false,
		OTPSecret:       otpCode,
	}
	createdUser, err := u.userRepository.Create(ctx, newUser)
	if err != nil {
		u.logger.Error("Failed to create user", valueobject.Field{Key: "error", Value: err})
		return nil, err
	}
	// assigning role to user in casbin
	_, err = u.casbinEnforcer.AddRoleForUser(createdUser.ID, createdUser.Role)
	if err != nil {
		u.logger.Error("Failed to assign role to user in Casbin", valueobject.Field{Key: "error", Value: err})
		return nil, err
	}

	u.logger.Info("User registered successfully", valueobject.Field{Key: "userID", Value: createdUser.ID})
	return createdUser, nil
}

func (u *UserUsecase) Login(ctx context.Context, identifier, password string) (string, error) {
	var user *entity.User
	var err error
	if strings.Contains(identifier, "@") {
		user, err = u.userRepository.FindByEmail(ctx, identifier)
	} else {
		user, err = u.userRepository.FindByUsername(ctx, identifier)
	}
	if err != nil {
		u.logger.Error("Failed to find user by identifier", valueobject.Field{Key: "error", Value: err})
		return "", err
	}

	if user == nil {
		u.logger.Info("User not found by email, trying username", valueobject.Field{Key: "identifier", Value: identifier})
		user, err = u.userRepository.FindByUsername(ctx, identifier)
		if err != nil {
			u.logger.Error("Failed to find user by username", valueobject.Field{Key: "error", Value: err})
			return "", err
		}
		if user == nil {
			u.logger.Info("User not found", valueobject.Field{Key: "identifier", Value: identifier})
			return "", fmt.Errorf("invalid credentials")
		}
	}


	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		u.logger.Error("Invalid password", valueobject.Field{Key: "error", Value: err})
		return "", fmt.Errorf("invalid credentials")
	}

	if user.IsAccountLocked {
		u.logger.Info("Account is locked", valueobject.Field{Key: "userID", Value: user.ID})
		return "", fmt.Errorf("account is locked")
	}

	token, err := u.jwtservice.GenerateToken(user)
	if err != nil {
		u.logger.Error("Failed to generate JWT token", valueobject.Field{Key: "userID", Value: user.ID})
		return "", fmt.Errorf("failed to generate token")
	}

	u.logger.Info("User logged in successfully", valueobject.Field{Key: "userID", Value: user.ID})
	return token, nil

}

// generateNumericOTP creates a 6-digit numeric OTP using crypto/rand.
func generateNumericOTP() string {
	var b [8]byte
	if _, err := rand.Read(b[:]); err != nil {
		return "000000"
	}
	code := binary.BigEndian.Uint64(b[:]) % 1000000
	return fmt.Sprintf("%06d", code)
}
