package main

import "github.com/kolibriee/trade-metrics/internal/app"

const (
	configName = "config"
	configsDir = "configs"
)

func main() {
	app.Run(configsDir, configName)
}
