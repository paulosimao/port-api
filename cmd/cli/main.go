package main

import  "github.com/paulosimao/ports-api/lib/cli"

func main() {
	err := cli.Run()
	if err != nil {
		panic(err)
	}
}
