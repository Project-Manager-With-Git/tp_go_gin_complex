package usernamespace

import (
	"net/http"
	"tp_go_gin_complex/models"
	"tp_go_gin_complex/models/user"

	"github.com/gin-gonic/gin"
)

type UserListSource struct {
}

func (s *UserListSource) Get(ctx *gin.Context) {
	cnt, err := user.Count(models.DB)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}
	result := map[string]interface{}{
		"Description": "测试api,User总览",
		"UserCount":   cnt,
		"Links": []map[string]interface{}{
			{
				"uri":         "/user",
				"method":      "POST",
				"description": "创建一个新用户",
			},
			{
				"uri":         "/user/<int:uid>",
				"method":      "GET",
				"description": "用户号为<id>的用户信息",
			},
			{
				"uri":         "/user/<int:uid>",
				"method":      "PUT",
				"description": "更新用户号为<id>用户信息",
			},
			{
				"uri":         "/user/<int:uid>",
				"method":      "DELETE",
				"description": "删除用户号为<id>用户",
			},
		},
	}
	ctx.JSON(200, gin.H{"Result": result})
}

type UserCreateQuery struct {
	Name string `json:"Name"`
}

func (s *UserListSource) Post(ctx *gin.Context) {
	// 请求参数校验
	uinput := UserCreateQuery{}
	err := ctx.ShouldBindJSON(&uinput)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
	}
	//创建用户
	u := &user.User{
		Name: uinput.Name,
	}
	err = u.Save(models.DB)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}
	ctx.JSON(200, u)
}
