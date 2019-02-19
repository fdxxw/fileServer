package db

import (
	"le5le.com/fileServer/db/mongo"
)

// Init 初始化数据库连接
func Init() bool {
	return mongo.Init()
}
