package svc

import (
	"easy-chat/apps/user/rpc/internal/config"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"os"
	"path/filepath"
)

type ServiceContext struct {
	Config config.Config
	DB     *gorm.DB
	Logger *zap.Logger
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		DB:     newDB(c),
		Logger: newLogger(c),
		Config: c,
	}
}

func newDB(c config.Config) *gorm.DB {
	db, err := gorm.Open(mysql.Open(c.Database.Source), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return db
}

func newLogger(c config.Config) *zap.Logger {
	logDir, appName := c.Log.Dir, c.Log.AppName
	// 确保日志目录存在
	if err := os.MkdirAll(logDir, 0755); err != nil {
		panic(err)
	}

	// 1. 定义日志级别到文件的映射
	levelFiles := map[zapcore.Level]string{
		zapcore.DebugLevel:  filepath.Join(logDir, appName+".debug.log"),
		zapcore.InfoLevel:   filepath.Join(logDir, appName+".info.log"),
		zapcore.WarnLevel:   filepath.Join(logDir, appName+".warn.log"),
		zapcore.ErrorLevel:  filepath.Join(logDir, appName+".error.log"),
		zapcore.DPanicLevel: filepath.Join(logDir, appName+".dpanic.log"),
		zapcore.PanicLevel:  filepath.Join(logDir, appName+".panic.log"),
		zapcore.FatalLevel:  filepath.Join(logDir, appName+".fatal.log"),
	}

	// 2. 创建各个级别的文件core
	cores := make([]zapcore.Core, 0, len(levelFiles)+1) // +1 给控制台

	// 控制台输出: 全级别输出，带颜色
	consoleEncoder := zapcore.NewConsoleEncoder(zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalColorLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	})
	cores = append(cores, zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), zapcore.DebugLevel))

	// 文件输出: 每个级别单独文件
	for level, filename := range levelFiles {
		// 使用lumberjack实现日志轮转
		// 每个文件最大100MB，最多保留5个备份，保留30天
		fileWriter := zapcore.AddSync(&lumberjack.Logger{
			Filename:   filename,
			MaxSize:    100, // MB
			MaxBackups: 5,
			MaxAge:     30, // days
			Compress:   true,
			LocalTime:  true,
		})

		fileEncoder := zapcore.NewJSONEncoder(zapcore.EncoderConfig{
			TimeKey:        "ts",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			FunctionKey:    zapcore.OmitKey,
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.StringDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		})

		// 每个core只处理特定级别的日志
		levelEnabler := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
			return lvl == level
		})

		cores = append(cores, zapcore.NewCore(fileEncoder, fileWriter, levelEnabler))
	}

	// 3. 创建最终Logger
	core := zapcore.NewTee(cores...)
	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))

	return logger
}
