package auth

import (
	"context"

	"apps/api/database/sqlc"
	"apps/api/util/auth"
)

type AuthService interface {
	CreateToken(ctx context.Context, ctr *CreateTokenRequest) (string, error)
}

type AuthServiceImpl struct {
	store         *sqlcstore.Queries
	authenticator auth.Authenticator
}

func NewAuthServiceImpl(store *sqlcstore.Queries, authenticator auth.Authenticator) (service AuthService) {
	return &AuthServiceImpl{
		store:         store,
		authenticator: authenticator,
	}
}

func (as AuthServiceImpl) CreateToken(ctx context.Context, ctr *CreateTokenRequest) (string, error) {
	// checks if the input's email exists in database
	user, err := as.store.GetUserByEmail(ctx, ctr.Email)
	if err != nil {
		return "", err
	}

	// checks if the input's password is same as the one in database
	if err := auth.ComparePassword(user.Password, ctr.Password); err != nil {
		return "", err
	}

	// generates JWT token
	claims := as.authenticator.CreateClaims(user.ID)
	token, err := as.authenticator.CreateToken(claims)
	if err != nil {
		return "", err
	}

	return token, err
}
