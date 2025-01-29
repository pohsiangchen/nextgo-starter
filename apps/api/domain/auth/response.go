package auth

import (
	"net/http"
)

type AuthTokenResponse struct {
	Token string `json:"token"`
}

func (atr AuthTokenResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func NewAuthTokenResponse(token string) AuthTokenResponse {
	return AuthTokenResponse{
		Token: token,
	}
}
