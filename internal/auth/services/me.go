package services

import (
	"bitbucket.org/sercide/data-ingestion/internal/auth"
	"context"
)

type MeDto struct {
	Token string
}

type MeService struct {
	authRepository auth.Repository
	oAuthClient    auth.OAuthClient
}

func NewMeService(authRepository auth.Repository, oAuthClient auth.OAuthClient) *MeService {
	return &MeService{authRepository: authRepository, oAuthClient: oAuthClient}
}

func (s MeService) Handler(ctx context.Context, dto MeDto) (auth.User, error) {
	id, err := s.oAuthClient.VerifyToken(ctx, dto.Token)

	if err != nil {
		return auth.User{}, err
	}

	user, err := s.authRepository.Me(ctx, id)

	return user, err
}
