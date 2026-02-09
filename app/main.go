package main

import (
	"github.com/slickip/Subscription-service/internal/config"
	"github.com/slickip/Subscription-service/internal/db"
)

func main() {
	cfg := config.Load()
	dbConn := db.New(cfg)
	_ = dbConn
}
