package database

import (
	"fmt"
	"time"
)

type RefreshToken struct {
	UserId    int       `json:"user_id"`
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
}

func (db *DB) SaveRefreshToken(userId int, refreshToken string) error {
	dbStructure, err := db.loadDB()
	if err != nil {
		return err
	}

	updatedRefreshToken := RefreshToken{
		UserId:    userId,
		Token:     refreshToken,
		ExpiresAt: time.Now().Add(time.Hour),
	}
	dbStructure.RefreshTokens[refreshToken] = updatedRefreshToken

	err = db.writeDB(dbStructure)
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) RevokeRefreshToken(refreshToken string) error {
	dbStructure, err := db.loadDB()
	if err != nil {
		return err
	}
	fmt.Printf("before delete - DB: %v\n", dbStructure.RefreshTokens)
	delete(dbStructure.RefreshTokens, refreshToken)
	fmt.Printf("after delete - DB: %v\n", dbStructure.RefreshTokens)
	err = db.writeDB(dbStructure)
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) GetUserForRefreshToken(refresh string) (User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	refreshToken, ok := dbStructure.RefreshTokens[refresh]
	if !ok {
		return User{}, nil
	}

	if refreshToken.ExpiresAt.Before(time.Now()) {
		return User{}, ErrNotExist
	}

	user, err := db.GetUser(refreshToken.UserId)
	if err != nil {
		return User{}, nil
	}
	return user, nil
}
