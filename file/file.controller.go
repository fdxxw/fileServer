package file

import (
	"log"
	"os"

	"github.com/kataras/iris"
	"gopkg.in/mgo.v2/bson"
	"le5le.com/fileServer/keys"
)

const (
	// CachePath 文件临时存储目录
	CachePath = "./out/uploads"
)

// Upload 文件上传
func Upload(ctx iris.Context) {
	ret := make(map[string]interface{})
	defer ctx.JSON(ret)

	file, info, err := ctx.FormFile("file")

	if err != nil {
		ret["error"] = keys.ErrorFile + err.Error()
		return
	}

	defer file.Close()
	filename := GetUniqueName(info.Filename)

	_, err = Put(filename, file, bson.M{"owner": bson.M{"id": ctx.Values().GetString("uid")}})
	if err == nil {
		ret["name"] = filename
		ret["url"] = "/image/" + filename
	} else {
		ret["error"] = keys.ErrorFileSave + err.Error()
	}
}

// Image 获取图片
func Image(ctx iris.Context) {
	filename := ctx.Params().Get("filename")
	if filename == "" {
		ctx.StatusCode(iris.StatusNotFound)
		ctx.JSON("")
		return
	}

	if !IsExist(CachePath) {
		err := os.MkdirAll(CachePath, os.ModePerm)
		if err != nil {
			log.Printf("[error]Can not create the dir of cache: err=%v; dir=%v\r\n", err, CachePath)
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.JSON("")
			return
		}
	}

	fullpath := CachePath + "/" + filename
	if !IsExist(fullpath) {
		err := Get(filename, CachePath)
		if err != nil {
			log.Printf("[error]file.Image: err=%v; filename=%v\r\n", err, filename)
			ctx.StatusCode(iris.StatusNotFound)
			ctx.JSON("")
			return
		}
	}

	w, _ := ctx.URLParamInt("w")
	h, _ := ctx.URLParamInt("h")

	if w > 0 || h > 0 {
		thumb, err := ImageThumbnail(fullpath, w, h)
		if err == nil {
			fullpath = thumb
		} else {
			log.Printf("[error]ImageThumbnail: err=%v; fullpath=%v\r\n", err, fullpath)
		}

	}

	ctx.ServeFile(fullpath, true)
}
