package main

import (
	"easy-chat/apps/user/models"
	"flag"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"easy-chat/apps/user/rpc/internal/config"
	"easy-chat/apps/user/rpc/internal/server"
	"easy-chat/apps/user/rpc/internal/svc"
	"easy-chat/apps/user/rpc/user"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "./etc/user.yaml", "the config file")

func main() {
	migrator := flag.Bool("db", false, "是否创建表")
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		user.RegisterUserServer(grpcServer, server.NewUserServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	if *migrator {
		db, err := gorm.Open(mysql.Open(c.Database.Source), &gorm.Config{})
		if err != nil {
			panic(err)
		}
		if err := db.AutoMigrate(&models.User{}); err != nil {
			panic(err)
		}
		return
	}
	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
