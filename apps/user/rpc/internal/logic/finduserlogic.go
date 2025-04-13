package logic

import (
	"context"
	"easy-chat/apps/user/models"
	"easy-chat/apps/user/rpc/internal/svc"
	"easy-chat/apps/user/rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type FindUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFindUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindUserLogic {
	return &FindUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FindUserLogic) FindUser(in *user.FindUserReq) (*user.FindUserResp, error) {
	var usersModel []models.User
	if in.Name != "" {
		if err := l.svcCtx.DB.Where("nickname = ?", in.Name).Find(&usersModel).Error; err != nil {
			return nil, err
		}
	} else if in.Phone != "" {
		if err := l.svcCtx.DB.Where("phone = ?", in.Phone).Find(&usersModel).Error; err != nil {
			return nil, err
		}
	} else if len(in.Ids) != 0 {
		if err := l.svcCtx.DB.Find(&usersModel, in.Ids).Error; err != nil {
			return nil, err
		}
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
	return &user.FindUserResp{
		User: usersRes,
	}, nil
}
