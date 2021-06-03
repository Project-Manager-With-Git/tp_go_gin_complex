// GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by swaggo/swag

package docs

import (
	"bytes"
	"encoding/json"
	"strings"

	"github.com/alecthomas/template"
	"github.com/swaggo/swag"
)

var doc = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{.Description}}",
        "title": "{{.Title}}",
        "contact": {
            "name": "hsz",
            "email": "hsz1273327@gmail.com"
        },
        "license": {
            "name": "Apache 2.0",
            "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/user": {
            "get": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "user"
                ],
                "summary": "获取用户列表信息",
                "responses": {
                    "200": {
                        "description": "用户列表响应信息,会展示用户数量",
                        "schema": {
                            "$ref": "#/definitions/usernamespace.UserListResponse"
                        }
                    },
                    "500": {
                        "description": "服务器处理失败",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            },
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "user"
                ],
                "summary": "创建新用户",
                "parameters": [
                    {
                        "description": "用户名信息",
                        "name": "name",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/usernamespace.UserCreateQuery"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "用户信息",
                        "schema": {
                            "$ref": "#/definitions/user.User"
                        }
                    },
                    "400": {
                        "description": "请求数据不符合要求",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "服务器处理失败",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/user/{uid}": {
            "get": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "user"
                ],
                "summary": "获取用户列表信息",
                "parameters": [
                    {
                        "minimum": 1,
                        "type": "integer",
                        "description": "User ID",
                        "name": "uid",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/user.User"
                        }
                    },
                    "400": {
                        "description": "请求数据不符合要求",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "404": {
                        "description": "未找到指定资源",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            },
            "put": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "user"
                ],
                "summary": "更新指定用户信息",
                "parameters": [
                    {
                        "minimum": 1,
                        "type": "integer",
                        "description": "User ID",
                        "name": "uid",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/user.User"
                        }
                    },
                    "400": {
                        "description": "请求数据不符合要求",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "404": {
                        "description": "未找到指定资源",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "服务器处理失败",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            },
            "delete": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "user"
                ],
                "summary": "删除指定用户",
                "parameters": [
                    {
                        "minimum": 1,
                        "type": "integer",
                        "description": "User ID",
                        "name": "uid",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/user.User"
                        }
                    },
                    "400": {
                        "description": "请求数据不符合要求",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "404": {
                        "description": "未找到指定资源",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "服务器处理失败",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "user.User": {
            "type": "object",
            "required": [
                "ID"
            ],
            "properties": {
                "ID": {
                    "type": "integer"
                },
                "Name": {
                    "type": "string"
                }
            }
        },
        "usernamespace.LinkResponse": {
            "type": "object",
            "properties": {
                "description": {
                    "type": "string"
                },
                "method": {
                    "type": "string"
                },
                "uri": {
                    "type": "string"
                }
            }
        },
        "usernamespace.UserCreateQuery": {
            "type": "object",
            "properties": {
                "Name": {
                    "type": "string"
                }
            }
        },
        "usernamespace.UserListResponse": {
            "type": "object",
            "properties": {
                "Description": {
                    "type": "string"
                },
                "Links": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/usernamespace.LinkResponse"
                    }
                },
                "UserCount": {
                    "type": "integer"
                }
            }
        }
    }
}`

type swaggerInfo struct {
	Version     string
	Host        string
	BasePath    string
	Schemes     []string
	Title       string
	Description string
}

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = swaggerInfo{
	Version:     "1.0.0",
	Host:        "localhost:5000",
	BasePath:    "/v1_0_0/api",
	Schemes:     []string{},
	Title:       "tp_go_gin_complex",
	Description: "测试",
}

type s struct{}

func (s *s) ReadDoc() string {
	sInfo := SwaggerInfo
	sInfo.Description = strings.Replace(sInfo.Description, "\n", "\\n", -1)

	t, err := template.New("swagger_info").Funcs(template.FuncMap{
		"marshal": func(v interface{}) string {
			a, _ := json.Marshal(v)
			return string(a)
		},
	}).Parse(doc)
	if err != nil {
		return doc
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, sInfo); err != nil {
		return doc
	}

	return tpl.String()
}

func init() {
	swag.Register(swag.Name, &s{})
}
