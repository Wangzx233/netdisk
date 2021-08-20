package main

import (
	"netdisk/dao"
	"netdisk/router"
)

func main() {
	dao.InitDB()
	router.Start()
}
