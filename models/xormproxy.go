package models

import (
	"errors"
	"fmt"
	"net/url"
	"strings"

	log "github.com/Golang-Tools/loggerhelper"
	_ "github.com/mattn/go-sqlite3"
	"github.com/xormplus/xorm"
)

//ErrProxyAllreadySettedClient 代理已经设置过redis客户端对象
var ErrProxyAllreadySettedClient = errors.New("代理不能重复设置客户端对象")

//ErrUnSupportSchema 不支持的数据库类型
var ErrUnSupportSchema = errors.New("不支持的数据库类型")

//Callback redis操作的回调函数
type Callback func(cli xorm.EngineInterface) error

//xormProxy redis客户端的代理
type xormProxy struct {
	xorm.EngineInterface
	parallelcallback bool
	callBacks        []Callback
}

// New 创建一个新的数据库客户端代理
func New() *xormProxy {
	proxy := new(xormProxy)
	return proxy
}

// IsOk 检查代理是否已经可用
func (proxy *xormProxy) IsOk() bool {
	if proxy.EngineInterface == nil {
		return false
	}
	return true
}

//SetConnect 设置连接的客户端
//@params cli UniversalClient 满足redis.UniversalClient接口的对象的指针
func (proxy *xormProxy) SetConnect(cli *xorm.Engine) error {
	if proxy.IsOk() {
		return ErrProxyAllreadySettedClient
	}

	proxy.EngineInterface = cli
	if proxy.parallelcallback {
		for _, cb := range proxy.callBacks {
			go func(cb Callback) {
				err := cb(proxy.EngineInterface)
				if err != nil {
					log.Error("regist callback get error", log.Dict{"err": err})
				} else {
					log.Debug("regist callback done")
				}
			}(cb)
		}
	} else {
		for _, cb := range proxy.callBacks {
			err := cb(proxy.EngineInterface)
			if err != nil {
				log.Error("regist callback get error", log.Dict{"err": err})
			} else {
				log.Debug("regist callback done")
			}
		}
	}
	return nil
}

//InitFromURL 从URL条件初始化代理对象
func (proxy *xormProxy) InitFromURL(u string) error {
	U, err := url.Parse(u)
	if err != nil {
		return err
	}
	dataSourceName := strings.ReplaceAll(u, fmt.Sprintf("%s://", U.Scheme), "")
	proxy.parallelcallback = false
	switch U.Scheme {
	case "postgres":
		{
			engine, err := xorm.NewEngine("postgres", u)
			if err != nil {
				return err
			}
			return proxy.SetConnect(engine)
		}
	case "mysql":
		{
			engine, err := xorm.NewEngine("mysql", dataSourceName)
			if err != nil {
				return err
			}
			return proxy.SetConnect(engine)
		}
	case "sqlserver":
		{
			engine, err := xorm.NewEngine("mssql", u)
			if err != nil {
				return err
			}
			return proxy.SetConnect(engine)
		}
	case "sqlite":
		{
			engine, err := xorm.NewEngine("sqlite3", dataSourceName)
			if err != nil {
				return err
			}
			return proxy.SetConnect(engine)
		}
	default:
		{
			return ErrUnSupportSchema
		}
	}
}

//InitFromURLParallelCallback 从URL条件初始化代理对象
func (proxy *xormProxy) InitFromURLParallelCallback(u string) error {
	U, err := url.Parse(u)
	if err != nil {
		return err
	}
	dataSourceName := strings.ReplaceAll(u, fmt.Sprintf("%s://", U.Scheme), "")
	proxy.parallelcallback = true
	switch U.Scheme {
	case "postgres":
		{
			engine, err := xorm.NewEngine("postgres", u)
			if err != nil {
				return err
			}
			return proxy.SetConnect(engine)
		}
	case "mysql":
		{
			engine, err := xorm.NewEngine("mysql", dataSourceName)
			if err != nil {
				return err
			}
			return proxy.SetConnect(engine)
		}
	case "sqlserver":
		{
			engine, err := xorm.NewEngine("mssql", u)
			if err != nil {
				return err
			}
			return proxy.SetConnect(engine)
		}
	case "sqlite":
		{
			engine, err := xorm.NewEngine("sqlite3", dataSourceName)
			if err != nil {
				return err
			}
			return proxy.SetConnect(engine)
		}
	default:
		{
			return ErrUnSupportSchema
		}
	}
}

// Regist 注册回调函数,在init执行后执行回调函数
//如果对象已经设置了被代理客户端则无法再注册回调函数
func (proxy *xormProxy) Regist(cb Callback) error {
	if proxy.IsOk() {
		return ErrProxyAllreadySettedClient
	}
	proxy.callBacks = append(proxy.callBacks, cb)
	return nil
}

//DB 默认的xorm代理对象
var DB = New()
