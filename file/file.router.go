package file

import (
	"fileServer/middlewares"

	"github.com/kataras/iris"
)

// Route file模块路由
func Route(route *iris.Application) {

	route.Get("/file/{path:path}", File)
	route.Post("/api/file", middlewares.Usr, Limit, UploadFile)
	route.Patch("/api/file/{path:path}", middlewares.Usr, FileAttr)
	route.Delete("/api/file/{path:path}", middlewares.Usr, DelFile)

	route.Get("/image/{path:path}", Image)
	route.Post("/api/image", middlewares.Auth, Limit, UploadImage)
	route.Patch("/api/image/{path:path}", middlewares.Auth, FileAttr)
	route.Delete("/api/image/{path:path}", middlewares.Auth, DelFile)
}
