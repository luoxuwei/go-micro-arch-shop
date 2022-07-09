package handler

import (
	"crypto/sha512"
	"fmt"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/codes"
	"github.com/anaskhan96/go-password-encoder"

	"context"

	"shop-srvs/user-srv/global"
	"shop-srvs/user-srv/model"
	"shop-srvs/user-srv/proto"
)

type UserServer struct {
	proto.UnimplementedUserServer
}

func ModelToRsponse(user model.User) proto.UserInfoResponse{
	//在grpc的message中字段有默认值，你不能随便赋值nil进去，容易出错
	//这里要搞清， 哪些字段是有默认值
	userInfoRsp := proto.UserInfoResponse{
		Id:       user.ID,
		PassWord: user.Password,
		NickName: user.NickName,
		Gender: user.Gender,
		Role: int32(user.Role),
		Mobile: user.Mobile,
	}
	if user.Birthday != nil {
		userInfoRsp.BirthDay = uint64(user.Birthday.Unix())
	}
	return userInfoRsp
}

func (s *UserServer) CreateUser(c context.Context, req *proto.CreateUserInfo) (*proto.UserInfoResponse, error){
	user := model.User{
		Mobile: req.Mobile,
	}

	result := global.DB.Where(&model.User{
		Mobile: req.Mobile,
	}).First(&user)

	if result.RowsAffected == 1 {
		return nil, status.Errorf(codes.AlreadyExists, "用户已存在")
	}

	user.Mobile = req.Mobile
	user.NickName = req.NickName

	//密码加密
	options := &password.Options{16, 100, 32, sha512.New}
	salt, encodedPwd := password.Encode(req.PassWord, options)
	user.Password = fmt.Sprintf("$pbkdf2-sha512$%s$%s", salt, encodedPwd)

	result = global.DB.Create(&user)
	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, result.Error.Error())
	}

	userInfoRsp := ModelToRsponse(user)
	return &userInfoRsp, nil
}