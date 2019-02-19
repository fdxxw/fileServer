package middlewares

import (
	"fmt"
	"log"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/kataras/iris"
	"le5le.com/fileServer/config"
	"le5le.com/fileServer/keys"
	"le5le.com/fileServer/utils"
)

// Auth 身份认证中间件
func Auth(ctx iris.Context) {
	if config.App.Jwt == "" {
		ctx.Values().Set("uid", "0")
		ctx.Values().Set("role", "system")
		ctx.Next()
		return
	}

	// 获取header
	data := ctx.GetHeader("Authorization")
	if data == "" {
		unAuth(ctx)
		return
	}

	// jwt校验
	token, err := jwt.Parse(data, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("签名方法错误: %v", token.Header["alg"])
		}
		return []byte(config.App.Jwt), nil
	})

	if err != nil {
		log.Printf("Jwt parse error: %s, token=%s, jwt=%s", err, data, config.App.Jwt)
		unAuth(ctx)
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		unAuth(ctx)
		return
	}

	// 设置uid和role
	uid := utils.String(claims["uid"])
	if uid == "" {
		unAuth(ctx)
		return
	}

	ctx.Values().Set("uid", uid)
	ctx.Values().Set("role", utils.String(claims["role"]))

	ctx.Next()
}

func unAuth(ctx iris.Context) {
	ctx.StatusCode(iris.StatusUnauthorized)
	ret := make(map[string]interface{})
	ret["error"] = keys.ErrorNeedSign
	ctx.JSON(ret)
}
