package logic

import (
	"context"
	"easy-chat/apps/user/models"
	"easy-chat/apps/user/rpc/internal/svc"
	"easy-chat/apps/user/rpc/user"
	utils2 "easy-chat/apps/user/rpc/utils"
	"errors"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *LoginLogic) Login(in *user.LoginReq) (*user.LoginResp, error) {
	logger := l.svcCtx.Logger

	// 1. 查询用户
	userModel := &models.User{}
	if err := l.svcCtx.DB.Where("phone = ?", in.Phone).Take(userModel).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Warn("用户不存在", zap.String("phone", in.Phone))
			return nil, status.Error(codes.NotFound, "用户不存在或密码错误")
		}
		logger.Error("查询用户失败", zap.String("phone", in.Phone), zap.Error(err))
		return nil, status.Error(codes.Internal, "系统错误")
	}

	// 2. 校验密码
	if !utils2.CheckPassword(userModel.Password, in.Password) {
		logger.Warn("密码验证失败", zap.String("phone", in.Phone))
		return nil, status.Error(codes.Unauthenticated, "用户不存在或密码错误")
	}

	// 3. 生成Token
	auth := l.svcCtx.Config.Jwt
	token, err := utils2.GenerateToken(auth.Secret, int(auth.Expire), userModel.ID)
	if err != nil {
		logger.Error("生成Token失败",
			zap.Uint("userID", userModel.ID),
			zap.Error(err))
		return nil, status.Error(codes.Internal, "系统错误")
	}

	logger.Info("登录成功", zap.Uint("userID", userModel.ID))
	return &user.LoginResp{
		Token:  token,
		Expire: auth.Expire,
	}, nil
}
