package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"tp_go_gin_complex/apis"
	"tp_go_gin_complex/models"
	"tp_go_gin_complex/models/user"

	log "github.com/Golang-Tools/loggerhelper"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

const DBURL = "sqlite://:memory:"

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.RedirectTrailingSlash = true
	apis.Init(r)
	models.Init(DBURL)
	return r
}

func TestPingRoute(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "pong", w.Body.String())
}
func TestUserListGetRoute(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1_0_0/api/user", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	log.Info("get response", log.Dict{"body": w.Body.String()})
	// json.Unmarshal(,)
	assert.Contains(t, w.Body.String(), "测试api,User总览")
}

func TestUserGetRoute(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1_0_0/api/user/1", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	log.Info("get response", log.Dict{"body": w.Body.String()})
	u := user.User{}
	json.Unmarshal(w.Body.Bytes(), &u)
	assert.Contains(t, u.Name, "admin")
}
