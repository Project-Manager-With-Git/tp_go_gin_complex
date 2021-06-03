package main

import (
	"os"

	"tp_go_gin_complex/serv"

	log "github.com/Golang-Tools/loggerhelper"
	s "github.com/Golang-Tools/schema-entry-go"
)

// @title tp_go_gin_complex
// @version 1.0
// @description 测试

// @contact.name hsz
// @contact.email hsz1273327@gmail.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost
// @BasePath /
func main() {
	serv, err := s.New(&s.EntryPointMeta{Name: "tp_go_gin_complex", Usage: "tp_go_gin_complex [options]"}, &serv.ServNode)
	if err != nil {
		log.Error("create entrypoint get error", log.Dict{"err": err})
		os.Exit(2)
	}
	serv.Parse(os.Args)
}
