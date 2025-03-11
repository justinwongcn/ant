package main

import (
	"fmt"

	"github.com/roovet/ant"
)

func main() {
	server := ant.NewHTTPServer()

	server.Handle("/hello/{name}", func(ctx *ant.Context) {
		name := ctx.Req.PathValue("name")
		if name != "" {
			ctx.RespData = []byte(fmt.Sprintf("hello %s!", name))
		} else {
			ctx.RespData = []byte("hello world!")
		}
	})

	err := server.Run(":8080")
	if err != nil {
		panic(err)
	}
}