package auth

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog"
	"golang.org/x/crypto/bcrypt"

	"apps/api/middleware"
	"apps/api/util/response"
)

type AuthController struct {
	authService AuthService
	validator   *validator.Validate
}

func NewAuthController(service AuthService, validator *validator.Validate) (ctrl *AuthController, err error) {
	if validator == nil {
		return nil, errors.New("validator instance cannot be nil")
	}
	return &AuthController{authService: service, validator: validator}, err
}

func (authCtrl *AuthController) CreateToken(w http.ResponseWriter, r *http.Request) {
	zlog := zerolog.Ctx(r.Context())

	data, ok := (middleware.Object(r.Context())).(*CreateTokenRequest)
	if !ok {
		render.Render(w, r, response.ErrMissingBindedReqObj)
		return
	}

	token, err := authCtrl.authService.CreateToken(r.Context(), data)
	if err != nil {
		switch err {
		case sql.ErrNoRows, bcrypt.ErrMismatchedHashAndPassword:
			render.Render(w, r, response.ErrUnauthorized())
		default:
			zlog.Error().Err(err).Msg("failed to create authentication token")
			render.Render(w, r, response.ErrInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusCreated)
	render.Render(w, r, NewAuthTokenResponse(token))
}
