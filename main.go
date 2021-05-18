package main

import (
	"os"

	"tp_go_gin_complex/serv"

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
	root, _ := s.New(&s.EntryPointMeta{Name: "tp_go_gin_complex", Usage: "tp_go_gin_complex [cmd [options]]"})
	serv, _ := s.New(&s.EntryPointMeta{Name: "serv", Usage: "tp_go_gin_complex serv"}, &serv.ServNode)
	// nodec, _ := s.New(&s.EntryPointMeta{Name: "ping", Usage: "tp_go_gin_complex "}, &C{
	// 	Field: []int{1, 2, 3},
	// })
	s.RegistSubNode(root, serv)
	root.Parse(os.Args)
}
