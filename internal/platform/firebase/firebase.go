package firebase

import (
	"bitbucket.org/sercide/data-ingestion/internal/auth"
	"context"
	"errors"
	firebase "firebase.google.com/go"
	auth_firebase "firebase.google.com/go/auth"
	"google.golang.org/api/option"
	"os"
)

type OAuthClient struct {
	AuthClient *auth_firebase.Client
}

func (n OAuthClient) DeleteUser(ctx context.Context, id string) error {
	return n.AuthClient.DeleteUser(ctx, id)
}

func NewOAuthClient() (auth.OAuthClient, error) {
	serviceAccountPath := os.Getenv("FIREBASE_SERVICE_ACCOUNT")
	opt := option.WithCredentialsFile(serviceAccountPath)

	instance, err := firebase.NewApp(context.Background(), nil, opt)

	if err != nil {
		return nil, err
	}

	if err != nil {
		return nil, err
	}

	authClient, err := instance.Auth(context.Background())
	if err != nil {
		return nil, err
	}

	return OAuthClient{AuthClient: authClient}, nil
}

func (n OAuthClient) VerifyToken(ctx context.Context, token string) (string, error) {
	t, err := n.AuthClient.VerifyIDToken(ctx, token)
	if err != nil {
		return "", err
	}

	if t.Claims["id"] == nil || t.Claims["id"] == "" {
		return "", errors.New("")
	}

	return t.Claims["id"].(string), err
}

func (n OAuthClient) CreateUser(ctx context.Context, user auth.User, password string) (string, error) {
	u := &auth_firebase.UserToCreate{}
	u.UID(user.ID).Email(user.Email).DisplayName(user.Name).Password(password)
	record, err := n.AuthClient.CreateUser(ctx, u)
	if err != nil {
		return "", err
	}

	err = n.AuthClient.SetCustomUserClaims(ctx, record.UID, map[string]interface{}{
		"id": user.ID,
	})

	if err != nil {
		return "", err
	}

	return record.UID, nil
}
