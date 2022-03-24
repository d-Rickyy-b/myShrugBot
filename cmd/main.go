package main

import (
	"flag"
	"myShrugBot/internal/bot"
	"myShrugBot/internal/config"
	"myShrugBot/internal/prometheus"
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
