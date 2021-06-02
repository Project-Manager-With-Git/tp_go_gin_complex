package usernamespace

import (
	"net/http"
	"tp_go_gin_complex/models"
	"tp_go_gin_complex/models/user"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type UserListSource struct {
}

// @Summary 获取用户列表信息
// @Tags user
// @Produce  json
// @Success 200 {object} json
// @Router /v1/api/user [get]
func (s *UserListSource) Get(c *gin.Context) {
	cnt, err := user.Count(models.DB)
	if err != nil {
		c.PureJSON(http.StatusInternalServerError, &ResultResponse{Message: err.Error()})
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
	c.PureJSON(200, gin.H{"Result": result})
}

// @Summary 创建新用户
// @Tags user
// @accept json
// @Produce json
// @Param name body UserCreateQuery true "用户名"
// @Success 200 {object} user.User "{"Name":"1234","ID":1}"
// @Router /v1/api/user [post]
func (s *UserListSource) Post(c *gin.Context) {
	// 请求参数校验
	uinput := &UserCreateQuery{}
	err := c.ShouldBindBodyWith(uinput, binding.JSON)
	if err != nil {
		c.PureJSON(http.StatusBadRequest, &ResultResponse{Message: err.Error()})
		return
	}

	//创建用户
	u := &user.User{
		Name: uinput.Name,
	}
	err = u.Save(models.DB)
	if err != nil {
		c.PureJSON(http.StatusInternalServerError, &ResultResponse{Message: err.Error()})
		return
	}
	c.PureJSON(200, u)
}
