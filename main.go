package main

import (
	"blogx_server/core"
	"blogx_server/flags"
)

func main() {
	flags.Parse()
	core.ReadConf()
}
