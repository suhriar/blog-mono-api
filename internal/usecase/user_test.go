package usecase

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/suhriar/blog-mono-api/config"
	"github.com/suhriar/blog-mono-api/internal/repository/mysql/mocks"
	"github.com/suhriar/blog-mono-api/model"
	"golang.org/x/crypto/bcrypt"
)

func TestSignUp(t *testing.T) {
	ctx := context.Background()

	req := model.SignUpRequest{
		Email:    "test@example.com",
		Username: "testuser",
		Password: "securepassword",
	}

	t.Run("Success SignUp", func(t *testing.T) {
		mockRepo := new(mocks.MockUserRepository)
		usecase := &userUsecase{userRepository: mockRepo}
		mockRepo.On("GetUser", ctx, req.Email, req.Username, int64(0)).Return(model.User{}, nil)
		mockRepo.On("CreateUser", ctx, mock.AnythingOfType("model.User")).Return(int64(1), nil)

		err := usecase.SignUp(ctx, req)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Fail SignUp - Username or Email Already Exists", func(t *testing.T) {
		mockRepo := new(mocks.MockUserRepository)
		usecase := &userUsecase{userRepository: mockRepo}
		mockRepo.On("GetUser", ctx, req.Email, req.Username, int64(0)).Return(model.User{ID: 1}, nil)

		err := usecase.SignUp(ctx, req)

		assert.Error(t, err)
		assert.Equal(t, "username or email already exist", err.Error())
		mockRepo.AssertExpectations(t)
	})
}

func TestLogin(t *testing.T) {
	ctx := context.Background()

	config.LoadConfig()

	req := model.LoginRequest{
		Email:    "test@example.com",
		Password: "securepassword",
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	mockUser := model.User{
		ID:       1,
		Email:    req.Email,
		Password: string(hashedPassword),
	}

	t.Run("Success Login", func(t *testing.T) {
		mockRepo := new(mocks.MockUserRepository)
		usecase := &userUsecase{userRepository: mockRepo}
		mockRepo.On("GetUser", ctx, req.Email, "", int64(0)).Return(mockUser, nil)
		mockRepo.On("GetRefreshToken", ctx, mockUser.ID, mock.Anything).Return(model.RefreshToken{}, nil)
		mockRepo.On("InsertRefreshToken", ctx, mock.AnythingOfType("model.RefreshToken")).Return(int64(1), nil)

		jwtToken, refreshToken, err := usecase.Login(ctx, req)

		assert.NoError(t, err)
		assert.NotEmpty(t, jwtToken)
		assert.NotEmpty(t, refreshToken)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Fail Login - Email Not Exist", func(t *testing.T) {
		mockRepo := new(mocks.MockUserRepository)
		usecase := &userUsecase{userRepository: mockRepo}
		mockRepo.On("GetUser", ctx, req.Email, "", int64(0)).Return(model.User{}, nil)

		jwtToken, refreshToken, err := usecase.Login(ctx, req)

		assert.Error(t, err)
		assert.Equal(t, "email not exist", err.Error())
		assert.Empty(t, jwtToken)
		assert.Empty(t, refreshToken)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Fail Login - Password Incorrect", func(t *testing.T) {
		mockRepo := new(mocks.MockUserRepository)
		usecase := &userUsecase{userRepository: mockRepo}
		mockUser.Password = "wrongpassword"
		mockRepo.On("GetUser", ctx, req.Email, "", int64(0)).Return(mockUser, nil)

		jwtToken, refreshToken, err := usecase.Login(ctx, req)

		assert.Error(t, err)
		assert.Equal(t, "email or password is invalid", err.Error())
		assert.Empty(t, jwtToken)
		assert.Empty(t, refreshToken)
		mockRepo.AssertExpectations(t)
	})
}

func TestValidateRefreshToken(t *testing.T) {
	ctx := context.Background()
	userID := int64(1)

	req := model.RefreshTokenRequest{
		Token: "valid-refresh-token",
	}

	mockRefreshToken := model.RefreshToken{
		UserID:       userID,
		RefreshToken: req.Token,
		ExpiredAt:    time.Now().Add(24 * time.Hour),
	}

	mockUser := model.User{
		ID:       userID,
		Username: "testuser",
		Email:    "test@example.com",
	}

	t.Run("Success Validate Refresh Token", func(t *testing.T) {
		mockRepo := new(mocks.MockUserRepository)
		usecase := &userUsecase{userRepository: mockRepo}
		mockRepo.On("GetRefreshToken", ctx, userID, mock.Anything).Return(mockRefreshToken, nil)
		mockRepo.On("GetUser", ctx, "", "", userID).Return(mockUser, nil)

		jwtToken, err := usecase.ValidateRefreshToken(ctx, userID, req)

		assert.NoError(t, err)
		assert.NotEmpty(t, jwtToken)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Fail Validate - Refresh Token Expired", func(t *testing.T) {
		mockRepo := new(mocks.MockUserRepository)
		usecase := &userUsecase{userRepository: mockRepo}
		mockExpiredToken := mockRefreshToken
		mockExpiredToken.RefreshToken = ""
		mockRepo.On("GetRefreshToken", ctx, userID, mock.Anything).Return(mockExpiredToken, nil)

		jwtToken, err := usecase.ValidateRefreshToken(ctx, userID, req)

		assert.Error(t, err)
		assert.Equal(t, "refresh token has expired", err.Error())
		assert.Empty(t, jwtToken)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Fail Validate - Refresh Token Mismatch", func(t *testing.T) {
		mockRepo := new(mocks.MockUserRepository)
		usecase := &userUsecase{userRepository: mockRepo}
		mockRefreshToken.RefreshToken = "different-token"
		mockRepo.On("GetRefreshToken", ctx, userID, mock.Anything).Return(mockRefreshToken, nil)

		jwtToken, err := usecase.ValidateRefreshToken(ctx, userID, req)

		assert.Error(t, err)
		assert.Equal(t, "refresh token is invalid", err.Error())
		assert.Empty(t, jwtToken)
		mockRepo.AssertExpectations(t)
	})
}
