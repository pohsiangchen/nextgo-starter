package main

import (
	"context"
	"database/sql"
	"log"

	"github.com/brianvoe/gofakeit/v7"
	_ "github.com/lib/pq"

	"apps/api/config"
	"apps/api/database"
	"apps/api/database/sqlc"
	"apps/api/util/auth"
)

var (
	USERS_SIZE        = 100
	FOLLOWS_MIN_SIZE  = USERS_SIZE / 10 * 4
	FOLLOWS_MAX_SIZE  = USERS_SIZE / 10 * 8
	POSTS_MIN_SIZE    = USERS_SIZE / 10 * 3
	POSTS_MAX_SIZE    = USERS_SIZE / 10 * 6
	COMMENTS_MIN_SIZE = USERS_SIZE / 10 * 3
	COMMENTS_MAX_SIZE = USERS_SIZE / 10 * 6
	USER_PASSWORD     = "testtest"
)

func main() {
	//
	// get configs
	//
	cfg := config.Get()

	//
	// get database connection and the instance of store layer
	//
	if cfg.Database.Driver == "" {
		log.Fatal("please fill in database credentials in .env file or set in environment variable")
	}
	dsn := database.DataSourceName(int(cfg.Database.Port), cfg.Database.Host, cfg.Database.User, cfg.Database.Password, cfg.Database.Name, cfg.Database.SslMode)
	db := database.NewDB(cfg.Database.Driver, dsn)
	db.SetMaxOpenConns(cfg.Database.MaxConnectionPool)
	db.SetMaxIdleConns(cfg.Database.MaxIdleConnections)
	db.SetConnMaxLifetime(cfg.Database.ConnectionsMaxLifeTime)
	store := sqlcstore.New(db)

	//
	// seed
	//
	ctx := context.Background()
	if err := seed(ctx, db, store); err != nil {
		log.Fatal("Seeding fail: ", err)
	} else {
		log.Println("Seeding complete")
	}
}

// TODO: speed up with Goroutines
func seed(ctx context.Context, db *sql.DB, queries *sqlcstore.Queries) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	qtx := queries.WithTx(tx)

	//
	// users
	//
	log.Println("Seeding users")
	userIDs := make([]int64, USERS_SIZE)
	for i := 0; i < USERS_SIZE; i++ {
		// replace with the following line to generate random password
		// hash, err := auth.HashPassword(gofakeit.Password(true, true, true, false, false, 8))
		hash, err := auth.HashPassword(USER_PASSWORD)
		if err != nil {
			return err
		}

		createdUser, err := qtx.CreateUser(
			ctx,
			sqlcstore.CreateUserParams{Email: gofakeit.Email(), Username: gofakeit.Username(), Password: hash},
		)
		if err != nil {
			return err
		}
		userIDs[i] = createdUser.ID
	}

	//
	// follows
	//
	log.Println("Seeding follows")
	for idx, userID := range userIDs {
		followsMaxSize := gofakeit.Number(FOLLOWS_MIN_SIZE, FOLLOWS_MAX_SIZE)
		noSelfUserIDs := make([]int64, len(userIDs))
		copy(noSelfUserIDs, userIDs)
		noSelfUserIDs = append(noSelfUserIDs[:idx], noSelfUserIDs[idx+1:]...)
		for i := 0; i < followsMaxSize; i++ {
			randIdx := gofakeit.Number(0, len(noSelfUserIDs)-1)
			if _, err := qtx.CreateFollow(
				ctx,
				sqlcstore.CreateFollowParams{
					FollowingUserID: userID,
					FollowedUserID:  noSelfUserIDs[randIdx],
				},
			); err != nil {
				return err
			} else {
				noSelfUserIDs = append(noSelfUserIDs[:randIdx], noSelfUserIDs[randIdx+1:]...)
			}
		}
	}

	//
	// posts
	//
	log.Println("Seeding posts and comments")
	for idx, userID := range userIDs {
		postsMaxSize := gofakeit.Number(POSTS_MIN_SIZE, POSTS_MAX_SIZE)
		for i := 0; i < postsMaxSize; i++ {
			if createdPost, err := qtx.CreatePost(
				ctx,
				sqlcstore.CreatePostParams{
					Title:   sql.NullString{String: gofakeit.HipsterSentence(gofakeit.Number(1, 10)), Valid: true},
					Content: sql.NullString{String: gofakeit.HipsterParagraph(gofakeit.Number(1, 10), gofakeit.Number(1, 10), gofakeit.Number(1, 10), ","), Valid: true},
					UserID:  userID,
				},
			); err != nil {
				return err
			} else {
				//
				// comments
				//
				commentsMaxSize := gofakeit.Number(COMMENTS_MIN_SIZE, COMMENTS_MAX_SIZE)
				noSelfUserIDs := make([]int64, len(userIDs))
				copy(noSelfUserIDs, userIDs)
				noSelfUserIDs = append(noSelfUserIDs[:idx], noSelfUserIDs[idx+1:]...)
				for j := 0; j < commentsMaxSize; j++ {
					randIdx := gofakeit.Number(0, len(noSelfUserIDs)-1)
					if _, err := qtx.CreateComment(
						ctx,
						sqlcstore.CreateCommentParams{
							Content: sql.NullString{String: gofakeit.Sentence(gofakeit.Number(1, 10)), Valid: true},
							UserID:  noSelfUserIDs[randIdx],
							PostID:  createdPost.ID,
						},
					); err != nil {
						return err
					} else {
						noSelfUserIDs = append(noSelfUserIDs[:randIdx], noSelfUserIDs[randIdx+1:]...)
					}
				}
			}
		}
	}

	return tx.Commit()
}
