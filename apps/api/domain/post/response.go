package post

import (
	"net/http"

	"apps/api/database/sqlc"

	"github.com/go-chi/render"
)

type User struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
}

type PostResponse struct {
	ID      int64  `json:"id"`
	Title   string `json:"title,omitempty"`
	Content string `json:"content,omitempty"`
	User    User   `json:"user,omitempty"`
}

type FeedResponse struct {
	Post          PostResponse `json:"post"`
	CommentsCount int          `json:"comments_count"`
}

func (hr PostResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (hr FeedResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func NewPostResponse(p sqlcstore.Post) PostResponse {
	return PostResponse{
		ID:      p.ID,
		Title:   p.Title.String,
		Content: p.Content.String,
	}
}

func NewFeedListResponse(feeds []sqlcstore.ListFeedsByUserIdRow) []render.Renderer {
	list := []render.Renderer{}
	for _, feed := range feeds {
		list = append(list, FeedResponse{
			Post: PostResponse{
				ID:      feed.Post.ID,
				Title:   feed.Post.Title.String,
				Content: feed.Post.Content.String,
				User: User{
					ID:       feed.Post.UserID,
					Username: feed.Username.String,
				},
			},
			CommentsCount: int(feed.CommentsCount),
		})
	}
	return list
}
