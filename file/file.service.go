package file

import (
	"errors"
	"fmt"
	"image"
	_ "image/gif"
	"image/jpeg"
	_ "image/png"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/nfnt/resize"
	"gopkg.in/mgo.v2/bson"

	"le5le.com/fileServer/config"
	"le5le.com/fileServer/db/mongo"
	"le5le.com/fileServer/keys"
)

// WalkDir 获取指定目录及所有子目录下的所有文件，可以匹配后缀过滤。
func WalkDir(dirPth, suffix string) ([]string, error) {
	files := make([]string, 0, 30)
	suffix = strings.ToUpper(suffix)
	err := filepath.Walk(dirPth, func(filename string, fi os.FileInfo, err error) error {
		// 忽略目录
		if fi.IsDir() {
			return nil
		}

		if strings.HasSuffix(strings.ToUpper(fi.Name()), suffix) {
			files = append(files, filename)
		}
		return nil
	})
	return files, err
}

// ReadFile 读取文件内容
func ReadFile(path string) (string, error) {
	fileHandle, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer fileHandle.Close()
	fileBytes, err := ioutil.ReadAll(fileHandle)
	return string(fileBytes), err
}

// IsExist 文件是否存在
func IsExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}

// GetUniqueName 获取一个唯一的文件名
func GetUniqueName(name string) string {
	id := bson.NewObjectId()
	newName := id.Hex()
	return newName
}

// Put 存储文件到数据库
func Put(name string, f interface{}, meta interface{}) (string, error) {
	if f == nil {
		return "", errors.New(keys.ErrorFileInfo)
	}
	mongoSession := mongo.Session.Clone()
	defer mongoSession.Close()

	mongoFile, err := mongoSession.DB(config.App.Mongo.Database).GridFS(mongo.Files).Create(name)
	if err != nil {
		return "", err
	}

	defer mongoFile.Close()
	var fileID string

	fid := mongoFile.Id()
	switch v := fid.(type) {
	case string:
		fileID = v
	case bson.ObjectId:
		fileID = v.Hex()
	}
	if name == "" {
		mongoFile.SetName(fileID)
	}
	mongoFile.SetMeta(meta)
	switch data := f.(type) {
	case []byte:
		_, err = mongoFile.Write(data)
	case io.Reader:
		_, err = io.Copy(mongoFile, data)
	default:
	}
	if err != nil {
		return "", err
	}
	return fileID, err
}

// Get 从数据库读取文件
func Get(name, p string) error {
	if name == "" {
		return errors.New(keys.ErrorParam)
	}
	mongoSession := mongo.Session.Clone()
	defer mongoSession.Close()

	file, err := mongoSession.DB(config.App.Mongo.Database).GridFS(mongo.Files).Open(name)
	if err != nil {
		return err
	}
	defer file.Close()
	fw, err := os.Create(p + "/" + name)
	if err != nil {
		return err
	}
	defer fw.Close()
	_, err = io.Copy(fw, file)
	return err
}

// ImageThumbnail 生成图片缩略图
func ImageThumbnail(src string, w, h int) (string, error) {
	filename := fmt.Sprintf("%s_%d_%d", src, w, h)

	if IsExist(filename) {
		return filename, nil
	}

	file, err := os.Open(src)
	if err != nil {
		return src, err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return src, err
	}

	thumb := resize.Thumbnail(uint(w), uint(h), img, resize.Lanczos3)
	out, err := os.Create(filename)
	if err != nil {
		return src, err
	}
	defer out.Close()

	// Write new image to file.
	err = jpeg.Encode(out, thumb, nil)

	if err != nil {
		return src, err
	}

	return filename, nil
}
