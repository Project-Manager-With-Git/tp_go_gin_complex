package models

import (
	"os"
	"tp_go_gin_complex/models/user"

	log "github.com/Golang-Tools/loggerhelper"
)

func Init(url string) *xormProxy {
	DB.Regist(user.RegistCallback)

	err := DB.InitFromURL(url)
	if err != nil {
		log.Error("db proxy InitFromURL error ", log.Dict{"err": err.Error()})
		os.Exit(21)
	}
	log.Info("model inited done")
	return DB
}
