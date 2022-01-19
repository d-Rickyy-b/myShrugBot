package main

import (
	"flag"
	"github.com/d-Rickyy-b/myShrugBot/internal/bot"
	"github.com/d-Rickyy-b/myShrugBot/internal/config"
)

func main() {
	var configFile = flag.String("config", "config.json", "Path to config file")
	flag.Parse()

	c, _ := config.ReadConfig(*configFile)
	bot.StartBot(c)
}
