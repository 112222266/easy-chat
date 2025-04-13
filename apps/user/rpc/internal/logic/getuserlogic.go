package logic

import (
	"context"
	"easy-chat/apps/user/models"
	"errors"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
	"strconv"

	"easy-chat/apps/user/rpc/internal/svc"
	"easy-chat/apps/user/rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserLogic {
	return &GetUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUserLogic) GetUser(in *user.GetUserReq) (*user.GetUserResp, error) {
	logger := l.svcCtx.Logger
	var userModel models.User
	if err := l.svcCtx.DB.Take(&userModel, in.Id).Error; err != nil {
		logger.Error("获取用户失败",
			zap.String("id", in.Id),
			zap.Error(err))
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Error(codes.NotFound, "用户不存在")
		}
		return nil, status.Error(codes.Internal, "系统错误")
	}

	// 成功获取用户信息时记录info日志
	logger.Info("成功获取用户信息",
		zap.String("userId", strconv.Itoa(int(userModel.ID))),
		zap.String("nickname", userModel.Nickname))

	return &user.GetUserResp{
		User: &user.UserEntity{
			Id:       int64(userModel.ID),
			Avatar:   userModel.Avatar,
			Nickname: userModel.Nickname,
			Phone:    userModel.Phone,
			Status:   int32(userModel.Status),
			Sex:      int32(userModel.Sex),
		},
	}, nil
}
