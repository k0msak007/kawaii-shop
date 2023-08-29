package main

import (
	"os"

	"github.com/k0msak007/kawaii-shop/config"
	"github.com/k0msak007/kawaii-shop/modules/servers"
	"github.com/k0msak007/kawaii-shop/pkg/databases"
)

func envPath() string {
	if len(os.Args) == 1 {
		return ".env"
	} else {
		return os.Args[1]
	}
}

func main() {
	cfg := config.LoadConfig(envPath())

	db := databases.DbConnect(cfg.Db())
	defer db.Close() // defer จะทำงานท้ายสุดก่อน return

	servers.NewServer(cfg, db).Start()
}
