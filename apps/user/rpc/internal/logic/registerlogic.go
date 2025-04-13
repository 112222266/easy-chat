package logic

import (
	"context"
	"easy-chat/apps/user/models"
	"easy-chat/apps/user/rpc/internal/svc"
	"easy-chat/apps/user/rpc/user"
	utils2 "easy-chat/apps/user/rpc/utils"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type RegisterLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}
func (l *RegisterLogic) Register(in *user.RegisterReq) (*user.RegisterResp, error) {
	logger := l.svcCtx.Logger

	// 1. 密码哈希
	pwHashed, err := utils2.HashPassword(in.Password)
	if err != nil {
		logger.Error("密码哈希失败", zap.Error(err))
		return nil, status.Error(codes.Internal, "系统错误")
	}

	// 2. 创建用户
	userModel := &models.User{
		Avatar:   in.Avatar,
		Nickname: in.Nickname,
		Phone:    in.Phone,
		Password: pwHashed,
		Sex:      int8(in.Sex),
	}
	if err := l.svcCtx.DB.Create(&userModel).Error; err != nil {
		logger.Error("创建用户失败", zap.String("phone", in.Phone), zap.Error(err))
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, status.Error(codes.AlreadyExists, "手机号已注册")
		}
		return nil, status.Error(codes.Internal, "系统错误")
	}
	fmt.Printf("%+v", userModel)
	// 3. 生成Token
	auth := l.svcCtx.Config.Jwt
	token, err := utils2.GenerateToken(auth.Secret, int(auth.Expire), userModel.ID)
	if err != nil {
		logger.Error("生成Token失败", zap.Uint("userID", userModel.ID), zap.Error(err))
		return nil, status.Error(codes.Internal, "系统错误")
	}

	logger.Info("注册成功", zap.Uint("userID", userModel.ID))
	return &user.RegisterResp{
		Token:  token,
		Expire: auth.Expire,
	}, nil
}
