package post

import (
	"net/http"

	"apps/api/database/sqlc"
)

type User struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
}

type PostResponse struct {
	ID      int64  `json:"id"`
	Title   string `json:"title,omitempty"`
	Content string `json:"content,omitempty"`
	User    *User  `json:"user,omitempty"`
}

type Feed struct {
	Post          *PostResponse `json:"post"`
	CommentsCount int           `json:"comments_count"`
}

type FeedResponse struct {
	Data     []*Feed `json:"data"`
	Metadata *Filter `json:"metadata"`
}

func (pr *PostResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (fr *FeedResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func NewPostResponse(p sqlcstore.Post) *PostResponse {
	return &PostResponse{
		ID:      p.ID,
		Title:   p.Title.String,
		Content: p.Content.String,
	}
}

func NewFeedListResponse(feeds []sqlcstore.ListFeedsByUserIdRow, filters *Filter) *FeedResponse {
	data := []*Feed{}
	for _, feed := range feeds {
		data = append(data, &Feed{
			Post: &PostResponse{
				ID:      feed.Post.ID,
				Title:   feed.Post.Title.String,
				Content: feed.Post.Content.String,
				User: &User{
					ID:       feed.Post.UserID,
					Username: feed.Username.String,
				},
			},
			CommentsCount: int(feed.CommentsCount),
		})
	}
	return &FeedResponse{
		Data:     data,
		Metadata: filters,
	}
}
