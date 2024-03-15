package service

import (
	"github.com/Dubjay18/gobank_auth/domain"
	"github.com/Dubjay18/gobank_auth/dto"
	"github.com/Dubjay18/gobank_auth/errs"
)

type AuthService interface {
	Login(request dto.LoginRequest) (*dto.LoginResponse, *errs.AppError)
}

type DefaultAuthService struct {
	repo            domain.AuthRepository
	rolePermissions domain.RolePermissions
}

func (s DefaultAuthService) Login(req dto.LoginRequest) (*dto.LoginResponse, *errs.AppError) {
	var appErr *errs.AppError
	var login *domain.SLogin

	if login, appErr = s.repo.FindById(req.Username, req.Password); appErr != nil {
		return nil, appErr
	}

	claims := login.ClaimsForAccessToken()
	authToken := domain.NewAuthToken(claims)

	var accessToken, refreshToken string
	if accessToken, appErr = authToken.NewAccessToken(); appErr != nil {
		return nil, appErr
	}

	if refreshToken, appErr = s.repo.GenerateAndSaveRefreshTokenToStore(authToken); appErr != nil {
		return nil, appErr
	}

	return &dto.LoginResponse{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}
func NewLoginService(repo domain.AuthRepository, permissions domain.RolePermissions) DefaultAuthService {
	return DefaultAuthService{repo, permissions}
}
