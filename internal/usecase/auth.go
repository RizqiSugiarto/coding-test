package usecase

import (
	"context"

	"github.com/RizqiSugiarto/coding-test/internal/dto"
	"github.com/RizqiSugiarto/coding-test/internal/repository"
	"github.com/RizqiSugiarto/coding-test/pkg/apperror"
	"github.com/RizqiSugiarto/coding-test/pkg/jwt"
	"golang.org/x/crypto/bcrypt"
)

type AuthUseCase struct {
	userRepo   repository.UserRepo
	jwtManager jwt.Manager
}

func NewAuthUseCase(up repository.UserRepo, jwtMng jwt.Manager) *AuthUseCase {
	return &AuthUseCase{
		userRepo:   up,
		jwtManager: jwtMng,
	}
}

func (au *AuthUseCase) Login(ctx context.Context, req dto.LoginRequestDTO) (*dto.LoginResponseDTO, error) {
	user, err := au.userRepo.GetByUsername(ctx, req.UserName)
	if err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, apperror.ErrInvalidCredentials
	}

	accToken, err := au.jwtManager.GenerateAccessToken(user.ID)
	if err != nil {
		return nil, err
	}

	refrToken, err := au.jwtManager.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, err
	}

	resp := &dto.LoginResponseDTO{
		AccessToken:  accToken,
		RefreshToken: refrToken,
	}

	return resp, nil
}
