package main

import (
	"flag"
	"github.com/d-Rickyy-b/myShrugBot/internal/bot"
	"github.com/d-Rickyy-b/myShrugBot/internal/config"
	"github.com/d-Rickyy-b/myShrugBot/internal/prometheus"
)

func main() {
	var configFile = flag.String("config", "config.json", "Path to config file")
	flag.Parse()

	c, _ := config.ReadConfig(*configFile)

	if c.Prometheus.Enabled {
		go prometheus.StartPrometheusExporter(c.Prometheus.Listen)
	}
	bot.StartBot(c)
}
