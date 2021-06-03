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
// @Accept application/json
// @Produce  application/json
// @Success 200 {object} UserListResponse "用户列表响应信息,会展示用户数量"
// @Failure 500 {string} ResultResponse "服务器处理失败"
// @Router /v1/api/user [get]
func (s *UserListSource) Get(c *gin.Context) {
	cnt, err := user.Count(models.DB)
	if err != nil {
		c.PureJSON(http.StatusInternalServerError, &ResultResponse{Message: err.Error()})
		return
	}
	result := UserListResponse{
		Description: "测试api,User总览",
		UserCount:   cnt,
		Links: []LinkResponse{
			{
				URI:         "/user",
				Method:      "POST",
				Description: "创建一个新用户",
			},
			{
				URI:         "/user/<int:uid>",
				Method:      "GET",
				Description: "用户号为<id>的用户信息",
			},
			{
				URI:         "/user/<int:uid>",
				Method:      "PUT",
				Description: "更新用户号为<id>用户信息",
			},
			{
				URI:         "/user/<int:uid>",
				Method:      "DELETE",
				Description: "删除用户号为<id>用户",
			},
		},
	}
	c.PureJSON(200, gin.H{"Result": result})
}

// @Summary 创建新用户
// @Tags user
// @Accept application/json
// @Produce application/json
// @Param name body UserCreateQuery true "用户名信息"
// @Success 200 {object} user.User "用户信息"
// @Failure 400 {string} ResultResponse "请求数据不符合要求"
// @Failure 500 {string} ResultResponse "服务器处理失败"
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
