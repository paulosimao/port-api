package main

import "github.com/paulosimao/ports-api/lib/server"

func main() {
	err := server.Run()
	if err != nil {
		panic(err)
	}
}
