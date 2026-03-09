package utils

import (
	"context"
	"database/sql"
	"errors"

	"github.com/BlitzStudio/blitzStudioAuth/out/repository"
	"github.com/BlitzStudio/blitzStudioAuth/types"
)

var log = GetLogger()

func createUser(user types.User, db *sql.DB) (int64, error) {
	repo := repository.New(db)
	ctx := context.Background()
	createdUser, err := repo.CreateUser(ctx, repository.CreateUserParams{
		Name:     user.Name,
		Email:    user.Email,
		Password: user.Password,
	})
	if err != nil {
		log.Error(err)
		return -1, err
	}
	return createdUser.LastInsertId()
}

func CreateUser(user types.User, db *sql.DB) (int64, error) {
	ctx := context.Background()
	repo := repository.New(db)
	userCount, err := repo.CountUserByEmail(ctx, user.Email)

	if err != nil {
		log.Error("Error counting user by email\n" + err.Error())
		return -1, err
	}

	if userCount != 0 {
		return -1, errors.New("User with: " + user.Email + " already exists")
	}

	user.Password, err = GenerateHash(user.Password)
	if err != nil {
		return -1, err
	}

	createdUser, err := createUser(user, db)
	if err != nil {
		return -1, errors.New("500")
	}

	return createdUser, nil
}
