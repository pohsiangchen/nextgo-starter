package user

import (
	"database/sql"
	"errors"
	"net/http"
	"strings"

	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog"

	"apps/api/middleware"
	"apps/api/util/response"
)

type UserController struct {
	userService UserService
	validator   *validator.Validate
}

func NewUserController(service UserService, validator *validator.Validate) (ctrl *UserController, err error) {
	if validator == nil {
		return nil, errors.New("validator instance cannot be nil")
	}
	return &UserController{userService: service, validator: validator}, err
}

func (userCtrl *UserController) Create(w http.ResponseWriter, r *http.Request) {
	zlog := zerolog.Ctx(r.Context())

	data, ok := (middleware.Object(r.Context())).(*CreateUserRequest)
	if !ok {
		render.Render(w, r, response.ErrMissingBindedReqObj)
		return
	}

	user, err := userCtrl.userService.Create(r.Context(), data)
	if err != nil {
		if err != nil && strings.Contains(err.Error(), `duplicate key value violates unique constraint "users_email_key"`) {
			render.Render(w, r, response.ErrDuplicateEmail(data.Email))
			return
		}
		zlog.Error().Err(err).Msg("failed to create a user record")
		render.Render(w, r, response.ErrInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	render.Render(w, r, NewUserResponse(user))
}

func (userCtrl *UserController) Update(w http.ResponseWriter, r *http.Request) {
	zlog := zerolog.Ctx(r.Context())
	userID := UserIDFromCtx(r.Context())

	data, ok := (middleware.Object(r.Context())).(*UpdateUserRequest)
	if !ok {
		render.Render(w, r, response.ErrMissingBindedReqObj)
		return
	}
	data.ID = userID

	user, err := userCtrl.userService.Update(r.Context(), data)
	if err != nil {
		if err == sql.ErrNoRows {
			render.Render(w, r, response.ErrNotFound)
			return
		}
		zlog.Error().Err(err).Msg("failed to update a user record")
		render.Render(w, r, response.ErrInternalServerError)
		return
	}

	render.Render(w, r, NewUserResponse(user))
}

func (userCtrl *UserController) Delete(w http.ResponseWriter, r *http.Request) {
	zlog := zerolog.Ctx(r.Context())
	userID := UserIDFromCtx(r.Context())
	if err := userCtrl.userService.Delete(r.Context(), userID); err != nil {
		if err == sql.ErrNoRows {
			render.Render(w, r, response.ErrNotFound)
			return
		}
		zlog.Error().Err(err).Msgf("failed to delete a user record with user ID '%d'", userID)
		render.Render(w, r, response.ErrInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(nil)
}

func (userCtrl *UserController) List(w http.ResponseWriter, r *http.Request) {
	zlog := zerolog.Ctx(r.Context())
	users, err := userCtrl.userService.FindAll(r.Context())
	if err != nil {
		zlog.Error().Err(err).Msg("failed to retrieve user records")
		render.Render(w, r, response.ErrInternalServerError)
		return
	}
	render.RenderList(w, r, NewUserListResponse(users))
}

func (userCtrl *UserController) Get(w http.ResponseWriter, r *http.Request) {
	zlog := zerolog.Ctx(r.Context())
	userID := UserIDFromCtx(r.Context())
	user, err := userCtrl.userService.FindById(r.Context(), userID)
	if err != nil {
		if err == sql.ErrNoRows {
			render.Render(w, r, response.ErrNotFound)
			return
		}
		zlog.Error().Err(err).Msgf("failed to retrieve a user record with user ID '%d'", userID)
		render.Render(w, r, response.ErrInternalServerError)
		return
	}
	render.Render(w, r, NewUserResponse(user))
}
