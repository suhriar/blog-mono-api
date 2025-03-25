package usecase

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/suhriar/blog-mono-api/model"
	"github.com/suhriar/blog-mono-api/pkg/utils"
	"golang.org/x/crypto/bcrypt"
)

func (u *userUsecase) SignUp(ctx context.Context, req model.SignUpRequest) (err error) {
	user, err := u.userRepository.GetUser(ctx, req.Email, req.Username, 0)
	if err != nil {
		return err
	}

	if user.ID != 0 {
		return errors.New("username or email already exist")
	}

	pass, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	now := time.Now()
	user = model.User{
		Email:     req.Email,
		Username:  req.Username,
		Password:  string(pass),
		CreatedAt: now,
		UpdatedAt: now,
		CreatedBy: req.Email,
		UpdatedBy: req.Email,
	}

	_, err = u.userRepository.CreateUser(ctx, user)
	if err != nil {
		return err
	}

	return
}

func (u *userUsecase) Login(ctx context.Context, req model.LoginRequest) (jwtTokenstring, refreshToken string, err error) {
	user, err := u.userRepository.GetUser(ctx, req.Email, "", 0)
	if err != nil {
		return "", "", err
	}

	if user.ID == 0 {
		return "", "", errors.New("email not exist")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return "", "", errors.New("email or password is invalid")
	}

	jwtToken, err := utils.GenerateJWT(user.ID, user.Username, user.Email)
	if err != nil {
		return "", "", err
	}

	existingRefreshToken, err := u.userRepository.GetRefreshToken(ctx, user.ID, time.Now())
	if err != nil {
		return "", "", err
	}

	if existingRefreshToken.RefreshToken != "" {
		return jwtToken, existingRefreshToken.RefreshToken, nil
	}

	refreshToken = utils.GenerateRefreshToken()
	if refreshToken == "" {
		return jwtToken, "", errors.New("failed to generate refresh token")
	}

	_, err = u.userRepository.InsertRefreshToken(ctx, model.RefreshToken{
		UserID:       user.ID,
		RefreshToken: refreshToken,
		ExpiredAt:    time.Now().Add(10 * 24 * time.Hour),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		CreatedBy:    strconv.FormatInt(user.ID, 10),
		UpdatedBy:    strconv.FormatInt(user.ID, 10),
	})
	if err != nil {
		return jwtToken, refreshToken, err
	}

	return jwtToken, refreshToken, nil
}

func (u *userUsecase) ValidateRefreshToken(ctx context.Context, userID int64, request model.RefreshTokenRequest) (jwtToken string, err error) {
	existingRefreshToken, err := u.userRepository.GetRefreshToken(ctx, userID, time.Now())
	if err != nil {
		return "", err
	}

	if existingRefreshToken.RefreshToken == "" {
		return "", errors.New("refresh token has expired")
	}

	// means the token in database is not matched with request token, throw error invalid refresh token
	if existingRefreshToken.RefreshToken != request.Token {
		return "", errors.New("refresh token is invalid")
	}

	user, err := u.userRepository.GetUser(ctx, "", "", userID)
	if err != nil {
		log.Error().Err(err).Msg("failed to get user")
		return "", err
	}
	if user.ID == 0 {
		return "", errors.New("user not exist")
	}

	jwtToken, err = utils.GenerateJWT(user.ID, user.Username, user.Email)
	if err != nil {
		return "", err
	}
	return jwtToken, nil
}
