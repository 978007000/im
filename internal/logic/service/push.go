package service

import (
	"context"
	"im/internal/logic/model"
	"im/pkg/gerrors"
	"im/pkg/grpclib"
	"im/pkg/logger"
	"im/pkg/pb"
	"time"

	"github.com/golang/protobuf/proto"

	"go.uber.org/zap"
)

type pushService struct{}

var PushService = new(pushService)

func (s *pushService) Push(ctx context.Context, userId int64, code pb.PushCode, message proto.Message, isPersist bool) error {
	logger.Logger.Debug("push",
		zap.Int64("request_id", grpclib.GetCtxRequstId(ctx)),
		zap.Int64("ser_id", userId),
		zap.Any("message", message))

	messageBuf, err := proto.Marshal(message)
	if err != nil {
		return gerrors.WrapError(err)
	}

	commandBuf, err := proto.Marshal(&pb.Command{Code: int32(code), Data: messageBuf})
	if err != nil {
		return gerrors.WrapError(err)
	}

	MessageService.SendToUser(ctx,
		model.Sender{
			SenderType: pb.SenderType_ST_SYSTEM,
			SenderId:   0,
			DeviceId:   0,
		},
		userId,
		0,
		pb.SendMessageReq{
			ReceiverType:   pb.ReceiverType_RT_USER,
			ReceiverId:     userId,
			ToUserIds:      nil,
			MessageType:    pb.MessageType_MT_COMMAND,
			MessageContent: commandBuf,
			SendTime:       time.Now().Unix(),
			IsPersist:      isPersist,
		},
	)
	return nil
}

func (s *pushService) PushToGroup(ctx context.Context, groupId int64, code pb.PushCode, message proto.Message, isPersist bool) error {
	logger.Logger.Debug("push",
		zap.Int64("request_id", grpclib.GetCtxRequstId(ctx)),
		zap.Int64("group_id", groupId),
		zap.Any("message", message))

	messageBuf, err := proto.Marshal(message)
	if err != nil {
		return gerrors.WrapError(err)
	}

	commandBuf, err := proto.Marshal(&pb.Command{Code: int32(code), Data: messageBuf})
	if err != nil {
		return gerrors.WrapError(err)
	}

	MessageService.SendToGroup(ctx,
		model.Sender{
			SenderType: pb.SenderType_ST_SYSTEM,
			SenderId:   0,
			DeviceId:   0,
		},
		pb.SendMessageReq{
			ReceiverType:   pb.ReceiverType_RT_NORMAL_GROUP,
			ReceiverId:     groupId,
			ToUserIds:      nil,
			MessageType:    pb.MessageType_MT_COMMAND,
			MessageContent: commandBuf,
			SendTime:       time.Now().Unix(),
			IsPersist:      isPersist,
		},
	)
	return nil
}
