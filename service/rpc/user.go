package rpc

import (
	"context"
	"fmt"

	"goshop/service-shop/model/user"
	"goshop/service-shop/pkg/db"

	"golang.org/x/crypto/bcrypt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/shinmigo/pb/basepb"
	"github.com/shinmigo/pb/shoppb"
)

type User struct {
}

func NewMUser() *User {
	return &User{}
}

func (u *User) Login(ctx context.Context, req *shoppb.LoginReq) (*shoppb.UserRes, error) {
	info, err := user.GetUserByUsername(req.Username)
	if err != nil {
		return nil, fmt.Errorf("用户不存在")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(info.Password), []byte(req.Password)); err != nil {
		return nil, fmt.Errorf("密码不正确")
	}

	return &shoppb.UserRes{
		UserId:   info.UserId,
		Username: info.Username,
		Name:     info.Name,
	}, nil
}

func (u *User) EditUser(ctx context.Context, req *shoppb.EditUserReq) (*basepb.AnyRes, error) {
	if _, err := user.GetOneByUserId(req.UserId); err != nil {
		return nil, err
	}

	aul := user.User{
		Password:  req.Password,
		Name:      req.Name,
		UpdatedBy: req.AdminId,
	}

	if err := db.Conn.Table(user.GetTableName()).Model(&user.User{UserId: req.UserId}).Updates(aul).Error; err != nil {
		return nil, err
	}

	if ctx.Err() == context.Canceled {
		return nil, status.Errorf(codes.Canceled, "The client canceled the request")
	}

	return &basepb.AnyRes{
		Id:    req.UserId,
		State: 1,
	}, nil
}
