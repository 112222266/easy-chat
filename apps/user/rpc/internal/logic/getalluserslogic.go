package logic

import (
	"context"
	"easy-chat/apps/user/models"
	"easy-chat/apps/user/rpc/internal/svc"
	"easy-chat/apps/user/rpc/user"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetAllUsersLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetAllUsersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAllUsersLogic {
	return &GetAllUsersLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetAllUsersLogic) GetAllUsers(in *user.GetAllUsersReq) (*user.GetAllUsersResp, error) {
	var usersModel []models.User

	// 查询所有用户
	if err := l.svcCtx.DB.Find(&usersModel).Error; err != nil {
		l.svcCtx.Logger.Error("查询所有用户失败", zap.Error(err))
		return nil, status.Error(codes.Internal, "系统错误")
	}

	// 转换用户数据
	usersRes := make([]*user.UserEntity, len(usersModel))
	for i, u := range usersModel {
		usersRes[i] = &user.UserEntity{
			Id:       int64(u.ID),
			Avatar:   u.Avatar,
			Nickname: u.Nickname,
			Phone:    u.Phone,
			Status:   int32(u.Status),
			Sex:      int32(u.Sex),
		}
	}

	// 记录成功日志
	l.svcCtx.Logger.Info("获取所有用户成功",
		zap.Int("count", len(usersRes)))

	return &user.GetAllUsersResp{
		Users: usersRes,
	}, nil
}
