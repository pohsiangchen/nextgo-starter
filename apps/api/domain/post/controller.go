package post

import (
	"errors"
	"net/http"

	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog"

	"apps/api/middleware"
	"apps/api/util/response"
)

type PostController struct {
	postService PostService
	validator   *validator.Validate
}

func NewPostController(service PostService, validator *validator.Validate) (ctrl *PostController, err error) {
	if validator == nil {
		return nil, errors.New("validator instance cannot be nil")
	}
	return &PostController{postService: service, validator: validator}, err
}

func (postCtrl *PostController) Create(w http.ResponseWriter, r *http.Request) {
	zlog := zerolog.Ctx(r.Context())

	data, ok := (middleware.Object(r.Context())).(*CreatePostRequest)
	if !ok {
		render.Render(w, r, response.ErrMissingBindedReqObj)
		return
	}

	post, err := postCtrl.postService.Create(r.Context(), data)
	if err != nil {
		zlog.Error().Err(err).Msg("failed to create a post")
		render.Render(w, r, response.ErrInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	render.Render(w, r, NewPostResponse(post))
}

func (postCtrl *PostController) Update(w http.ResponseWriter, r *http.Request) {
	zlog := zerolog.Ctx(r.Context())
	post := PostFromCtx(r.Context())

	data, ok := (middleware.Object(r.Context())).(*UpdatePostRequest)
	if !ok {
		render.Render(w, r, response.ErrMissingBindedReqObj)
		return
	}
	data.ID = post.ID

	post, err := postCtrl.postService.Update(r.Context(), data)
	if err != nil {
		zlog.Error().Err(err).Msgf("failed to update the post with ID: %d", data.ID)
		render.Render(w, r, response.ErrInternalServerError)
		return
	}

	render.Render(w, r, NewPostResponse(post))
}

func (postCtrl *PostController) Delete(w http.ResponseWriter, r *http.Request) {
	zlog := zerolog.Ctx(r.Context())
	post := PostFromCtx(r.Context())

	if err := postCtrl.postService.Delete(r.Context(), post.ID); err != nil {
		zlog.Error().Err(err).Msgf("failed to delete the post with ID '%d'", post.ID)
		render.Render(w, r, response.ErrInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(nil)
}

func (postCtrl *PostController) Get(w http.ResponseWriter, r *http.Request) {
	post := PostFromCtx(r.Context())
	render.Render(w, r, NewPostResponse(post))
}

func (postCtrl *PostController) ListFeeds(w http.ResponseWriter, r *http.Request) {
	zlog := zerolog.Ctx(r.Context())
	feeds, err := postCtrl.postService.ListFeeds(r.Context())
	if err != nil {
		zlog.Error().Err(err).Msg("failed to retrieve feeds")
		render.Render(w, r, response.ErrInternalServerError)
		return
	}
	render.RenderList(w, r, NewFeedListResponse(feeds))
}
