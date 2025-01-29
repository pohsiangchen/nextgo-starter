package post

import (
	"context"
	"database/sql"
	"errors"
	"math"

	"apps/api/database/sqlc"
	"apps/api/util/auth"
)

type PostService interface {
	Create(ctx context.Context, post *CreatePostRequest) (sqlcstore.Post, error)
	Update(ctx context.Context, post *UpdatePostRequest) (sqlcstore.Post, error)
	Delete(ctx context.Context, postID int64) error
	FindById(ctx context.Context, postID int64) (sqlcstore.Post, error)
	ListFeeds(ctx context.Context) ([]sqlcstore.ListFeedsByUserIdRow, error)
}

type PostServiceImpl struct {
	store         *sqlcstore.Queries
	authenticator auth.Authenticator
}

func NewPostServiceImpl(store *sqlcstore.Queries, authenticator auth.Authenticator) (service PostService) {
	return &PostServiceImpl{
		store:         store,
		authenticator: authenticator,
	}
}

func (ps PostServiceImpl) Create(ctx context.Context, post *CreatePostRequest) (sqlcstore.Post, error) {
	user, ok := ps.authenticator.UserFromCtx(ctx)
	if !ok {
		return sqlcstore.Post{}, errors.New("cannot find current signed-in user from context")
	}
	createdPost, err := ps.store.CreatePost(
		ctx,
		sqlcstore.CreatePostParams{
			Title:   sql.NullString{String: post.Title, Valid: post.Title != ""},
			Content: sql.NullString{String: post.Content, Valid: post.Content != ""},
			UserID:  user.ID,
		},
	)
	return createdPost, err
}

func (ps PostServiceImpl) Update(ctx context.Context, post *UpdatePostRequest) (sqlcstore.Post, error) {
	updatedPost, err := ps.store.UpdatePost(
		ctx,
		sqlcstore.UpdatePostParams{
			ID:      post.ID,
			Title:   sql.NullString{String: post.Title, Valid: post.Title != ""},
			Content: sql.NullString{String: post.Content, Valid: post.Content != ""},
		},
	)
	return updatedPost, err
}

func (ps PostServiceImpl) Delete(ctx context.Context, postID int64) error {
	return ps.store.DeletePost(ctx, postID)
}

func (ps PostServiceImpl) FindById(ctx context.Context, postID int64) (sqlcstore.Post, error) {
	post, err := ps.store.GetPost(ctx, postID)
	return post, err
}

func (ps PostServiceImpl) ListFeeds(ctx context.Context) ([]sqlcstore.ListFeedsByUserIdRow, error) {
	user, ok := ps.authenticator.UserFromCtx(ctx)
	if !ok {
		return []sqlcstore.ListFeedsByUserIdRow{}, errors.New("cannot find current signed-in user from context")
	}
	rows, err := ps.store.ListFeedsByUserId(ctx, sqlcstore.ListFeedsByUserIdParams{
		FollowingUserID: user.ID,
		LastPostsID:     math.MaxInt64,
		Limit:           100,
	})
	return rows, err
}
