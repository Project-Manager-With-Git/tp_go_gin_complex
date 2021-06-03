package serv

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"time"

	"tp_go_gin_complex/apis"
	"tp_go_gin_complex/models"

	_ "tp_go_gin_complex/docs"

	log "github.com/Golang-Tools/loggerhelper"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	ginlogrus "github.com/toorop/gin-logrus"
)

type Serv struct {
	App_Version            string   `json:"app_version" jsonschema:"required,title=v,description=应用版本"`
	App_Name               string   `json:"app_name" jsonschema:"required,title=n,description=应用名"`
	Log_Level              string   `json:"log_level" jsonschema:"required,title=l,description=log等级,enum=TRACE,enum=DEBUG,enum=INFO,enum=WARN,enum=ERROR"`
	Address                string   `json:"address" jsonschema:"required,title=a,description=启动地址"`
	Published_Address      string   `json:"published_address" jsonschema:"title=p,description=外部访问地址"`
	DB_URL                 string   `json:"db_url" jsonschema:"required,description=数据库连接url"`
	Cros_Allow_Origins     []string `json:"cros_allow_origins" jsonschema:"description=跨域允许的域名"`
	Cros_Allow_Credentials bool     `json:"cros_allow_credentials" jsonschema:"description=跨域是否需要证书"`
	Cros_Allow_Headers     []string `json:"cros_allow_headers" jsonschema:"description=跨域允许的头"`
	Cros_Expose_Headers    []string `json:"cros_expose_headers" jsonschema:"description=跨域暴露的头"`
	Static_Page_Dir        string   `json:"static_page_dir" jsonschema:"description=静态页面存放的文件夹"`
	Static_Source_Dir      string   `json:"static_source_dir" jsonschema:"description=静态资源存放的文件夹"`
	Serv_Cert_Path         string   `json:"serv_cert_path" jsonschema:"description=服务证书位置"`
	Serv_Key_Path          string   `json:"serv_key_path" jsonschema:"description=服务证书的私钥位置"`
	Ca_Cert_Path           string   `json:"ca_cert_path" jsonschema:"description=根证书位置"`
	Client_Crl_Path        string   `json:"client_crl_path" jsonschema:"description=客户端证书黑名单"`

	app *gin.Engine
}

func (s *Serv) runserv() {
	srv := &http.Server{
		Addr:    s.Address,
		Handler: s.app,
	}
	usetls := false
	if s.Serv_Cert_Path != "" && s.Serv_Key_Path != "" {
		usetls = true
		//双向认证
		if s.Ca_Cert_Path != "" {
			capool := x509.NewCertPool()
			caCrt, err := ioutil.ReadFile(s.Ca_Cert_Path)
			if err != nil {
				log.Error("read pem file error:", log.Dict{"err": err.Error(), "path": s.Ca_Cert_Path})
				os.Exit(2)
			}
			capool.AppendCertsFromPEM(caCrt)
			tlsconf := &tls.Config{
				RootCAs:    capool,
				ClientAuth: tls.RequireAndVerifyClientCert, // 检验客户端证书
			}
			//指定client名单
			if s.Client_Crl_Path != "" {
				clipool := x509.NewCertPool()
				cliCrt, err := ioutil.ReadFile(s.Client_Crl_Path)
				if err != nil {
					log.Error("read pem file error:", log.Dict{"err": err.Error(), "path": s.Client_Crl_Path})
					os.Exit(2)
				}
				clipool.AppendCertsFromPEM(cliCrt)
				tlsconf.ClientCAs = clipool
			}
			srv.TLSConfig = tlsconf
		}
	}

	//启动服务器
	go func() {
		// 服务连接
		log.Info("servrt start", log.Dict{"config": s})
		var err error
		if usetls {
			err = srv.ListenAndServeTLS(s.Serv_Cert_Path, s.Serv_Key_Path)

		} else {
			err = srv.ListenAndServe()
		}
		if err != nil && err != http.ErrServerClosed {
			log.Error("listen error", log.Dict{"err": err.Error()})
			os.Exit(2)
		}
	}()
	// 等待中断信号以优雅地关闭服务器（设置 5 秒的超时时间）
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Info("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := srv.Shutdown(ctx)
	if err != nil {
		log.Error("server shutdown error", log.Dict{"err": err.Error()})
		os.Exit(2)
	}
	log.Info("Server exiting")
}

func (s *Serv) Main() {
	// 初始化log
	log.Init(s.Log_Level, map[string]interface{}{
		"app_version": s.App_Version,
		"app_name":    s.App_Name,
	})

	s.app = gin.New()

	//初始化数据模型
	models.Init(s.DB_URL)

	// 跨域
	// router.Router.Use(cors.Default())
	corsconf := cors.DefaultConfig()
	schema := "http"
	if s.Serv_Cert_Path != "" && s.Serv_Key_Path != "" {
		schema = "https"
	}
	corsconf.AllowOrigins = []string{fmt.Sprintf("%s://%s", schema, s.Address)}
	if s.Published_Address != "" {
		corsconf.AllowOrigins = append(corsconf.AllowOrigins, fmt.Sprintf("%s://%s", schema, s.Published_Address))
	}
	if s.Cros_Allow_Origins != nil {
		corsconf.AllowOrigins = append(corsconf.AllowOrigins, s.Cros_Allow_Origins...)
	}
	corsconf.AllowCredentials = s.Cros_Allow_Credentials
	if s.Cros_Allow_Headers != nil {
		for _, header := range s.Cros_Allow_Headers {
			corsconf.AddAllowHeaders(header)
		}
	}
	if s.Cros_Expose_Headers != nil {
		for _, header := range s.Cros_Expose_Headers {
			corsconf.AddExposeHeaders(header)
		}
	}
	log.Info("cros config", log.Dict{"config": corsconf})
	s.app.Use(cors.New(corsconf))

	// log
	s.app.Use(ginlogrus.Logger(log.Logger), gin.Recovery())

	//注册静态路由
	if s.Static_Page_Dir != "" {
		s.app.Use(static.Serve("/", static.LocalFile(s.Static_Page_Dir, false)))
	}
	if s.Static_Source_Dir != "" {
		s.app.Use(static.Serve("/static", static.LocalFile(s.Static_Source_Dir, false)))
	}
	if s.Log_Level == "DEBUG" {
		url := ginSwagger.URL("http://localhost:5000/swagger/doc.json") // The url pointing to API definition
		s.app.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
	}
	apis.Init(s.app)
	// 启动服务
	s.runserv()
}

var ServNode = Serv{
	App_Version:         "1.0.0",
	App_Name:            "tp_go_gin_complex",
	Log_Level:           "DEBUG",
	Address:             "0.0.0.0:5000",
	DB_URL:              "sqlite://:memory:",
	Cros_Allow_Origins:  []string{},
	Cros_Allow_Headers:  []string{},
	Cros_Expose_Headers: []string{},
}
