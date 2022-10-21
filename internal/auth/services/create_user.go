package services

import (
	"bitbucket.org/sercide/data-ingestion/internal/auth"
	"context"
)

type CreateUserDto struct {
	ID          string
	Email       string
	Name        string
	IsAdmin     bool
	Password    string
	Permissions []struct {
		DistributorID string
		RoleID        string
	}
}

type CreateUserService struct {
	authRepository auth.Repository
	oAuthClient    auth.OAuthClient
}

func NewCreateUserService(authRepository auth.Repository, oAuthClient auth.OAuthClient) *CreateUserService {
	return &CreateUserService{authRepository: authRepository, oAuthClient: oAuthClient}
}

func (s CreateUserService) Handler(ctx context.Context, dto CreateUserDto) error {
	user := auth.User{
		ID:           dto.ID,
		Name:         dto.Name,
		Email:        dto.Email,
		AuthID:       "",
		IsAdmin:      dto.IsAdmin,
		Distributors: make([]auth.Distributor, 0, cap(dto.Permissions)),
	}

	authId, err := s.oAuthClient.CreateUser(ctx, user, dto.Password)
	if err != nil {
		return err
	}

	user.AuthID = authId

	for _, p := range dto.Permissions {
		distributor, err := s.authRepository.GetDistributor(ctx, p.DistributorID)
		if err != nil {
			continue
		}
		role, err := s.authRepository.GetRole(ctx, p.RoleID)
		if err != nil {
			continue
		}
		distributor.Role = role
		user.Distributors = append(user.Distributors, distributor)
	}

	_, err = s.authRepository.SaveUser(ctx, user)

	if err != nil {
		_ = s.oAuthClient.DeleteUser(ctx, user.ID)
	}

	return err
}
