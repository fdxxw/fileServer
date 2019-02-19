package file

import (
	"github.com/kataras/iris"
	"gopkg.in/mgo.v2/bson"
	"le5le.com/fileServer/keys"
	"le5le.com/fileServer/middlewares"
)

// Route file模块路由
func Route(route *iris.Application) {

	route.Post("/api/file", func(ctx iris.Context) {
		if ctx.GetContentLength() > 10<<20 {
			ctx.StatusCode(iris.StatusRequestEntityTooLarge)
			ctx.JSON(bson.M{"error": keys.ErrorFileMaxSize})
			return
		}
		ctx.Next()
	}, middlewares.Auth, Upload)

	route.Get("/image/{filename}", Image)
	route.Get("/file/{filename}", Image)
}
