package api

import (
	"context"
	"im/internal/user/service"
	"im/pkg/grpclib"
	"im/pkg/pb"
)

type UserExtServer struct{}

func (s *UserExtServer) SignIn(ctx context.Context, in *pb.SignInReq) (*pb.SignInResp, error) {
	userId, token, err := service.AuthService.SignIn(ctx, in.PhoneNumber, in.Code, in.DeviceId)
	if err != nil {
		return nil, err
	}
	return &pb.SignInResp{
		UserId: userId,
		Token:  token,
	}, nil
}

func (s *UserExtServer) GetUser(ctx context.Context, in *pb.GetUserReq) (*pb.GetUserResp, error) {
	userId, _, err := grpclib.GetCtxData(ctx)
	if err != nil {
		return nil, err
	}
	user, err := service.UserService.Get(ctx, userId)
	if err != nil {
		return nil, err
	}
	return &pb.GetUserResp{
		User: &pb.User{
			UserId:     user.Id,
			Nickname:   user.Nickname,
			Sex:        user.Sex,
			AvatarUrl:  user.AvatarUrl,
			Extra:      user.Extra,
			CreateTime: user.CreateTime.Unix(),
			UpdateTime: user.UpdateTime.Unix(),
		},
	}, nil
}
