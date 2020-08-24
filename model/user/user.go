package user

import (
	"fmt"

	"goshop/service-shop/pkg/db"
	"goshop/service-shop/pkg/utils"

	"github.com/jinzhu/gorm"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	UserId    uint64 `json:"user_id" gorm:"PRIMARY_KEY"`
	Username  string
	Password  string
	Name      string
	CreatedBy uint64
	UpdatedBy uint64
	CreatedAt utils.JSONTime
	UpdatedAt utils.JSONTime
}

func GetTableName() string {
	return "user"
}

func GetField() []string {
	return []string{
		"user_id", "username", "password", "name",
		"created_by", "updated_by", "created_at", "updated_at",
	}
}

func (u *User) BeforeSave(scope *gorm.Scope) (err error) {
	if len(u.Password) > 0 {
		if pw, err := bcrypt.GenerateFromPassword([]byte(u.Password), 0); err == nil {
			scope.SetColumn("password", pw)
		}
	}
	return
}

func GetOneByUserId(userId uint64) (*User, error) {
	if userId == 0 {
		return nil, fmt.Errorf("user_id is null")
	}
	row := &User{}
	err := db.Conn.Table(GetTableName()).
		Select(GetField()).
		Where("user_id = ?", userId).
		First(row).Error

	if err != nil {
		return nil, fmt.Errorf("err: %v", err)
	}
	return row, nil
}

func GetUserByUsername(username string) (*User, error) {
	if len(username) == 0 {
		return nil, fmt.Errorf("username is null")
	}
	row := &User{}
	err := db.Conn.Table(GetTableName()).
		Select(GetField()).
		Where("username = ?", username).
		First(row).Error

	if err != nil {
		return nil, fmt.Errorf("err: %v", err)
	}
	return row, nil
}
