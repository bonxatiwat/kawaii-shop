package main

import (
	"fmt"
	"os"

	"github.com/bonxatiwat/kawaii-shop-tutortial/config"
	"github.com/bonxatiwat/kawaii-shop-tutortial/modules/servers"
	"github.com/bonxatiwat/kawaii-shop-tutortial/pkg/database"
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
	db := database.DbConnect(cfg.Db())
	defer db.Close()

	fmt.Println(db)

	servers.NewServer(cfg, db).Start()
}
