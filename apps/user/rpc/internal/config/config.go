package config

import "github.com/zeromicro/go-zero/zrpc"

type Config struct {
	zrpc.RpcServerConf
	Database struct {
		Source string
	}
	Log struct {
		Dir     string
		AppName string
	}
	Jwt struct {
		Expire int64
		Secret string
	}
}
