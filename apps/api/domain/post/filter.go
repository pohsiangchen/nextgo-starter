package post

import (
	"math"
	"net/url"
	"strconv"
)

var (
	DEFAULT_LIMIT = 20
)

type Filter struct {
	Limit      int   `json:"limit" validate:"gte=1,lte=50"`
	LastPostID int64 `json:"last_post_id" validate:"gte=0"`
}

func NewFilter() *Filter {
	return &Filter{
		Limit:      DEFAULT_LIMIT,
		LastPostID: math.MaxInt64,
	}
}

func (f *Filter) Parse(qs url.Values) *Filter {
	if limit, err := strconv.Atoi(qs.Get("limit")); err == nil {
		f.Limit = limit
	}

	if lastPostID, err := strconv.ParseInt(qs.Get("last_post_id"), 10, 64); err == nil {
		f.LastPostID = lastPostID
	}

	return f
}
