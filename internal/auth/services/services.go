package services

import "bitbucket.org/sercide/data-ingestion/internal/auth"

type Services struct {
	MeService         *MeService
	CreateUserService *CreateUserService
}

func NewService(authRepository auth.Repository, oAuthClient auth.OAuthClient) *Services {
	return &Services{
		MeService:         NewMeService(authRepository, oAuthClient),
		CreateUserService: NewCreateUserService(authRepository, oAuthClient),
	}
}
