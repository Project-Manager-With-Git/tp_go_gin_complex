package models

import (
	"os"
	"tp_go_gin_complex/models/user"

	log "github.com/Golang-Tools/loggerhelper"
)

func Init(url string, opts ...Option) {
	DB.Regist(user.RegistCallback)
	err := DB.InitFromURL(url)
	if err != nil {
		log.Error("db proxy InitFromURL error ", log.Dict{"err": err.Error()})
		os.Exit(21)
	}
	DefaultOption := &Options{
		ShowSQL: true,
	}
	for _, opt := range opts {
		opt.Apply(DefaultOption)
	}
	if DefaultOption.MaxIdle > 0 {
		DB.SetMaxIdleConns(DefaultOption.MaxIdle)
	}
	if DefaultOption.MaxOpenConns > 0 {
		DB.SetMaxOpenConns(DefaultOption.MaxOpenConns)
	}
	log.Info("model inited done", log.Dict{"url": url, "options": DefaultOption})
}
