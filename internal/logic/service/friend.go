package service

import (
	"context"
	"im/internal/logic/dao"
	"im/internal/logic/model"
	"im/pkg/gerrors"
	"im/pkg/pb"
	"im/pkg/rpc"
)

type friendService struct{}

var FriendService = new(friendService)

// GetUsers 获取群组用户
func (s *friendService) List(ctx context.Context, groupId int64) ([]*pb.Friend, error) {
	friends, err := dao.FriendDao.List(groupId)
	if err != nil {
		return nil, err
	}

	userIds := make([]int64, len(friends))
	for i := range friends {
		userIds[i] = friends[i].FriendId
	}
	resp, err := rpc.UserIntClient.GetUsers(ctx, &pb.GetUsersReq{UserIds: userIds})
	if err != nil {
		return nil, err
	}

	var infos = make([]*pb.Friend, len(friends))
	for i := range friends {
		friend := pb.Friend{
			UserId:  friends[i].FriendId,
			Remarks: friends[i].Remarks,
			Extra:   friends[i].Extra,
		}

		user, ok := resp.Users[friends[i].FriendId]
		if ok {
			friend.Nickname = user.Nickname
			friend.Sex = user.Sex
			friend.AvatarUrl = user.AvatarUrl
			friend.UserExtra = user.Extra
		}
		infos[i] = &friend
	}

	return infos, nil
}

func (*friendService) AddFriend(ctx context.Context, userId, friendId int64, remarks, description string) error {
	err := dao.FriendDao.Add(model.Friend{
		UserId:   userId,
		FriendId: friendId,
		Remarks:  remarks,
		Status:   model.FriendStatusApply,
	})
	if err != nil {
		return err
	}

	resp, err := rpc.UserIntClient.GetUser(ctx, &pb.GetUserReq{UserId: userId})
	if err != nil {
		return err
	}

	PushService.Push(ctx, friendId, pb.PushCode_PC_ADD_FRIEND, &pb.AddFriendPush{
		FriendId:    userId,
		Nickname:    resp.User.Nickname,
		AvatarUrl:   resp.User.AvatarUrl,
		Description: description,
	}, true)
	return nil
}

func (*friendService) AgreeAddFriend(ctx context.Context, userId, friendId int64, remarks string) error {
	friend, err := dao.FriendDao.Get(friendId, userId)
	if err != nil {
		return err
	}
	if friend == nil {
		return gerrors.ErrBadRequest
	}

	if friend.Status == model.FriendStatusAgree {
		return nil
	}

	err = dao.FriendDao.UpdateStatus(friendId, userId, model.FriendStatusAgree)
	if err != nil {
		return err
	}

	err = dao.FriendDao.Add(model.Friend{
		UserId:   userId,
		FriendId: friendId,
		Remarks:  remarks,
		Status:   model.FriendStatusApply,
	})
	if err != nil {
		return err
	}

	resp, err := rpc.UserIntClient.GetUser(ctx, &pb.GetUserReq{UserId: userId})
	if err != nil {
		return err
	}

	PushService.Push(ctx, friendId, pb.PushCode_PC_AGREE_ADD_FRIEND, &pb.AgreeAddFriendPush{
		FriendId:  userId,
		Nickname:  resp.User.Nickname,
		AvatarUrl: resp.User.AvatarUrl,
	}, true)
	return nil
}
