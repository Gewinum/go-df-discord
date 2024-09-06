package server

import (
	"errors"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type User struct {
	Discord string
	XUID    string
}

type Repository interface {
	GetUserByDiscord(discordId string) (*User, error)
	GetUserByXUID(xuid string) (*User, error)
	CreateUser(discordId, xuid string) (*User, error)
	DeleteUserByDiscord(discordId string) error
	DeleteUserByXUID(xuid string) error
}

type UserData struct {
	*gorm.Model
	Discord string
	XUID    string `gorm:"column:xuid"`
}

func (u *UserData) ToUser() *User {
	return &User{
		Discord: u.Discord,
		XUID:    u.XUID,
	}
}

type defaultRepository struct {
	db *gorm.DB
}

func NewDefaultRepository() (Repository, error) {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	err = db.AutoMigrate(&UserData{})
	if err != nil {
		return nil, err
	}

	return &defaultRepository{
		db: db,
	}, nil
}

func (r *defaultRepository) GetUserByDiscord(discordId string) (*User, error) {
	var user UserData
	err := r.db.First(&user, "discord = ?", discordId).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, NewApplicationError(40400, "User not found")
		}
		return nil, err
	}
	return user.ToUser(), nil
}

func (r *defaultRepository) GetUserByXUID(xuid string) (*User, error) {
	var user UserData
	err := r.db.First(&user, "xuid = ?", xuid).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, NewApplicationError(40400, "User not found")
		}
		return nil, err
	}
	return user.ToUser(), nil
}

func (r *defaultRepository) CreateUser(discordId, xuid string) (*User, error) {
	var user UserData
	err := r.db.First(&user, "discord = ? OR xuid = ?", discordId, xuid).Error
	if err == nil {
		return nil, NewApplicationError(40000, "Either discord or XUID are already bound")
	}
	user.Discord = discordId
	user.XUID = xuid
	err = r.db.Create(&user).Error
	if err != nil {
		return nil, err
	}
	return user.ToUser(), nil
}

func (r *defaultRepository) DeleteUserByDiscord(discordId string) error {
	user, err := r.GetUserByDiscord(discordId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return NewApplicationError(40400, "User not found")
		}
		return err
	}
	r.db.Delete(&UserData{}, "discord = ?", user.Discord)
	return nil
}

func (r *defaultRepository) DeleteUserByXUID(xuid string) error {
	user, err := r.GetUserByXUID(xuid)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return NewApplicationError(40400, "User not found")
		}
		return err
	}
	r.db.Delete(&UserData{}, "xuid = ?", user.XUID)
	return nil
}
