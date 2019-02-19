package main

import (
	"log"
	"runtime"

	lumberjack "gopkg.in/natefinch/lumberjack.v2"
	"le5le.com/fileServer/config"
	"le5le.com/fileServer/db"
	"le5le.com/fileServer/db/mongo"
	"le5le.com/fileServer/router"
)

func main() {
	// 处理panic
	defer func() {
		if err := recover(); err != nil {
			log.Printf("[panic] %v\r\n", err)
		}
	}()

	// 初始化配置
	config.Init()

	// 设置日志
	if config.App.Log.Filename != "" {
		log.SetOutput(&lumberjack.Logger{
			Filename:   config.App.Log.Filename,
			MaxSize:    config.App.Log.MaxSize, // mb
			MaxBackups: config.App.Log.MaxBackups,
			MaxAge:     config.App.Log.MaxAge, // days
		})
	}

	// 最大cpu使用核心数
	runtime.GOMAXPROCS(config.App.CPU)

	// 数据库连接
	if !db.Init() {
		return
	}
	defer mongo.Session.Close()

	// 监听路由
	router.Listen()
}
