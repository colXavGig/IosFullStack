package models

import (
	"context"
	"errors"
	"fmt"
	"ios_full_stack/data"
	"ios_full_stack/dto"
	"log"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Errors
var (
	ErrDuplicatedUser error = errors.New("user already exists")
	ErrPasswdTooLong  error = errors.New("password too long")
	ErrUserNotFound   error = errors.New("user not found")
)

type AppUser struct {
	dto.User
}

func TryFindUserById(ctx context.Context, id uint) (*AppUser, error) {
	tx := data.GetTransaction(ctx)

	userData, err := gorm.G[dto.User](tx).Where("ID = ?", id).First(ctx)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrUserNotFound
	} else if err != nil {
		log.Panicln(err)
	}

	return &AppUser{User: userData}, nil
}

func TryFindUserByUsername(ctx context.Context, username string) (*AppUser, bool) {
	tx := data.GetTransaction(ctx)

	userData, err := gorm.G[dto.User](tx).Where("username = ?", username).First(ctx)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, false
	} else if err != nil {
		log.Panicln(err)
	}

	return &AppUser{
		User: userData,
	}, true
}

func SearchUsersByUsername(ctx context.Context, query string) ([]dto.User, error) {
	var (
		tx    = data.GetTransaction(ctx)
		users []dto.User
	)

	err := tx.Where("username LIKE ?", "%"+query+"%").Find(&users).Error
	if err != nil {
		return nil, err
	}

	return users, nil
}

func RegisterNewAppUser(ctx context.Context, newUser *dto.User) error {
	hashedPasswd, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), 10)
	if errors.Is(err, bcrypt.ErrPasswordTooLong) {
		return ErrPasswdTooLong
	} else if err != nil {
		log.Panicln(fmt.Errorf("unable to hash the password: %w", err))
	}

	newUser.Password = string(hashedPasswd)

	tx := data.GetTransaction(ctx)

	err = gorm.G[dto.User](tx).Create(ctx, newUser)
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return ErrDuplicatedUser
	} else if err != nil {
		log.Fatalln(fmt.Errorf("was not able to create user entry in db: %w", err))
	}

	return nil
}

func (u *AppUser) IsCorrectPassword(passwd string) bool {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(passwd)) == nil
}

func (u *AppUser) UpdatePassword(ctx context.Context, newPasswd string) error {

	tx := data.GetTransaction(ctx)

	hashedPasswd, err := bcrypt.GenerateFromPassword([]byte(newPasswd), 10)
	if err != nil {
		log.Panicln(fmt.Errorf("unable to hash the password: %w", err))
	}

	rowsAffected, err := gorm.G[dto.User](tx).Where("username = ?", u.Username).Update(ctx, "password", hashedPasswd)
	if rowsAffected < 1 {
		return errors.New("user not found")
	} else if rowsAffected > 1 {
		tx.Rollback()
		return errors.New("user not found")
	} else if err != nil {
		return fmt.Errorf("unable to update password: %w", err)
	}

	return nil
}
