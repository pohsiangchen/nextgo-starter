package user

import (
	"context"
	"database/sql"
	"strings"

	"golang.org/x/crypto/bcrypt"

	"apps/api/database"
	"apps/api/database/sqlc"
)

type UserService interface {
	Create(ctx context.Context, user *CreateUserRequest) (sqlcstore.User, error)
	Update(ctx context.Context, user *UpdateUserRequest) (sqlcstore.User, error)
	UpdatePassword(ctx context.Context, user *UpdateUserPasswordRequest) (sqlcstore.User, error)
	Delete(ctx context.Context, userID int64) error
	FindById(ctx context.Context, userID int64) (sqlcstore.User, error)
	FindAll(ctx context.Context) ([]sqlcstore.User, error)
}

type UserServiceImpl struct {
	store *sqlcstore.Queries
}

func NewUserServiceImpl(store *sqlcstore.Queries) (service UserService) {
	return &UserServiceImpl{
		store: store,
	}
}

func (us UserServiceImpl) Create(ctx context.Context, user *CreateUserRequest) (sqlcstore.User, error) {
	hash, err := HashPassword(user.Password)
	if err != nil {
		return sqlcstore.User{}, err
	}

	createdUser, err := us.store.CreateUser(
		ctx,
		sqlcstore.CreateUserParams{Email: user.Email, Username: user.Username, Password: hash},
	)
	if err != nil && strings.Contains(err.Error(), `duplicate key value violates unique constraint "users_email_key"`) {
		return createdUser, database.ErrDuplicateEmail
	}
	return createdUser, err
}

func (us UserServiceImpl) Update(ctx context.Context, user *UpdateUserRequest) (sqlcstore.User, error) {
	if foundUser, err := us.FindById(ctx, user.ID); err != nil {
		return foundUser, err
	}

	updatedUser, err := us.store.UpdateUser(
		ctx,
		sqlcstore.UpdateUserParams{ID: user.ID, Username: sql.NullString{String: user.Username, Valid: user.Username != ""}},
	)
	return updatedUser, err
}

// TDOD: update user password
func (us UserServiceImpl) UpdatePassword(ctx context.Context, user *UpdateUserPasswordRequest) (sqlcstore.User, error) {
	if foundUser, err := us.FindById(ctx, user.ID); err != nil {
		return foundUser, err
	}

	hash, err := HashPassword(user.Password)
	if err != nil {
		return sqlcstore.User{}, err
	}

	return us.store.UpdateUserPassword(
		ctx,
		sqlcstore.UpdateUserPasswordParams{ID: user.ID, Password: hash},
	)
}

func (us UserServiceImpl) Delete(ctx context.Context, userID int64) error {
	if _, err := us.FindById(ctx, userID); err != nil {
		return err
	}
	return us.store.DeleteUser(ctx, userID)
}

func (us UserServiceImpl) FindById(ctx context.Context, userID int64) (sqlcstore.User, error) {
	user, err := us.store.GetUser(ctx, userID)
	if err != nil && err == sql.ErrNoRows {
		return user, database.ErrRecordNotFound
	}
	return user, err
}

func (us UserServiceImpl) FindAll(ctx context.Context) ([]sqlcstore.User, error) {
	users, err := us.store.ListUsers(ctx)
	return users, err
}

// Hash password
func HashPassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

// Compares a bcrypt hashed password with the given password
func ComparePassword(hashedPassword []byte, password string) error {
	return bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
}
